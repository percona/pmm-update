#!/usr/bin/env bash

set -o errexit
set -o pipefail

export VICTORIAMETRICS_HOST=${VICTORIAMETRICS_HOST:-127.0.0.1:9091}
envsubst -v VICTORIAMETRICS_HOST < /etc/nginx/conf.d/pmm.conf.template > /etc/nginx/conf.d/pmm.conf

exec nginx
