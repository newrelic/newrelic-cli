FROM alpine:3.11

# Default configs live in .newrelic
RUN mkdir -p /root/.newrelic/plugins
COPY configs/ /root/.newrelic

# Add the binary
COPY ./bin/linux/newrelic /bin

CMD ["/bin/newrelic"]
