FROM golang:1.22-alpine3.19

WORKDIR /src/chat-app-backend-server
ADD . /src/chat-app-backend-server

ENV CGO_ENABLED=0

RUN go mod vendor
RUN cd ./cmd/ && go build -o ../app

EXPOSE 5000

ENV WAIT_VERSION 2.7.2
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

CMD ["/src/chat-app-backend-server/app"]
