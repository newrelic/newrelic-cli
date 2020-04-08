#!/bin/bash

curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/repos/ctrombley/newrelic-cli/releases/latest | jq -r '.upload_url'