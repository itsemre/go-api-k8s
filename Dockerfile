ARG GO_VERSION=1.20
FROM golang:${GO_VERSION}-alpine AS build
ARG GOOS=linux
ARG GOARCH=amd64
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ENV CGO_ENABLED=0 \
    GOOS=${GOOS} \
    GOARCH=${GOARCH}
RUN go build -ldflags "-s -w" -a -installsuffix cgo -o api .

FROM alpine:3.18.0
ENV GIN_MODE="release"
WORKDIR /home/app
COPY --from=build /app/api /home/app/api
ENTRYPOINT ["/home/app/api"]