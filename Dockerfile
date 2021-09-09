FROM golang:1.17-alpine AS builder

ARG GOARCH
ARG GOARM

WORKDIR /go/src/github.com/matrinos/mainflux-agent
COPY . .

RUN apk update \
    && apk add make \
    && make \
    && mv build/mainflux-agent /mainflux-agent

FROM scratch
COPY --from=builder /mainflux-agent /

ENTRYPOINT ["/mainflux-agent"]