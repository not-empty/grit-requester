# grit-requester

[![Test Coverage](https://img.shields.io/badge/test-100%25-brightgreen)](#)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

**grit-requester** is a Go library to abstract requests to microservices built using Grit.

Features:

- 🔁 Automatic retry on `401 Unauthorized`
- 🔐 Per-service token cache with concurrency safety
- 💉 Config and HTTP client injection (perfect for testing)
- 📦 Full support for generics (`any`) in request/response
- 🧠 Context-aware: all requests support context.Context for cancellation, timeouts, and APM tracing

---

## ✨ Installation

```bash
go get github.com/not-empty/grit-requester
```

---

## 🚀 Usage Example

```go
import "github.com/not-empty/grit-requester"

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
```

---

## 🧪 Testing

Test coverage: **100%**

Run tests:

```bash
./test.sh
```

Visualize coverage:

```bash
open ./coverage/coverage-unit.html
```

---

## 🧠 Design Overview

- `MsRequest`: generic structure for describing a request
- `ResponseData[T]`: generic expected response wrapper
- `RequesterObj`: manages tokens, configs, and the HTTP client
- `TokenCache`: thread-safe in-memory token cache
- `DoMsRequest`: core function to execute the request

---

## 🔧 License

MIT © [Not Empty](https://github.com/not-empty)

**Not Empty Foundation - Free codes, full minds**
