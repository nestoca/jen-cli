FROM golang:1.14-buster as build

ENV GO111MODULE="on"
WORKDIR /app
COPY src/go.mod .
COPY src/go.sum .
RUN go mod download
COPY src .
RUN go build

# -----------

FROM debian:13.4

# Explicitly install libcap2 to pick up fix for SNYK-DEBIAN13-LIBCAP2-15960217 from trixie-security (PL-5740).
RUN apt-get update && apt-get install -y --no-install-recommends libcap2 && rm -rf /var/lib/apt/lists/*

COPY --from=build /app/jen /jen
CMD ["bash"]
