# build stage
FROM golang:alpine AS build-env
WORKDIR /go/src/github.com/gdunghi/go_hook
COPY . /go/src/github.com/gdunghi/go_hook
RUN apk add --no-cache git
RUN apk add --no-cache bash
RUN go get github.com/labstack/echo
RUN go get github.com/utahta/go-linenotify

RUN go build -o main .

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/github.com/gdunghi/go_hook/main .
RUN apk update
RUN apk add --no-cache openssh-keygen
COPY id_rsa /root/.ssh/id_rsa
COPY id_rsa.pub /root/.ssh/id_rsa.pub

RUN touch info.log

ENTRYPOINT ["./main"]
