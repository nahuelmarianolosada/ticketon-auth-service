FROM golang:1.19-alpine as builder
RUN mkdir /build
WORKDIR /build
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /ticketon-auth-service .

FROM alpine
RUN apk add libc6-compat
COPY --from=builder /build/ticketon-auth-service /ticketon-auth-service
COPY --from=builder /build/.env.example /.env
ENTRYPOINT ["./ticketon-auth-service"]