FROM golang:latest as builder
WORKDIR /app/
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build github.com/aplulu/aws-ecs-deploy/cmd/deploy

FROM alpine:latest
MAINTAINER Aplulu <aplulu.liv@gmail.com>
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /app/deploy /deploy
ENTRYPOINT ["/deploy"]
