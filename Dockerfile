FROM golang:alpine AS builder

RUN apk add --no-cache git
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./out/cup-predictor-server

FROM alpine:latest

COPY --from=builder /app/out/cup-predictor-server /app/cup-predictor-server
COPY --from=builder /app/assets /assets

EXPOSE 7000

CMD [ "/app/cup-predictor-server" ]