#!/bin/bash

set -euo pipefail
cd ${0%/*}

find . -type f \( -name "*.md" -or -name '*.txt' -or -name '*.html' \) -exec rm -v {} \;
