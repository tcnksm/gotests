#!/bin/bash
set -e

DIR=$(cd $(dirname ${0})/.. && pwd)
cd ${DIR}

make build
DEBUG=1 ./bin/goatg -d ./test-fixtures/basic.go
DEBUG=1 ./bin/goatg -d ./test-fixtures/no-test.go
