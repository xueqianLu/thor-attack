#!/bin/bash
testcase=${1:-"base"}

composefile=docker-compose-${testcase}.yml

cp docker-compose.yml ${composefile}_bak

make clean

sed "26s/base/${testcase}/g" ${composefile}_bak > $composefile
rm -f ${composefile}_bak

docker compose -f ${composefile} up -d
#sleep 60
sleep 1800

docker exec -it thor-node-6 /usr/bin/query -url http://127.0.0.1:8669 > data/query.log

docker compose down

sudo cp -r data testdata_${testcase}
