#!/bin/bash
OLD_GOPATH=${GOPATH}
export GOPATH=${PWD}
which go1
if [ $? -ne 0 ];then
    echo "go command not found" 
    exit 1
fi

go clean all
rm -rf td-client td-server td-proxy pkg bin
echo 'install td-client'
go install td-client
echo 'install td-guard'
go install td-guard
echo 'install td-proxy'
go install td-proxy
echo 'install td-server'
go install td-server

mkdir -p {td-client,td-server,td-proxy}/conf
mkdir -p packages
mv bin/td-client td-client/
mv bin/td-guard td-client/
cp src/etc/client.conf td-client/conf/
tar zcf packages/td-client.tar.gz td-client
mv bin/td-server td-server/
cp src/etc/server.conf td-server/conf/
tar zcf packages/td-server.tar.gz td-server
mv bin/td-proxy td-proxy/
cp src/etc/proxy.conf td-proxy/conf/
tar zcf packages/td-proxy.tar.gz td-proxy
rm -rf td-client td-server td-proxy pkg bin
export GOPATH=$OLD_GOPATH
echo 'complete，see ./packages'
echo
ls -l ./packages
