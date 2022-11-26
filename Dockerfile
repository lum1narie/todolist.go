FROM golang:1.17 as builder

ENV GOPATH=/go/src
ENV WORKSPACE=${GOPATH}/app
RUN mkdir -p ${WORKSPACE}

WORKDIR ${WORKSPACE}

COPY . ${WORKSPACE}

RUN go mod download && \
    go mod tidy -compat=1.17

RUN go build main.go

FROM alpine:3.17 as runner

RUN apk --no-cache add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    mkdir /lib64 && \
    ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

WORKDIR /app

COPY --from=builder /go/src/app/main ./

ENV PATH $PATH:/app

ENTRYPOINT ["/bin/sleep", "INF"]
