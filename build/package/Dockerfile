FROM alpine:3.24.1

RUN apk upgrade --no-cache

# Add the binary
COPY newrelic /bin/newrelic

ENTRYPOINT ["/bin/newrelic"]
