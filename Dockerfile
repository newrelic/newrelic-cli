FROM alpine:3.11

# Add the binary
COPY ./bin/linux/newrelic /bin

ENTRYPOINT ["/bin/newrelic"]
