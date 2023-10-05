FROM golang:1.21-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -a -ldflags "-w -s" -o ftpcam ./cmd/ftpcam/main.go

FROM scratch
COPY --from=builder /etc/ssl/cert.pem /etc/ssl/
COPY --from=builder /app/ftpcam /ftpcam
ENTRYPOINT ["/ftpcam"]
