#!/bin/bash

CURRENT_DIR=$(pwd)
CMD_DIR=${CURRENT_DIR}/cmd
PROD_BINARY_NAME=gen-model

# Build to binary.
cd ${CMD_DIR} && go build -o ${CMD_DIR}/bin/${PROD_BINARY_NAME} -v .

if [ -d "$GOPATH/bin" ]; then
	# Expose binary path to go bin.
	mv ${CMD_DIR}/bin/${PROD_BINARY_NAME} ${GOPATH}/bin

	# check if sqlc cli exists
	if !command -v sqlc &> /dev/null; then
		echo "go package github.com/kyleconroy/sqlc is not installed. installing..."
		go install github.com/kyleconroy/sqlc/cmd/sqlc@latest

		echo << EndOfMessage
		dont forget to create 'sqlc.yaml' and 'query' directory to proceed using 'sqlc' command line.
		you can copy the default content of 'sqlc.yaml' from:

			https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html

		Modify the content to suite your need.
EndOfMessage
	fi

	echo "build successed!"
else
	# https://stackoverflow.com/questions/2990414/echo-that-outputs-to-stderr
	>&2 echo "go bin path not found, please makesure go bin path exists"
fi
