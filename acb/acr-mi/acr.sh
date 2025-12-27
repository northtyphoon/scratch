#!/usr/bin/env bash
set -Eeuo pipefail

az login --identity

az acr login -n $Registry --expose-token --output tsv --query accessToken | acr login $Registry -u 00000000-0000-0000-0000-000000000000 --password-stdin

acr "$@"