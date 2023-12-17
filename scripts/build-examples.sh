#!/bin/bash

set -euo pipefail
cd ${0%/*}

find ../_examples -type f -name "*.go" -exec go generate -v {} \;
