package gritrequester

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestNewRequestObj(t *testing.T) {
	conf := StaticConfig{
		"auth": {
			Token:   "abc",
			Secret:  "xyz",
			Context: "ctx",
			BaseUrl: "http://localhost",
		},
	}

	r := NewRequestObj(conf)

	if r == nil {
		t.Fatal("expected non-nil RequesterObj")
	}

	if r.Client == nil {
		t.Error("expected non-nil HTTP client")
	}

	if r.Token == nil {
		t.Error("expected token cache to be initialized")
	}

	got, err := r.Confs.Get("auth")
	if err != nil {
		t.Errorf("expected config for 'auth' to exist, got error: %v", err)
	}

	if got.BaseUrl != "http://localhost" {
		t.Errorf("expected BaseUrl to be 'http://localhost', got '%s'", got.BaseUrl)
	}
}

func TestDoMsRequestSuccess(t *testing.T) {
	mockClient := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			var resp string
			var statusCode int

			if req.URL.Path == "/user/create" {
				resp = `{"id":"01JRVBXTGFHF9137A738S7GMFV"}`
				statusCode = 200
			}

			if req.URL.Path == "/auth/generate" {
				resp = ""
				statusCode = 204
			}

			// Simula resposta JSON válida

			return &http.Response{
				StatusCode: statusCode,
				Body:       io.NopCloser(strings.NewReader(resp)),
				Header:     http.Header{"X-Token": []string{"mock-token"}},
			}, nil
		},
	}

	conf := StaticConfig{
		"mock": MSAuthConf{
			Token:   "token",
			Secret:  "secret",
			Context: "ctx",
			BaseUrl: "http://fake.com",
		},
	}

	r := &RequesterObj{
		Client: mockClient,
		Token:  NewTokenCache(),
		Confs:  conf,
	}

	type Output struct {
		ID string `json:"id"`
	}

	msReq := MsRequest{
		MSName: "mock",
		Method: "POST",
		Path:   "/user/create",
		Body:   map[string]string{"email": "test@test.com"},
	}

	resp, err := DoMsRequest[Output](r, msReq, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "01JRVBXTGFHF9137A738S7GMFV" {
		t.Errorf("expected id 01JRVBXTGFHF9137A738S7GMFV, got %s", resp.ID)
	}
}

var countRetryUnauthorized int = 0

func TestDoMsRequestRetryUnauthorized(t *testing.T) {
	mockClient := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			var resp string
			var statusCode int

			if req.URL.Path == "/user/create" {
				if countRetryUnauthorized == 0 {
					resp = "Unathorized"
					statusCode = 401
					countRetryUnauthorized++
				} else {
					resp = `{"id":"01JRVBXTGFHF9137A738S7GMFV"}`
					statusCode = 200
				}
			}

			if req.URL.Path == "/auth/generate" {
				resp = ""
				statusCode = 204
			}

			return &http.Response{
				StatusCode: statusCode,
				Body:       io.NopCloser(strings.NewReader(resp)),
				Header:     http.Header{"X-Token": []string{"mock-token"}},
			}, nil
		},
	}

	conf := StaticConfig{
		"mock": MSAuthConf{
			Token:   "token",
			Secret:  "secret",
			Context: "ctx",
			BaseUrl: "http://fake.com",
		},
	}

	r := &RequesterObj{
		Client: mockClient,
		Token:  NewTokenCache(),
		Confs:  conf,
	}

	type Output struct {
		ID string `json:"id"`
	}

	msReq := MsRequest{
		MSName: "mock",
		Method: "POST",
		Path:   "/user/create",
		Body:   map[string]string{"email": "test@test.com"},
	}

	resp, err := DoMsRequest[Output](r, msReq, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.ID != "01JRVBXTGFHF9137A738S7GMFV" {
		t.Errorf("expected id 01JRVBXTGFHF9137A738S7GMFV, got %s", resp.ID)
	}

	if countRetryUnauthorized != 1 {
		t.Errorf("Request was not executed twice")
	}
}

func TestNewRequestReturnError(t *testing.T) {
	mockClient := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     http.Header{"X-Token": []string{"mock-token"}},
			}, nil
		},
	}

	conf := StaticConfig{
		"mock": MSAuthConf{
			Token:   "token",
			Secret:  "secret",
			Context: "ctx",
			BaseUrl: "://invalid-url",
		},
	}

	r := &RequesterObj{
		Client: mockClient,
		Token:  NewTokenCache(),
		Confs:  conf,
	}

	msReq := MsRequest{
		MSName: "mock",
		Method: "POST",
		Path:   "/user/create",
		Body:   map[string]string{"email": "test@test.com"},
	}

	_, err := newRequest(r, msReq)

	if err == nil {
		t.Error("Expected error to create request object but got nil")
	}
}

func TestNewRequestErrorToSetHeaders(t *testing.T) {
	mockClient := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("Error to execute request")
		},
	}

	conf := StaticConfig{
		"mock": MSAuthConf{
			Token:   "token",
			Secret:  "secret",
			Context: "ctx",
			BaseUrl: "http://fake-url.com",
		},
	}

	r := &RequesterObj{
		Client: mockClient,
		Token:  NewTokenCache(),
		Confs:  conf,
	}

	msReq := MsRequest{
		MSName: "mock",
		Method: "POST",
		Path:   "/user/create",
		Body:   map[string]string{"email": "test@test.com"},
	}

	_, err := newRequest(r, msReq)

	if err == nil {
		t.Error("Expected error to create request object but got nil")
	}
}

func TestDoMsRequestWithoutConfs(t *testing.T) {
	mockClient := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			var resp string = "test"
			var statusCode int = 200
			return &http.Response{
				StatusCode: statusCode,
				Body:       io.NopCloser(strings.NewReader(resp)),
				Header:     http.Header{"X-Token": []string{"mock-token"}},
			}, nil
		},
	}

	conf := StaticConfig{}

	r := &RequesterObj{
		Client: mockClient,
		Token:  NewTokenCache(),
		Confs:  conf,
	}

	msReq := MsRequest{
		MSName: "mock",
		Method: "GET",
		Path:   "/user/list",
		Body:   nil,
	}

	_, err := DoMsRequest[any](r, msReq, true)

	if err == nil || err.Error() != "config map is empty" {
		t.Errorf("expeted an error but but got none")
	}
}

func TestDoMsRequestWithInvalidBody(t *testing.T) {
	mockClient := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			var resp string = "test"
			var statusCode int = 200
			return &http.Response{
				StatusCode: statusCode,
				Body:       io.NopCloser(strings.NewReader(resp)),
				Header:     http.Header{"X-Token": []string{"mock-token"}},
			}, nil
		},
	}

	conf := StaticConfig{
		"mock": MSAuthConf{
			Token:   "token",
			Secret:  "secret",
			Context: "ctx",
			BaseUrl: "http://fake.com",
		},
	}

	r := &RequesterObj{
		Client: mockClient,
		Token:  NewTokenCache(),
		Confs:  conf,
	}

	type Invalid struct {
		Fn func()
	}

	v := Invalid{
		Fn: func() {},
	}

	msReq := MsRequest{
		MSName: "mock",
		Method: "POST",
		Path:   "/user/list",
		Body:   v,
	}

	_, err := DoMsRequest[any](r, msReq, true)

	if err == nil || err.Error() != "json: unsupported type: func()" {
		t.Errorf("expeted an error but but got none")
	}
}

func TestSetRequestHeadersGetFromCache(t *testing.T) {
	mockClient := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, nil
		},
	}

	conf := StaticConfig{"mock": MSAuthConf{
		Token:   "token",
		Secret:  "secret",
		Context: "ctx",
		BaseUrl: "http://fake.com",
	}}

	r := &RequesterObj{
		Client: mockClient,
		Token:  NewTokenCache(),
		Confs:  conf,
	}

	r.Token.Set("mock", "token test")

	request, _ := http.NewRequest(
		"POST",
		"http://fake-url.com",
		strings.NewReader(`{"user":"test"}`),
	)

	err := setRequestHeaders(r, request, "mock")

	if err != nil {
		t.Errorf("Error to set request headers")
	}
}

func TestRequestMSTokenRetuningInvalidResponse(t *testing.T) {
	mockClient := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			var resp string = ""
			var statusCode int = 401
			return &http.Response{
				StatusCode: statusCode,
				Body:       io.NopCloser(strings.NewReader(resp)),
			}, nil
		},
	}

	conf := StaticConfig{
		"mock": MSAuthConf{
			Token:   "token",
			Secret:  "secret",
			Context: "ctx",
			BaseUrl: "http://fake.com",
		},
	}

	r := &RequesterObj{
		Client: mockClient,
		Token:  NewTokenCache(),
		Confs:  conf,
	}

	token, err := requestMSToken(
		r,
		"mock",
	)

	if len(token) > 0 ||
		err != nil &&
			err.Error() != "error to authenticate, empty token or invalid status code" {
		t.Errorf("Expected to get an error but got nil")
	}
}

func TestExecRequestFailRequest(t *testing.T) {
	mockClient := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("Error to execute request")
		},
	}

	conf := StaticConfig{"mock": MSAuthConf{
		Token:   "token",
		Secret:  "secret",
		Context: "ctx",
		BaseUrl: "http://fake.com",
	}}

	r := &RequesterObj{
		Client: mockClient,
		Token:  NewTokenCache(),
		Confs:  conf,
	}

	r.Token.Set("mock", "token test")

	request, _ := http.NewRequest(
		"POST",
		"http://fake-url.com",
		strings.NewReader(`{"user":"test"}`),
	)

	status, _, err := execRequest[any](r, "mock", request)

	if status != 0 || err == nil {
		t.Errorf("Expect to fail request but it was suscceful")
	}
}

type ErrorReader struct{}

func (e ErrorReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("erro simulado ao ler body")
}

func (e ErrorReader) Close() error {
	return nil
}

func TestExecRequestErrorInReadAll(t *testing.T) {
	mockClient := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			var statusCode int = 200
			return &http.Response{
				StatusCode: statusCode,
				Body:       ErrorReader{},
			}, nil
		},
	}

	conf := StaticConfig{"mock": MSAuthConf{
		Token:   "token",
		Secret:  "secret",
		Context: "ctx",
		BaseUrl: "http://fake.com",
	}}

	r := &RequesterObj{
		Client: mockClient,
		Token:  NewTokenCache(),
		Confs:  conf,
	}

	r.Token.Set("mock", "token test")

	request, _ := http.NewRequest(
		"POST",
		"http://fake-url.com",
		strings.NewReader(`{"user":"test"}`),
	)

	_, _, err := execRequest[any](r, "mock", request)

	if err == nil || err.Error() != "erro simulado ao ler body" {
		t.Errorf("Expect to fail request but it was suscceful")
	}
}

func TestExecRequestResponseGreaterThan299(t *testing.T) {
	mockClient := &MockClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			var resp string = `{"test":"test"}`
			var statusCode int = 301
			return &http.Response{
				StatusCode: statusCode,
				Body:       io.NopCloser(strings.NewReader(resp)),
			}, nil
		},
	}

	conf := StaticConfig{"mock": MSAuthConf{
		Token:   "token",
		Secret:  "secret",
		Context: "ctx",
		BaseUrl: "http://fake.com",
	}}

	r := &RequesterObj{
		Client: mockClient,
		Token:  NewTokenCache(),
		Confs:  conf,
	}

	r.Token.Set("mock", "token test")

	request, _ := http.NewRequest(
		"POST",
		"http://fake-url.com",
		strings.NewReader(`{"user":"test"}`),
	)

	_, _, err := execRequest[any](r, "mock", request)

	if err == nil {
		t.Errorf("Expect to fail request but it was suscceful")
	}
}

func TestUpdateServiceToken_NoTokenInCache(t *testing.T) {
	tokenCache := NewTokenCache()

	r := &RequesterObj{
		Token: tokenCache,
	}

	headers := http.Header{}
	headers.Set("X-Token", "new-token")

	updateServiceToken(r, "auth", headers)

	_, exists := tokenCache.Get("auth")
	if exists {
		t.Error("expected token to not be set because it wasn't in cache initially")
	}
}

func TestUpdateServiceToken_HeaderEmpty(t *testing.T) {
	tokenCache := NewTokenCache()
	tokenCache.Set("auth", "old-token")

	r := &RequesterObj{
		Token: tokenCache,
	}

	headers := http.Header{}
	headers.Set("X-Token", "") // vazio

	updateServiceToken(r, "auth", headers)

	token, _ := tokenCache.Get("auth")
	if token != "old-token" {
		t.Errorf("expected token to remain unchanged, got '%s'", token)
	}
}

func TestUpdateServiceToken_UpdatesIfDifferent(t *testing.T) {
	tokenCache := NewTokenCache()
	tokenCache.Set("auth", "old-token")

	r := &RequesterObj{
		Token: tokenCache,
	}

	headers := http.Header{}
	headers.Set("X-Token", "new-token")

	updateServiceToken(r, "auth", headers)

	token, _ := tokenCache.Get("auth")
	if token != "new-token" {
		t.Errorf("expected token to be updated to 'new-token', got '%s'", token)
	}
}
