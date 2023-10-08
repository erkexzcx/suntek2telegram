FROM --platform=$BUILDPLATFORM golang:1.21-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
ARG version
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GOARM=${TARGETVARIANT#v} go build -a -ldflags "-w -s -X main.version=$version" -o suntek2telegram ./cmd/ftpcam/main.go

FROM scratch
COPY --from=builder /etc/ssl/cert.pem /etc/ssl/
COPY --from=builder /app/suntek2telegram /suntek2telegram
ENTRYPOINT ["/suntek2telegram"]
