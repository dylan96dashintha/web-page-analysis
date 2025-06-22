FROM golang:1.24.2-alpine as builder
WORKDIR /app
COPY . .
RUN go build .

# Final image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app .
EXPOSE 8080
CMD ["./web-page-analysis"]