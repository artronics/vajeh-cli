FROM golang:1.18-alpine as go_build
ENV GO111MODULE=on

WORKDIR /go/src/github.com/artronics/vajeh-cli/
COPY . ./
RUN #CGO_ENABLED=0 go build -a -installsuffix cgo -o app .
#
#FROM alpine:latest
#WORKDIR /root/
#COPY --from=go_build /app ./
#
#CMD ["./app"]
