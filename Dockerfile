FROM golang:1.23-alpine AS builder
WORKDIR /gau_truyen
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -tags 'prod' -o main .

FROM alpine:latest
WORKDIR /gau_truyen
COPY --from=builder /gau_truyen/main .
EXPOSE 8084
CMD ["./main"]