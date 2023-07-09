FROM golang:1.20

COPY . /app

WORKDIR /app

RUN go mod tidy && go build -o versionup cmd/cli/main.go

CMD ["./versionup"]