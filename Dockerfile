# Go کا 1.24 ورژن استعمال کریں
FROM golang:1.24-alpine

WORKDIR /app

# Git انسٹال کرنا ضروری ہے تاکہ گو ماڈیولز ڈاؤن لوڈ ہو سکیں
RUN apk add --no-cache git

# ماڈیول انیشلائزیشن
RUN go mod init whatsmeow-full-extractor

# لیٹیسٹ ورژن فیچ کرنا
RUN go get -u go.mau.fi/whatsmeow@latest
RUN go mod tidy

COPY main.go .

# Railway کو بتانے کے لیے کہ کونسی پورٹ ایکسپوز کرنی ہے
ENV PORT=8080
EXPOSE 8080

CMD ["go", "run", "main.go"]
