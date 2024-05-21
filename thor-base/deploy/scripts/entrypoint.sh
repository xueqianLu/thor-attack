#! /bin/sh

# exit script on any error
set -e

echo "start sender with account $ACCOUNT_IDX"
/usr/bin/sender -index ${ACCOUNT_IDX} >> /root/node/sender.log 2>&1 &

echo "start thor with beneficiary $BENEFICIARY and bootnode $BOOTNODE"
/usr/bin/thor --data-dir /root/node --verbosity 4 --network /root/genesis.json --beneficiary ${BENEFICIARY} --api-addr "0.0.0.0:8669" --bootnode ${BOOTNODE} >> /root/node/node.log 2>&1
