FROM golang:1.13.0-alpine3.10 as builder

# Prepare buildement
RUN apk add --update --no-cache git \
   && mkdir -p /output

RUN go get -v gopkg.in/yaml.v2
RUN go get -v github.com/hashicorp/vault/api

# Build
WORKDIR /go/src/vault-backup
COPY ./code/. .
RUN go build -i -o /output/vault-backup .

# Lets create toolbox for use it
FROM alpine:3.10

RUN apk add --no-cache \
    ca-certificates    \
    bash               \
    xz                 \
    gzip               \
    gnupg              \
    python3

RUN pip3 install awscli

COPY --from=builder /output/. /

ENTRYPOINT ["/bin/bash"]
