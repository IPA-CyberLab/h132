#!/bin/bash

SWTPM_DIR=$(pwd)/var/swtpm_dir
echo "Starting SWTPM at ${SWTPM_DIR}/server.sock"
swtpm socket --tpm2 --tpmstate dir=${SWTPM_DIR} --server type=unixio,path=${SWTPM_DIR}/server.sock --flags startup-clear,not-need-init --log file=-,level=5
