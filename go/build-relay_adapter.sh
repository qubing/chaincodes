mkdir $GOPATH/src/github.com
rm $GOPATH/src/github.com/chaincodes
ln -s 
ln -s $PWD $GOPATH/src/github.com/chaincodes

#balance_mgr
cd $GOPATH/src/github.com/chaincodes/relay_adapter
rm -rf vendor
govendor init
govendor add +external
govendor add github.com/hyperledger/fabric/peer
go build
