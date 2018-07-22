#!/bin/sh

target=admin
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${target}
if [ $? -ne 0 ]; then
    echo "编译失败"
    exit 1
fi

os=$1
dir=/Volumes/WorkHD/workspace/go
if [ "${os}" == "win" ]; then
    dir=/F/project
fi

mkdir ${dir}/bin
rm ${dir}/bin/${target}
mv ./${target} ${dir}/bin/


ssh root@39.108.186.82 "cd /data && ./del.sh"
scp ${dir}/bin/${target} root@39.108.186.82:/data/
ssh root@39.108.186.82 "cd /data && ./run.sh"