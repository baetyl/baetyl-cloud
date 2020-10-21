FROM golang:1.13.5-stretch as devel
COPY / /go/src/
RUN cd /go/src/ && make build

FROM busybox
COPY /scripts/native/templates /etc/templates
COPY --from=devel /go/src/output/baetyl-cloud /bin/
ENTRYPOINT ["baetyl-cloud"]