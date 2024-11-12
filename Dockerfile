FROM golang:1.23.3-bullseye

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN go build -o /httpserver

EXPOSE 8080

CMD [ "/httpserver" ]