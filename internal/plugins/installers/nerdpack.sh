#! /bin/bash

wget https://github.com/ctrombley/newrelic-cli-plugin-nerdpack/archive/master.zip -O temp.zip && \ 
    mkdir -p ~/newrelic/plugins/nerdpack && \
    unzip -oj  temp.zip -d ~/.newrelic/plugins/nerdpack && \
    rm temp.zip && \
    npm --prefix $HOME/.newrelic/plugins/nerdpack install $HOME/.newrelic/plugins/nerdpack