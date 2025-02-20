FROM golang:1.23

RUN go install github.com/air-verse/air@latest
RUN alias air='~/.air'

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates nodejs npm \
    && update-ca-certificates


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air", "-c", ".air.toml"]