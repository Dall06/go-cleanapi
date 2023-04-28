FROM golang:latest
WORKDIR /app

# manage dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENV PORT 8080

RUN go build

CMD [ "./go-cleanapi" ]