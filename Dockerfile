FROM golang:1.13 AS build
WORKDIR /work
COPY . .
RUN go get ./...
RUN CGO_ENABLED=0 go build -ldflags="-w -s -X \"github.com/florinutz/filme/pkg.Version=built-by-docker\"" -mod=readonly -v -o filme

FROM alpine:3
RUN apk --no-cache add ca-certificates
COPY --from=build /work/filme /
CMD ["/filme",  "serve"]
