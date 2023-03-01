FROM golang:1.18.2-bullseye

RUN apt-get update


WORKDIR /source
COPY go.mod .
RUN go mod download
COPY . .


#websocket
EXPOSE 8088  
#grpc
EXPOSE 8000

RUN go build ./cmd/main.go

ENV VALIDATION_URL "http://authenticator/token/challenge/session"
ENV SCHEMA         "thinkshare"
ENTRYPOINT [ "./main" ]