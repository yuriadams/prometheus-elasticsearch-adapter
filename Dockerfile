FROM golang:1.9
ADD . /go/src/github.com/yuriadams/prometheus-elasticsearch-adapter
ENV CGO_ENABLED=0
RUN go install github.com/yuriadams/prometheus-elasticsearch-adapter

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/bin/prometheus-elasticsearch-adapter /usr/local/bin/prometheus-elasticsearch-adapter
ENTRYPOINT ["/usr/local/bin/prometheus-elasticsearch-adapter"]
