FROM golang:1.11 as builder
COPY . /work
COPY ./pkg/filme/tpl/google /work/templates/google
WORKDIR /work
RUN useradd app
RUN CGO_ENABLED=0 go build -ldflags="-w -s -X \"github.com/florinutz/filme/pkg.Version=built-by-docker\"" -o filme

FROM scratch
LABEL maintainer="florinutz@gmail.com"
EXPOSE 14051/tcp
COPY --from=builder /work/filme /usr/bin/filme
COPY --from=builder /etc/passwd /etc/
COPY --from=builder /work/templates/google /tpl/google
ENV FILME_TPL_GOOGLE /tpl/google/*
USER app
ENTRYPOINT ["/usr/bin/filme"]
