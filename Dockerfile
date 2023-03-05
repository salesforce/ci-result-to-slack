# syntax=docker/dockerfile:1.4
# Use Wolfi/Chainguard Images
# https://www.chainguard.dev/unchained/introducing-wolfi-the-first-linux-un-distro
# https://github.com/chainguard-images/images#chainguard-images
FROM cgr.dev/chainguard/static:latest

COPY ./ci-result-to-slack /ci-result-to-slack

ENTRYPOINT ["/ci-result-to-slack"]
