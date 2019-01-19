# Build
FROM golang:alpine AS build

RUN apk add --no-cache -U git make build-base

WORKDIR /src/autodock
COPY . /src/autodock
RUN make install

# Runtime
FROM alpine:latest

COPY --from=build /go/bin/autodock /autodock

EXPOSE 8000/tcp

ENTRYPOINT ["/autodock"]
CMD []
