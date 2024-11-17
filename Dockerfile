FROM golang:1.23.3-bullseye

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN go build -o /http_server

EXPOSE 3333

CMD [ "/http_server" ]