#!/bin/sh
set -eu
base="${BASE_URL:-http://127.0.0.1:8086}"
curl -fsS "$base/healthz"
curl -fsS "$base/readyz"
curl -fsS "$base/metrics"
