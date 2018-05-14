#!/bin/bash
set -e

arch="amd64"
os="linux"
if [[ $TRAVIS_OS_NAME = osx ]]; then
    os="darwin"
fi

tag=$(wget -qSO- --method=HEAD --max-redirect=0 https://github.com/smola/ci-tricks/releases/latest 2>&1 | grep Location: | sed -e 's/.*\///g')
wget -qO ci-tricks https://github.com/smola/ci-tricks/releases/download/${tag}/ci-tricks_${os}_${arch}
chmod +x ci-tricks
exec ./ci-tricks