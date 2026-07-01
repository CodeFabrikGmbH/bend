FROM golang:1.23-alpine AS build

WORKDIR /workdir

RUN apk add --no-cache ca-certificates git

COPY . /workdir

RUN CGO_ENABLED=0 go build -o /app

FROM alpine:3.9
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app /usr/local/bin/app

# resources, static assets and README are embedded into the binary (embed.FS),
# so they no longer need to be copied into the runtime image.

ENTRYPOINT ["/usr/local/bin/app"]
