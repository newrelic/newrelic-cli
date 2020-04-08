#!/bin/bash

curl -H "Authorization: token $GITHUB_TOKEN" https://api.github.com/repos/newrelic/newrelic-cli/releases/latest | jq -r '.url'