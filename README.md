# go-cleanapi

This is a REST API template project, based on Uncle's Bob Clean Architecture principles, using go fiber as rest server package alongside another complementary packages for middleware and utilities.

The project contain different packages based on the layers of the consulted literature.

Principal unit tests can be found, as well as CI implementation with docker and GitHub actions.

## Run

Use the go command [go](https://go.dev/) to run the project locally.

```bash
go run main.go -p <port_flag_value> -v -p <version_flag_value>
```

## Principal commands

```bash
# test with coverage and save in file
go test ./... -coverprofile=coverage.out -coverpkg=./... && go tool cover -func=coverage.out

# export go path env
export PATH=$PATH:$(go env GOPATH)/bin

# runs golangci linter
golangci-lint run

# runs swagger for the endpoints
swag init
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)
