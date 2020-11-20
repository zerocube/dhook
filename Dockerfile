FROM golang:1.15 as builder


WORKDIR /go/src/github.com/zerocube/dhook

COPY go.mod go.sum ./

RUN go mod download -x

COPY . ./

ENV CGO_ENABLED=0
ARG WEBHOOK_URL

RUN go test -v ./... && go build -ldflags "-X main.webhookURL=${WEBHOOK_URL}" -o /dhook -v .

ENTRYPOINT ["go", "run", "."]

FROM scratch

# Copy across the CA certificate bundle
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.certs

# Copy across the compiled binary
COPY --from=builder /dhook ./

ENTRYPOINT ["./dhook"]
