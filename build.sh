# dir=/F/project
dir=/Volumes/WorkHD/workspace/go

target=admin

GOOS=linux GOARCH=amd64 go build -o ${target}

mkdir ${dir}/bin
rm ${dir}/bin/${target}
mv ./${target} ${dir}/bin/

ssh root@39.108.186.82 "cd /data/ && rm ${target}"
scp ${dir}/bin/${target} root@39.108.186.82:/data/
ssh root@39.108.186.82 "cd /data && ./run.sh"