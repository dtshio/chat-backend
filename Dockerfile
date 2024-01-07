FROM golang:1.21.3

WORKDIR /app
COPY go.mod go.sum ./
COPY cmd ./cmd
COPY pkg ./pkg
RUN go mod download
RUN GOOS=linux go build -o /api ./cmd
EXPOSE 3333
CMD ["/api"]
