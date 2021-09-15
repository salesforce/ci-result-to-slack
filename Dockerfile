FROM golang:1.17-bullseye as gobuild

WORKDIR /go/src/app/

COPY . .
RUN make ci

FROM gcr.io/distroless/base-debian11

COPY --from=gobuild /go/src/app/bin/ci-result-to-slack /ci-result-to-slack

ENTRYPOINT ["/ci-result-to-slack"]
