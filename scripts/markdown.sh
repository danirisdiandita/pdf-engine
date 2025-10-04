#!/bin/bash

# pandoc ../markdown/sample.md -o output.pdf


pandoc ../markdown/sample.md -o output.pdf \
  --pdf-engine=xelatex \
  -V geometry:margin=1in \
  -V mainfont=Helvetica

