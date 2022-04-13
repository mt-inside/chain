FROM docker.io/library/golang:1.18 AS build

# no Debian packages
RUN curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to /usr/local/bin

ADD . /go/src/github.com/mt-inside/chain
WORKDIR /go/src/github.com/mt-inside/chain

RUN just install-static


FROM gcr.io/distroless/static:latest

COPY --from=build /go/bin/chain /
CMD ["/chain"]
