FROM golang:1.26.2-bookworm

ENV GOTOOLCHAIN=local
WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build ./...

EXPOSE 8080
CMD ["go", "run", "./cmd/maintenance-server"]
