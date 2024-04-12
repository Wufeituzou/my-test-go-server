# Step1. modules caching
FROM golang:1.22.1 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step2. Builder
FROM golang:1.22.1  as builder
COPY --from=modules /go/pkg /go/pkg
# 将项目代码复制到工作目录
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/app ./cmd/app

# Step 3: Final
FROM scratch
# COPY --from=builder /app/config /config
# COPY --from=builder /app/migrations /migrations
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["/app"]