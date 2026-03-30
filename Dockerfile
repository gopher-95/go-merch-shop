FROM golang:1.25

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN ls -R /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /my_app ./cmd/api

CMD ["/my_app"]