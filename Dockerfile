FROM golang:1.14-alpine

RUN mkdir /api

COPY . /api

WORKDIR /api

RUN go build -o backend

CMD [ "/api/backend", "-db", "mongodb://mongo", "-addr", ":80" ]

