#!/bin/sh
PLUGINS=~/.terraform.d/plugins

make build-only && \
  $GOBIN/dlv exec --headless --listen=:2345 --api-version=2 $PLUGINS/terraform-provider-chronicle -- --debug
