# dir=/F/project
dir=/Volumes/WorkHD/workspace/go

GOOS=linux GOARCH=amd64 go build

mkdir $dir/bin
rm $dir/bin/*
mv ./admin $dir/bin/
cd $dir/bin/
ssh root@39.108.186.82 "cd /data/ && rm admin"
scp -r -2 $dir/bin/* root@39.108.186.82:/data/
ssh root@39.108.186.82 "cd /data && ./run.sh"