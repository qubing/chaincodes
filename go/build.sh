mkdir $GOPATH/src/github.com
rm $GOPATH/src/github.com/chaincodes
ln -s 
ln -s $PWD $GOPATH/src/github.com/chaincodes

#balance_mgr
cd $GOPATH/src/github.com/chaincodes/balance_mgr
rm -rf vendor
govendor init
govendor add +external
govendor add github.com/hyperledger/fabric/peer
govendor add github.com/hyperledger/fabric/core/chaincode/lib/cid
go build

#crypto_ledger
cd $GOPATH/src/github.com/chaincodes/crypto_ledger
rm -rf vendor
govendor init
govendor add +external
govendor add github.com/hyperledger/fabric/peer
govendor add github.com/hyperledger/fabric/core/chaincode/lib/cid
go build