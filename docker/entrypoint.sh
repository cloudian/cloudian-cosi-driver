#!/usr/bin/env sh

set -euxo pipefail

update-ca-certificates

exec /cosi-driver
