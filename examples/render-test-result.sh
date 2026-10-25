#!/usr/bin/env bash
# Run the actual test command, render its result, and preserve its exit status.
set -eu
renderer=${1:?Usage: render-test-result.sh RENDERER OUTPUT TEST_COMMAND [ARGS...]}
output=${2:?Set the SVG output path}
shift 2
if [ "$#" -eq 0 ]; then
  echo 'A test command is required' >&2
  exit 2
fi

result=0
"$@" || result=$?
status=passing
if [ "$result" -ne 0 ]; then status=failing; fi
"$renderer" render --status "$status" --character kappa --theme dark --output "$output"
exit "$result"
