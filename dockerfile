FROM golang:alpine AS builder
MAINTAINER wumiao
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct
WORKDIR /build
COPY . .
RUN go build -o app .

FROM scratch
COPY --from=builder /build/app /app/app
ENTRYPOINT ["/app/app"]
