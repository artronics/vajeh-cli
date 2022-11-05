FROM golang:alpine as go_build
ENV GO111MODULE=on

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN CGO_ENABLED=0 go build -o app .

FROM alpine:latest
WORKDIR /root/
COPY --from=go_build /app ./

ENTRYPOINT ["./app"]
