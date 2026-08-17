FROM golang:1.24.3-alpine

WORKDIR /app

RUN apk add --no-cache git

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./

ARG GITHUB_TOKEN
RUN git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"


RUN go mod download

COPY . .

EXPOSE 3004

CMD ["air", "-c", ".air.toml"]

