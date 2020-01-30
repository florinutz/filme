FROM golang:1.13 AS builder
WORKDIR /work
COPY . ./
RUN go get ./...
RUN ls -al
RUN CGO_ENABLED=0 go build -ldflags="-w -s -X \"github.com/florinutz/filme/pkg.Version=built-by-docker\"" -mod=readonly -v -o filme

FROM gcr.io/distroless/base
EXPOSE 14051
WORKDIR /
COPY --from=builder /work/filme /
ENTRYPOINT ["/filme",  "serve"]
