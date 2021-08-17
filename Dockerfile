FROM ubuntu:18.04

# Install basic programs
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends vim.tiny wget curl sudo tar git-core make ca-certificates gcc

# Install and set up Go
RUN wget https://golang.org/dl/go1.15.10.linux-amd64.tar.gz --no-check-certificate
RUN tar xf go1.15.10.linux-amd64.tar.gz
RUN mv go /usr/local/go
RUN echo "export GOROOT=/usr/local/go" >> ~/.bashrc
RUN echo "export PATH=\$GOROOT/bin:\$PATH" >> ~/.bashrc

# Create a directory to mount our source code.
# Local changes will be reflected in the running container.
RUN mkdir /newrelic-cli
RUN sudo chown -R $USER:$USER /newrelic-cli

WORKDIR /newrelic-cli
