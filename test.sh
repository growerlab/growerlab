#!/bin/bash

B=`cat ./.gitignore`

A=$(
cat <<- EOF
  HAHA
  $B
EOF
)
echo "===="
echo "$A"
