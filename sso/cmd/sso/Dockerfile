FROM golang

WORKDIR /app

#COPY ./sso/go.mod ./sso/go.sum  ./
COPY ./sso .

RUN go mod download

RUN go build -o main ./cmd/sso/main.go

CMD ["./main", "--config=./config/config.yaml"]

EXPOSE 44444