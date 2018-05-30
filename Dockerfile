FROM golang:alpine
LABEL description="Elastic write/read adapter for Prometheus remote storage."

ENV APP_PATH /go/src/app

RUN mkdir $APP_PATH

WORKDIR $APP_PATH

COPY . .

#Install dep, Git and dependencies
RUN apk --update add git openssh && \
    apk add --update ca-certificates && \
    apk add --no-cache curl && \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && \
    dep ensure && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

CMD go run main.go
