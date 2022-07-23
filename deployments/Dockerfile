FROM golang:alpine as builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o build/main cmd/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /build .

CMD ["./build/main"]