# builder image
FROM golang:alpine AS builder
WORKDIR /build
COPY . /build
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o main /build/main.go


FROM alpine as proxy_pool
COPY --from=builder /build/main  /usr/bin/main
