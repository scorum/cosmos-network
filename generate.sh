#!/usr/bin/env bash
# Script generates genesis.json with DEVNET validators and supervisors
# The output of the script is
# - genesis.json - the genesis file itself
# - accounts.txt - validator and supervisor accounts: name, address, public key, mnemonic
# - keyring - keyring directory

set -e
set -o pipefail

if [[ ${#PASSWORD} -lt 8 ]]; then
  echo "PASSWORD variable must be set and not less than 8 symbols" && exit 1
fi

if [[ -z "${CHAIN_ID}" ]]; then
  CHAIN_ID="devnet-1"
fi

if [[ -z "${VALIDATORS}" ]]; then
  VALIDATORS="validator1"
fi

if [[ -z "${SUPERVISORS}" ]]; then
  SUPERVISORS="orion pulsar aurora meteor nebula quasar supernova eclipse comet galaxy"
fi

echo "CHAIN_ID is set to $CHAIN_ID"
echo "VALIDATORS is set to $VALIDATORS"
echo "SUPERVISORS is set to $SUPERVISORS"

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
  for account in ${VALIDATORS} ${SUPERVISORS}
  do
    echo $PASSWORD$'\n'$PASSWORD | scorumd keys add ${account} --keyring-backend=file --home=${homed} --keyring-dir ${keyring}/${account} >> accounts.txt 2>&1
  done
}

function gentx() {
  for account in ${VALIDATORS}
  do
    addr=$(echo $PASSWORD | scorumd keys show ${account} -a --home=${homed} --keyring-backend=file --keyring-dir ${keyring}/${account})
    rm -rf ${homed}
    scorumd init ${account} --staking-bond-denom nsp --chain-id=${CHAIN_ID} --home=${homed}
    scorumd add-genesis-account ${addr} 1000000000000nscr,1000000000000nsp --home=${homed}
    scorumd add-genesis-supervisor ${addr} --home=${homed}

    echo $PASSWORD | scorumd gentx ${account} 100000000000nsp \
      --website "https://scorum.com" --home=${homed} \
      --chain-id=${CHAIN_ID} --keyring-backend=file \
      --keyring-dir ${keyring}/${account} \
      --commission-max-change-rate "1" --commission-max-rate "1" --commission-rate "1"

    mkdir -p ${gentx} && cp -a ${homed}/config/gentx/* ${gentx}/

    mkdir -p "${keys}/${account}"
    cp ${homed}/config/priv_validator_key.json ${keys}/${account}/priv_validator_key.json
  done
}

function genesis() {
   rm -rf ${homed}
   moniker=`echo $VALIDATORS | awk '{print $1}'`

   # add genesis accounts
   scorumd init "${moniker}" --staking-bond-denom nsp --chain-id=${CHAIN_ID} --home=${homed}

   for account in ${VALIDATORS} ${SUPERVISORS}
   do
     addr=$(echo $PASSWORD | scorumd keys show ${account} -a --keyring-backend=file --home=${homed} --keyring-dir ${keyring}/${account})

     scorumd add-genesis-account ${addr} 1000000000000nscr,1000000000000nsp --home=${homed}
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
sed -i -e 's/"stake"/"nsp"/g' ${homed}/config/genesis.json
# set inflation fields to 0, because there will be another reward mechanism
sed -i 's/"\(inflation[^"]*\)": "[0-9.]\+",/"\1": "0",/g' ${homed}/config/genesis.json
# set min_commission_rate to 1
sed -i 's/"min_commission_rate": "[0-9.]\+"/"min_commission_rate": "1"/' ${homed}/config/genesis.json

cp ${homed}/config/genesis.json genesis.json

#rm ./genesis.json.bak
scorumd validate-genesis genesis.json

