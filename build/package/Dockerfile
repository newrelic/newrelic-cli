FROM alpine:3.24.1

# Refresh installed packages so we pick up any libcrypto3/libssl3 (and other)
# security fixes published to the v3.24 repo after this base image was tagged.
# Without this, scanners flag CVEs whose fixes are in the live repo but not
# yet baked into a new alpine point release.
RUN apk upgrade --no-cache

# Add the binary
COPY newrelic /bin/newrelic

ENTRYPOINT ["/bin/newrelic"]
