FROM gcr.io/distroless/base-debian11

COPY ./ci-result-to-slack /ci-result-to-slack

ENTRYPOINT ["/ci-result-to-slack"]
