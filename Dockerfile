FROM alpine:3.11

# Add the binary
COPY ./bin/linux/newrelic /bin

CMD ["/bin/newrelic"]
