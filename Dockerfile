# Ingester service image: a static pure-Go binary built reproducibly with
# -trimpath. Both bases are digest-pinned; re-resolve the digests when
# bumping either base (docker buildx imagetools inspect <ref>).
#
# Base digests resolved 2026-10-28.
FROM golang:1.26.10-alpine@sha256:a2012b4a1d1f9f7e9f44e0d39b805a901e856b1dd80b4e6ee4338db547050308 AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/ingester ./cmd/ingester

FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab
COPY --from=build /out/ingester /usr/local/bin/ingester
USER 65532:65532
ENTRYPOINT ["/usr/local/bin/ingester"]
