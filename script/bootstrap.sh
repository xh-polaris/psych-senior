#!/bin/bash
CURDIR=$(cd $(dirname $0); pwd)
BinaryName=psych-senior
echo "$CURDIR/bin/${BinaryName}"
exec $CURDIR/bin/${BinaryName}