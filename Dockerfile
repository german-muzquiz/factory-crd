FROM golang:1.18.1 as builder

ENV GOOS=linux GOARCH=amd64
WORKDIR /opt/factory-control/build/
ADD ./ /opt/factory-control/build/
RUN make build

FROM debian:stable-slim
COPY --from=builder "/opt/factory-control/build/bin/linux_amd64/factory-control" "/usr/local/bin/factory-control"

ENTRYPOINT ["/usr/local/bin/factory-control"]
