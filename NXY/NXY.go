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
//{"%机构编号%": [
//	{
//		"no": "",
//		"attr": "",
//		"type": "",
//		"issuer_name": "",
//		"issuer_account": "",
//		"issuer_bank": "",
//		"custodian_name": "",
//		"custodian_account": "",
//		"custodian_bank": "",
//		"face_amount": "",
//		"acceptor_name": "",
//		"acceptor_account": "",
//		"acceptor_bank": "",
//		"issue_date": "",
//		"due_date": "",
//		"accept_date": "",
//		"pay_bank": "",
//		"trans_enable": "",
//	},
//	{...}, ...
//	],
//"%机构编号%": [...], ...}
const KEY_BILLS = "BILLS"

//现金总账
//{"%机构编号%": "1,000,000", "%机构编号%": "1,000,000", ...}
const KEY_CASHES = "CASHES"

//交易总账
//{"%交易编号%": {
//	"no": "",
//	"from": "%卖方机构编号%",
//	"to": "买方机构编号",
//	"bills": [
//		{"no": "%票据编号%",
//			"amount": "950,000"
//		}, {...}, ...]},
//	"%交易编号%": {...}, ...}
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

	var bills map[string] map[string] model.Bill
	bytes, err := stub.GetState(KEY_BILLS)
	if err != nil {
		bytes = []byte("{}")
		fmt.Println("current bills:\n{}")
	} else {
		fmt.Printf("current bills:\n%s\n", string(bytes))
	}
	err = json.Unmarshal(bytes, &bills)
	if err != nil {
		bills = make(map[string]map[string]model.Bill)
		fmt.Println("current bills:\n Unmarshalling failed. ")
	} else {
		fmt.Println("current bills: \n Unmarshalling failed. ")
	}

	bill := model.NewBill(args)
	//机构编码
	bills[args[0]] = make(map[string] *model.Bill)
	bills[args[0]][bill.No] = *bill
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

func (t *SimpleChaincode) viewBill(stub *shim.ChaincodeStub, party string, no string) ([]byte, error) {
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
	 	bills = make(map[string]map[string]model.Bill)
	 }

	 bytes, err = json.Marshal(bills[party][no])
	 if err != nil {
	 	fmt.Println("Bill not found.")
	 	return nil, errors.New("Bill JSON marshalling failed.")
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

	} else if function == "viewBills" {

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

//func main() {
//	//trade := newTrade("0000000001", "user01")
//	//fmt.Println(trade.toJSON())
//	//trade.addSign("02", "user02")
//	//fmt.Println(trade.toJSON())
//	//trade.addBill("01", "")
//	//fmt.Println(trade.toJSON())
//	//init bills
//	bills := make(map[string]model.Bill, 0)
//	bill := model.NewBill([]string{"01","","","","","","","","","","","","","","","","","","","",""})
//	bills["01"] = *bill
//	bytes, err := json.Marshal(bills)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(string(bytes))
//
//	err = json.Unmarshal(bytes, &bills)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	bytes, err = json.Marshal(bills["01"])
//
//	fmt.Println(string(bytes))
//}
