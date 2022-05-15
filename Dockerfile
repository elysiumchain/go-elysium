# Support setting various labels on the final image
ARG COMMIT=""
ARG VERSION=""
ARG BUILDNUM=""

# Build Gely in a stock Go builder container
FROM golang:1.18-alpine as builder

RUN apk add --no-cache gcc musl-dev linux-headers git

ADD . /go-elysium
RUN cd /go-elysium && go run build/ci.go install ./cmd/gely

# Pull Gely into a second stage deploy alpine container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go-elysium/build/bin/gely /usr/local/bin/

EXPOSE 2022 2023 20202 20202/udp
ENTRYPOINT ["gely"]

# Add some metadata labels to help programatic image consumption
ARG COMMIT=""
ARG VERSION=""
ARG BUILDNUM=""

LABEL commit="$COMMIT" version="$VERSION" buildnum="$BUILDNUM"
