#!/bin/bash
set -euo pipefail

export GO_VERSION=1.16

pip install Pygments

go run .
