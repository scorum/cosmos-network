#!/usr/bin/env bash
# Script generates genesis.json with DEVNET validators
# The ouput of the script is
# - genesis.json - the genesis file itself
# - accounts.txt - validator accounts: name, address, public key, mnemonic
# - keyring - keyring directory

set -e
set -o pipefail

if [[ ${#PASSWORD} -lt 8 ]]; then
  echo "PASSWORD variable must be set and not less than 8 symbols" && exit 1
fi

if [[ -z "${CHAIN_ID}" ]]; then
  CHAIN_ID="devnet-1"
fi

if [[ -z "${ACCOUNTS}" ]]; then
  ACCOUNTS="validator1 validator2 validator3 validator4"
fi

echo "CHAIN_ID is set to $CHAIN_ID"
echo "ACCOUNTS is set to $ACCOUNTS"

out=output
keys=keys
keyring=keyring

tmp=$(mktemp -d)
homed=${tmp}/homed
gentx=${tmp}/gentx

echo "homed dir is set to ${homed}"
echo "gentx dir is set to ${gentx}"

function join_by { local IFS="$1"; shift; echo "$*"; }

scorumd config output json --home=${homed}

function addKeys() {
  for account in ${ACCOUNTS}
  do
    echo $PASSWORD$'\n'$PASSWORD | scorumd keys add ${account} --keyring-backend=file --home=${homed} --keyring-dir ${keyring}/${account} >> accounts.txt 2>&1
  done
}

function gentx() {
  for account in ${ACCOUNTS}
  do
    addr=$(echo $PASSWORD | scorumd keys show ${account} -a --home=${homed} --keyring-backend=file --keyring-dir ${keyring}/${account})
    rm -rf ${homed}
    scorumd init ${account} --staking-bond-denom sp --chain-id=${CHAIN_ID} --home=${homed}
    scorumd add-genesis-account ${addr} 1000000000000scr,1000000000000sp --home=${homed}
    scorumd add-genesis-supervisor ${addr} --home=${homed}
    echo $PASSWORD | scorumd gentx ${account} 100000000000sp --website "https://scorum.com" --home=${homed} --chain-id=${CHAIN_ID} --keyring-backend=file --keyring-dir ${keyring}/${account}
    mkdir -p ${gentx} && cp -a ${homed}/config/gentx/* ${gentx}/

    mkdir -p "${keys}/${account}"
    cp ${homed}/config/priv_validator_key.json ${keys}/${account}/priv_validator_key.json
  done
}

function genesis() {
   rm -rf ${homed}
   moniker=`echo $ACCOUNTS | awk '{print $1}'`

   # add genesis accounts
   scorumd init "${moniker}" --staking-bond-denom sp --chain-id=${CHAIN_ID} --home=${homed}

   for account in ${ACCOUNTS}
   do
     addr=$(echo $PASSWORD | scorumd keys show ${account} -a --keyring-backend=file --home=${homed} --keyring-dir ${keyring}/${account})

     scorumd add-genesis-account ${addr} 1000000000000scr,1000000000000sp --home=${homed}
     scorumd add-genesis-supervisor ${addr} --home=${homed}
   done

   # generate genesis.json
   scorumd collect-gentxs --gentx-dir ${gentx} --home=${homed}
}

# prepare output folder
rm -rf $out
mkdir $out && cd $out || exit 1

# do
addKeys
gentx
genesis

# extra update genesis
sed -i -e 's/"stake"/"sp"/g' ${homed}/config/genesis.json
# set inflation fields to 0, because there will be another reward mechanism
sed -i 's/"\(inflation[^"]*\)": "[0-9.]\+",/"\1": "0",/g' ${homed}/config/genesis.json

cp ${homed}/config/genesis.json genesis.json

#rm ./genesis.json.bak
scorumd validate-genesis genesis.json

