#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v ./...

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
ENTRYPOINT /app
LABEL Name=go-cleanapi Version=0.0.1
EXPOSE 8080


# FROM golang:latest
# WORKDIR /app

# # manage dependencies
# COPY go.mod .
# COPY go.sum .
# RUN go mod download

# COPY . .

# ENV PORT 8080

# RUN go build

# CMD [ "./pssword-api-mysql" ]