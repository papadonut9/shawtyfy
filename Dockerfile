# syntax=docker/dockerfile:1

FROM golang:1.21-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o shawtyfy ./main.go

FROM gcr.io/distroless/static
COPY --from=build /src/shawtyfy /shawtyfy
EXPOSE 9808
ENTRYPOINT ["/shawtyfy"]
