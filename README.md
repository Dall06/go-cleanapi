# go-cleanapi

go test ./... -coverprofile=coverage.out -coverpkg=./... && go tool cover -func=coverage.out
export PATH=$PATH:$(go env GOPATH)/bin
golangci-lint run  
