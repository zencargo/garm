FROM docker.io/golang:alpine AS builder
ARG GARM_REF

LABEL stage=builder

RUN apk add musl-dev gcc libtool m4 autoconf g++ make libblkid util-linux-dev git linux-headers upx
RUN git config --global --add safe.directory /build

ADD . /build/garm
RUN cd /build/garm && git checkout ${GARM_REF}
RUN git clone --depth 1 --branch v0.1.2 https://github.com/cloudbase/garm-provider-gcp /build/garm-provider-gcp

RUN cd /build/garm && go build -o /bin/garm \
    -tags osusergo,netgo,sqlite_omit_load_extension \
    -ldflags "-linkmode external -extldflags '-static' -s -w -X github.com/cloudbase/garm/util/appdefaults.Version=$(git describe --tags --match='v[0-9]*' --always --abbrev=0)" \
    /build/garm/cmd/garm && upx /bin/garm
RUN mkdir -p /opt/garm/providers.d
RUN cd /build/garm-provider-gcp && go build -ldflags="-linkmode external -extldflags '-static' -s -w -X main.Version=v0.1.0" -o /opt/garm/providers.d/garm-provider-gcp . && upx /opt/garm/providers.d/garm-provider-gcp

FROM busybox

COPY --from=builder /bin/garm /bin/garm
COPY --from=builder /opt/garm/providers.d/garm-provider-gcp /opt/garm/providers.d/garm-provider-gcp
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/bin/garm", "-config", "/etc/garm/config.toml"]
