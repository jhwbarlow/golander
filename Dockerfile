FROM golang:1.14-alpine AS builder
COPY . /tmp/src/
RUN cd /tmp/src/cmd \
  && go build -o /tmp/bin/golander main.go \
  && chmod 500 /tmp/bin/golander

FROM scratch
COPY docker/passwd docker/group /etc/
USER golander:golander
COPY --from=builder --chown=golander:golander /tmp/bin/golander /usr/bin/golander
ENTRYPOINT [ "/usr/bin/golander" ]