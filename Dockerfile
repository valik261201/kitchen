From golang:1.19

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN go build main

EXPOSE 8060