#!/bin/sh
#mac下运行./build.sh
#widnows下运行./build.sh win

target=admin
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${target}


os=$1
dir=/Volumes/WorkHD/workspace/go
if [ "${os}" == "win" ]; then
    dir=/F/project
fi

mkdir ${dir}/bin
rm ${dir}/bin/${target}
mv ./${target} ${dir}/bin/


scp ${dir}/bin/${target} root@39.108.186.82:/data/
ssh root@39.108.186.82 "cd /data && ./run.sh"
