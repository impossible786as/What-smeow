FROM golang:1.25-alpine

WORKDIR /app

RUN apk add --no-cache git

RUN go mod init whatsmeow-full-extractor

# 🚨 پہلے کوڈ کو کاپی کریں تاکہ Go کو imports مل سکیں
COPY main.go .

# 🚨 پھر whatsmeow گیٹ اور ٹائیڈی کریں
RUN go get -u go.mau.fi/whatsmeow@latest
RUN go mod tidy

ENV PORT=8080
EXPOSE 8080

CMD ["go", "run", "main.go"]
