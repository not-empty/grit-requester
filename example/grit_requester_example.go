package main

import (
	"context"
	"fmt"
	"log"

	gritrequester "github.com/not-empty/grit-requester"
)

func main() {
	type RequestPayload struct {
		Example string `json:"example"`
	}

	type ResponseData struct {
		ID string `json:"id"`
	}

	conf := gritrequester.StaticConfig{}

	conf.Set("example", gritrequester.MSAuthConf{
		Token:   "your-integration-token",
		Secret:  "your-integration-secret",
		Context: "example-test",
		BaseUrl: "https://example.microservice.local",
	})

	client := gritrequester.NewRequestObj(conf)

	msReq := gritrequester.MsRequest{
		MSName: "example",
		Method: "POST",
		Path:   "/example/add",
		Body: RequestPayload{
			Example: "Example test",
		},
	}

	resp, err := gritrequester.DoMsRequest[ResponseData](context.TODO(), client, msReq, true)
	if err != nil {
		log.Fatal("Request failed:", err)
	}

	fmt.Println("Received ID:", resp.Data.ID)
}
