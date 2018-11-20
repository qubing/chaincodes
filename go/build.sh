mkdir $GOPATH/src/github.com
rm $GOPATH/src/github.com/chaincode
ln -s 
ln -s $PWD $GOPATH/src/github.com/chaincode
cd $GOPATH/src/github.com/chaincode/balance_mgr
rm -rf vendor
govendor init
govendor add +external
govendor add github.com/hyperledger/fabric/peer
govendor add github.com/hyperledger/fabric/core/chaincode/lib/cid
go build