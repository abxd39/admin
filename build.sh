GOOS=linux GOARCH=amd64 go build
mv ./amdin ../bin/
cd ../bin/
ssh root@39.108.186.82   "cd /data/ && rm admin"
scp -r -2 /F/project/src/bin/* root@39.108.186.82:/data/
chmod +x ./admin
ps -ef | grep admin | awk '{print $2}' | xargs kill -9
nohup ./amdin &