package gritrequester

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RequesterObj struct {
	Client HTTPClient
	Token  *TokenCache
	Confs  StaticConfig
}

type MsRequest struct {
	MSName string
	Method string
	Path   string
	Body   any
}

type ResponseData[T any] struct {
	Data       T
	PageCursor string
}

func NewRequestObj(conf StaticConfig) *RequesterObj {
	return &RequesterObj{
		Client: &http.Client{},
		Token:  NewTokenCache(),
		Confs:  conf,
	}
}

func DoMsRequest[ResponseBody any](
	requester *RequesterObj,
	msRequest MsRequest,
	retry bool,
) (ResponseData[ResponseBody], error) {
	var responseData ResponseData[ResponseBody]

	request, err := newRequest(
		requester,
		msRequest,
	)

	if err != nil {
		return responseData, err
	}

	statusCode, response, err := execRequest[ResponseBody](requester, msRequest.MSName, request)

	if statusCode == http.StatusUnauthorized && retry {
		requester.Token.Delete(msRequest.MSName)
		return DoMsRequest[ResponseBody](requester, msRequest, false)
	}

	return response, err
}

func execRequest[ResponseBody any](
	requester *RequesterObj,
	msName string,
	request *http.Request,
) (int, ResponseData[ResponseBody], error) {
	var responseData ResponseData[ResponseBody]

	response, err := requester.Client.Do(request)
	if err != nil {
		return 0, responseData, err
	}

	defer response.Body.Close()

	responseData.PageCursor = response.Header.Get("X-Page-Cursor")

	updateServiceToken(requester, msName, response.Header)

	responseBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return response.StatusCode, responseData, err
	}

	if response.StatusCode != 204 {
		err = json.Unmarshal(responseBodyBytes, &responseData.Data)
	}

	if err == nil && response.StatusCode > 299 {
		err = fmt.Errorf(
			"request to %s failed [%d]: %s",
			request.URL.String(),
			response.StatusCode,
			err,
		)
	}

	return response.StatusCode, responseData, err
}

func updateServiceToken(requester *RequesterObj, msName string, responseHeader http.Header) {
	actualToken, ok := requester.Token.Get(msName)
	if !ok {
		return
	}

	responseToken := responseHeader.Get("X-Token")

	if len(responseToken) > 0 && actualToken != responseToken {
		requester.Token.Set(msName, responseToken)
	}
}

func newRequest(
	requester *RequesterObj,
	msRequest MsRequest,
) (*http.Request, error) {
	conf, err := requester.Confs.Get(msRequest.MSName)
	if err != nil {
		return nil, err
	}

	encodedBody, err := json.Marshal(msRequest.Body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(
		msRequest.Method,
		conf.BaseUrl+msRequest.Path,
		bytes.NewBuffer(encodedBody),
	)

	if err != nil {
		return nil, err
	}

	err = setRequestHeaders(requester, request, msRequest.MSName)
	if err != nil {
		return nil, err
	}

	return request, nil
}

func setRequestHeaders(
	requester *RequesterObj,
	request *http.Request,
	MSName string,
) error {
	conf, _ := requester.Confs.Get(MSName)
	token, ok := requester.Token.Get(MSName)
	if ok {
		request.Header.Set("Authorization", token)
		request.Header.Set("Context", conf.Context)
		return nil
	}

	token, err := requestMSToken(requester, MSName)
	if err != nil {
		return err
	}

	requester.Token.Set(MSName, token)

	request.Header.Set("Authorization", token)
	request.Header.Set("Context", conf.Context)
	return nil
}

func requestMSToken(requester *RequesterObj, MSName string) (string, error) {
	conf, _ := requester.Confs.Get(MSName)

	authPayload := map[string]string{
		"token":  conf.Token,
		"secret": conf.Secret,
	}
	body, _ := json.Marshal(authPayload)

	req, _ := http.NewRequest(
		"POST",
		conf.BaseUrl+"/auth/generate",
		bytes.NewBuffer(body),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := requester.Client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	token := resp.Header.Get("X-Token")

	if len(token) == 0 || resp.StatusCode != 204 {
		return "", fmt.Errorf("error to authenticate, empty token or invalid status code")
	}

	return token, err
}
