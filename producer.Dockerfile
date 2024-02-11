FROM golang:1.21-alpine as builder
WORKDIR /build
COPY go.mod .
RUN go mod download && go mod tidy
COPY . .
RUN go build -o /main main.go

FROM alpine:3
COPY --from=builder main /bin/main
ENTRYPOINT ["/bin/main"]