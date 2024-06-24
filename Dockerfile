FROM golang:alpine AS builder
ENV CGO_ENABLED 0
ENV GOOS linux
WORKDIR /build
COPY . .
RUN go mod download
RUN cd
RUN cd daserver/cmd && go build -o /build/da-server

FROM alpine
RUN apk update --no-cache
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache tzdata
ENV TZ UTC
WORKDIR /app
COPY --from=builder /build/da-server /app/da-server
ENTRYPOINT ["/app/da-server"]