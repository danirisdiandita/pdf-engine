#!/usr/bin/env bash
set -euo pipefail

# If no args passed, show usage
if [ $# -eq 0 ]; then
  echo "Usage: docker run --rm -v \$(pwd):/data pandoc-multilang input.md output.pdf"
  exit 1
fi

# If exactly two args, treat as input and output
if [ $# -ge 2 ]; then
  INPUT="$1"
  OUTPUT="$2"
  shift 2
else
  # fallback
  INPUT="$1"
  OUTPUT="${INPUT%.*}.pdf"
  shift 1
fi

# Default pdf engine is xelatex (good Unicode/CJK/RTL support)
# Pass additional pandoc flags after the input/output if you want
pandoc "$INPUT" -o "$OUTPUT" --pdf-engine=xelatex "$@"

# exit with pandoc's exit code
exit $?
