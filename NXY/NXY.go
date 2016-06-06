package main

import (
	"github.com/qubing/chaincodes/NXY/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//==============================================================================
//	 Hyperledger Chaincode Method Definitions
//==============================================================================
type SimpleChaincode struct {
}

//票据总账
//{
//	"%机构编号%": [
//		{%Bill%},
//		...
//	],
//	...}
const KEY_BILLS = "BILLS"

//现金总账
//{"%机构编号%": "1,000,000", "%机构编号%": "1,000,000", ...}
const KEY_CASHES = "CASHES"

//交易总账
//{
//  "%交易编号%": {%TRADE%},
//  ...
//}
const KEY_TRADES = "TRADES"

//==============================================================================================================================
//	Init Function - Called when the user deploys the chaincode
//==============================================================================================================================
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	//init bills
	stub.PutState(KEY_BILLS, []byte("{}"))

	//init cashes
	stub.PutState(KEY_CASHES, []byte("{}"))

	//init cashes
	stub.PutState(KEY_TRADES, []byte("{}"))

	return nil, nil
}

//==============================================================================================================================
//	 Router Functions
//==============================================================================================================================
func (t *SimpleChaincode) inputBill(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 19 {
		return nil, errors.New("Parameter count is not correct.")
	}

	var bills map[string][]*model.Bill
	bytes, err := stub.GetState(KEY_BILLS)
	if err != nil {
		bytes = []byte("{}")
		fmt.Println("current bills:\n{}")
	} else {
		fmt.Printf("current bills:\n%s\n", string(bytes))
	}
	err = json.Unmarshal(bytes, &bills)
	if err != nil {
		bills = make(map[string][]*model.Bill)
		fmt.Println("current bills:\n Unmarshalling failed. ")
	} else {
		fmt.Println("current bills: \n Unmarshalling failed. ")
	}

	bill := model.NewBill(args)
	//机构编码
	bills[args[0]] = make([]*model.Bill, 0)
	bills[args[0]] = append(bills[args[0]], bill)
	bytes, err = json.Marshal(bills)
	if err != nil {
		fmt.Println("Bill JSON marshalling failed. ")
		return nil, errors.New("Bill JSON marshalling failed.")
	}
	stub.PutState(KEY_BILLS, bytes)

	return bytes, nil
}

//------------------------------------------------------------------------------------------------------------------------------
//	Invoke Router - Called on chaincode invoke. Takes a function name passed and calls that function. Converts some
//		  initial arguments passed to other things for use in the called function e.g. name -> ecert
//------------------------------------------------------------------------------------------------------------------------------
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "inputBill" {
		//录入票据
		bytes, err := t.inputBill(stub, args)
		if err != nil {
			return nil, err
		}
		return bytes, nil
	} else if function == "inputCash" {
		//输入现金
	} else if function == "tradeBill" {
		//申请卖出票据
	} else if function == "doSign" {
		//执行审批
	}
	return nil, errors.New("not valid invoke method")
}

func (t *SimpleChaincode) viewBill(stub *shim.ChaincodeStub, partyID string, billID string) ([]byte, error) {
	 var bills map[string] []*model.Bill
	 bytes, err := stub.GetState(KEY_BILLS)
	 if err != nil {
	 	bytes = []byte("{}")
	 }
	 if bytes != nil {
	 	fmt.Printf("current bills:\n %s \n", string(bytes))
	 }
	 err = json.Unmarshal(bytes, &bills)
	 if err != nil {
	 		bills = make(map[string] []*model.Bill)
	 }
	 for i:= 0; i < len(bills[partyID]); i++ {
		 if bills[partyID][i].No == billID {
			 bytes, err = json.Marshal(bills[partyID][i])
			 if err != nil {
			 		fmt.Println("Bill not found.")
			 		return nil, errors.New("Bill JSON marshalling failed.")
			 }
			 fmt.Printf("view bill:\n %s \n", string(bytes))
			 return bytes, nil
		 }
	 }
	 return nil, errors.New("Bill not found.")
}

func (t *SimpleChaincode) viewBills(stub *shim.ChaincodeStub, partyID string) ([]byte, error) {
	 var bills map[string] map[string] model.Bill
	 bytes, err := stub.GetState(KEY_BILLS)
	 if err != nil {
	 	bytes = []byte("{}")
	 }
	 if bytes != nil {
	 	fmt.Printf("current bills:\n %s \n", string(bytes))
	 }
	 err = json.Unmarshal(bytes, &bills)
	 if err != nil {
	 		bills = make(map[string] map[string] model.Bill)
			bills[partyID] = make(map[string] model.Bill)
	 }

	 bytes, err = json.Marshal(bills[partyID])
	 if err != nil {
	 		fmt.Println("Bill not found.")
			bytes = []byte("{}")
	 		//return nil, errors.New("Bill JSON marshalling failed.")
	 }
	 fmt.Printf("view bill:\n %s \n", string(bytes))
	 return bytes, nil
}

func (t *SimpleChaincode) viewTrades(stub *shim.ChaincodeStub, partyID string, tradeID string) ([]byte, error) {
	 var trades map[string] map[string] model.Trade
	 bytes, err := stub.GetState(KEY_TRADES)
	 if err != nil {
	 		bytes = []byte("{}")
	 }
	 if bytes != nil {
	 		fmt.Printf("current trades:\n %s \n", string(bytes))
	 }
	 err = json.Unmarshal(bytes, &trades)
	 if err != nil {
	 		trades = make(map[string] map[string] model.Trade)
			trades[partyID] = make(map[string] model.Trade)
			fmt.Println("Trades not found.")
	 		return nil, errors.New("Trade JSON marshalling failed.")
	 }

	 bytes, err = json.Marshal(trades[partyID][tradeID])
	 if err != nil {
	 		fmt.Println("Trades not found.")
	 		return nil, errors.New("Trade JSON marshalling failed.")
	 }
	 fmt.Printf("view bill:\n %s \n", string(bytes))
	 return bytes, nil
}

//=================================================================================================================================
//	Query - Called on chaincode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
//=================================================================================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "viewBill" {
			fmt.Println("view Bill......>")
			return t.viewBill(stub, args[0], args[1])
	} else if function == "viewTrades" {
			return t.viewTrades(stub, args[0], args[1])
	} else if function == "viewBills" {
			return t.viewBills(stub, args[0])
	} else if function == "viewCashes" {

	}
	return nil, errors.New("Received unknown function invocation")
}

//=================================================================================================================================
//	 Main - main - Starts up the chaincode
//=================================================================================================================================
 func main() {
 	err := shim.Start(new(SimpleChaincode))
 	if err != nil {
 		fmt.Printf("Error starting Chaincode: %s", err)
 	}

	 //args := []string{"P1", "00001", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q"}
	 //var bills = make(map[string] map[string] *model.Bill)
	 //bill := model.NewBill(args)
	 //bills[args[0]] = make(map[string] *model.Bill)
	 //bills[args[0]][bill.No] = bill
	 //fmt.Println(bill.ToJSON())
 }
