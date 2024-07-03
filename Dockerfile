FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:3.16

RUN apk --no-cache add ca-certificates
RUN apk update && \
    apk add --no-cache imagemagick pngquant curl ffmpeg libreoffice librsvg potrace

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env .env
RUN mkdir -p /app/assets/documents /app/assets/images /app/assets/videos

EXPOSE 3000

CMD ["./main"]
