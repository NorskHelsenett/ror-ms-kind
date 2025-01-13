ARG GCR_MIRROR=gcr.io/
FROM ${GCR_MIRROR}distroless/static:nonroot
LABEL org.opencontainers.image.source https://github.com/norskhelsenett/ror
WORKDIR /

COPY dist/ror-ms-kind /bin/ror-ms-kind
ENTRYPOINT ["/bin/ror-ms-kind"]
