FROM golang:1.24.3-alpine

WORKDIR /app

#ENV PATH="/go/bin:$PATH"
#RUN go install github.com/cosmtrek/air@v1.51.0

RUN apk add --no-cache git


RUN go install github.com/air-verse/air@latest

#RUN #mkdir -p /app/tmp && chmod -R 777 /app/tmp/


COPY go.mod go.sum ./

ARG GITHUB_TOKEN
RUN git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"


RUN go mod download

COPY . .

#RUN #go build -o main .

EXPOSE 3002

#CMD ["/app/start.sh"]

#CMD ["ls", "-la", "./backend/galactu?s"]
CMD ["air", "-c", ".air.toml"]
#CMD ["./main"]