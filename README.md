# grit-requester

[![Go Reference](https://pkg.go.dev/badge/github.com/not-empty/grit-requester.svg)](https://pkg.go.dev/github.com/not-empty/grit-requester)
[![Test Coverage](https://img.shields.io/badge/test-100%25-brightgreen)](#)
[![Go Report Card](https://goreportcard.com/badge/github.com/not-empty/grit-requester)](https://goreportcard.com/report/github.com/not-empty/grit-requester)

**grit-requester** is a Go library to abstract requests to microservices built using Grit.

Features:

- 🔁 Automatic retry on `401 Unauthorized`
- 🔐 Per-service token cache with concurrency safety
- 💉 Config and HTTP client injection (perfect for testing)
- 📦 Full support for generics (`any`) in request/response

---

## ✨ Installation

```bash
go get github.com/not-empty/grit-requester
```

---

## 🚀 Usage Example

```go
import "github.com/not-empty/grit-requester"

type LoginInput struct {
	Email string `json:"email"`
	Pass  string `json:"password"`
}

type LoginOutput struct {
	Token string `json:"token"`
	Name  string `json:"name"`
}

conf := gritrequester.StaticConfig{
	"auth": {
		Token:   "your-integration-token",
		Secret:  "your-integration-secret",
		Context: "app-test",
		BaseUrl: "https://auth.microservice.local",
	},
}

client := gritrequester.NewRequestObj(conf)

msReq := gritrequester.MsRequest{
	MSName: "auth",
	Method: "POST",
	Path:   "/auth/login",
	Body: LoginInput{
		Email: "test@example.com",
		Pass:  "123456",
	},
}

resp, err := gritrequester.DoMsRequest[LoginOutput](client, msReq, true)
if err != nil {
	log.Fatal("Request failed:", err)
}

fmt.Println("Received token:", resp.Name)
```

---

## 🧪 Testing

Test coverage: **100%**

Run tests:

```bash
go test -v -cover ./...
```

Visualize coverage:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
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