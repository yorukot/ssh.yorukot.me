#!/bin/sh

set -eu

printf '%s\n' 'Running backend checks...'
(
	make lint
)

printf '%s\n' 'Pre-commit checks passed.'