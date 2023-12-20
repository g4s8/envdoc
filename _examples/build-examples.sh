#!/bin/bash

set -euo pipefail
cd ${0%/*}

find . -type f -name "*.go" -exec go generate -v {} \;
