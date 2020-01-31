FROM golang:1.13 AS builder
WORKDIR /work
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 go build -ldflags="-w -s -X \"github.com/florinutz/filme/pkg.Version=built-by-docker\"" -mod=readonly -v -o filme

FROM gcr.io/distroless/base
WORKDIR /
COPY --from=builder /work/filme /
ENTRYPOINT ["/filme",  "serve"]
ENV PORT=14051
EXPOSE 14051
