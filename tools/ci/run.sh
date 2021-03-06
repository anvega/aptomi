#!/bin/bash
set -exou pipefail

sudo rm -rf /usr/local/bin/*aptomi* /etc/aptomi /var/lib/aptomi ~/.aptomi

hostname

ls -la

export PATH=$PATH:"$WORKSPACE"/bin
export GOPATH="$WORKSPACE"
export GOROOT=/usr/local/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
mkdir -p "$GOPATH/bin"

pushd src/github.com/Aptomi/aptomi

make vendor

tools/demo-ldap.sh

make lint

make smoke

source /jenkins/aptomi-coveralls.io
make coverage-full coverage-publish

tools/test-install.sh

tools/publish-docker.sh

tools/publish-charts.sh

popd
