FROM golang:1.11.5-alpine as builder

WORKDIR /go/src

ENV VERSION="1.1.4"

RUN apk --update add alpine-sdk && \
    curl -L -o- https://github.com/adfinis-sygroup/vault-client/archive/v${VERSION}.tar.gz  | tar xvzf - && \
    cd vault-client-${VERSION} && \
    make build  && mkdir /go/build && \
    mv /go/src/vault-client-${VERSION}/vc /go/build/vc

FROM alpine

COPY --from=builder /go/build/vc /vc

CMD ["/vc"]


