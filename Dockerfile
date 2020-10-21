FROM --platform=$TARGETPLATFORM golang:1.13.5-stretch as devel
COPY / /go/src/
RUN cd /go/src/ && make all

FROM --platform=$TARGETPLATFORM busybox
COPY /scripts/native/templates /etc/baetyl/templates
COPY --from=devel /go/src/output/baetyl-cloud /bin/
ENTRYPOINT ["baetyl-cloud"]