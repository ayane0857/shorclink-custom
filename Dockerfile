# ビルドステージ
FROM golang:1.24-alpine AS builder

WORKDIR /app

# 必要なパッケージのインストール
RUN apk add --no-cache git ca-certificates

# 依存関係のコピーとダウンロード
COPY go.mod go.sum ./
ENV GOPROXY=https://proxy.golang.org,direct
RUN go mod download

# ソースコードのコピー
COPY . .

# アプリケーションのビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shorlink .

# 実行ステージ
FROM alpine:latest

WORKDIR /app

# 必要な証明書のコピー
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# ビルドしたバイナリのコピー
COPY --from=builder /app/shorlink .

# 実行ユーザーの設定
RUN adduser -D -g '' appuser
USER appuser

# アプリケーションの実行
CMD ["./shorlink"]