#!/bin/bash

CURRENT_DIR=$(pwd)
CMD_DIR=${CURRENT_DIR}/cmd
PROD_BINARY_NAME=gogenm

# Build to binary.
cd ${CMD_DIR} && go build -o ${CMD_DIR}/bin/${PROD_BINARY_NAME} -v .

if [ -d "$GOPATH/bin" ]; then
	# Expose binary path to go bin.
	mv ${CMD_DIR}/bin/${PROD_BINARY_NAME} ${GOPATH}/bin

	echo "build successed!"
else
	# https://stackoverflow.com/questions/2990414/echo-that-outputs-to-stderr
	>&2 echo "go bin path not found, please makesure go bin path exists"
fi
