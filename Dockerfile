FROM golang:latest as builder

LABEL maintainer = "Gustaf Pahlevi"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy && go mod vendor

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# start new stage

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app .

EXPOSE 8081

CMD ["./main"]