FROM alpine:latest

WORKDIR /app

# 必要な証明書のコピー
COPY --from=local/shorlink-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# ビルドしたバイナリのコピー
COPY --from=local/shorlink-builder /app/shorlink .

# 実行ユーザーの設定
RUN adduser -D -g '' appuser
USER appuser

# アプリケーションの実行
CMD ["./shorlink"]