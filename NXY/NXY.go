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

//现金&票据总账
//{"%机构编号%": {"NO": "%机构编号%", "AMOUNT":"1,000,000", "BILLS": ["%票据编号%", ...]}, ...}
const KEY_ACCOUNTS = "ACCOUNTS"

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
	stub.PutState(KEY_ACCOUNTS, []byte("{}"))

	//init cashes
	stub.PutState(KEY_TRADES, []byte("{}"))

	return nil, nil
}

//==============================================================================================================================
//	 Router Functions
//==============================================================================================================================
func (t *SimpleChaincode) inputBill(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	// if len(args) != 19 {
	// 	return nil, errors.New("Parameter count is not correct.")
	// }

	var bills map[string][]*model.Bill
	bytes, err := stub.GetState(KEY_BILLS)
	if err != nil {
		bytes = []byte("{}")
		fmt.Println("current bills:\n{}")
	}
	err = json.Unmarshal(bytes, &bills)
	if err != nil {
		bills = make(map[string][]*model.Bill)
		fmt.Println("current bills:\n Unmarshalling failed. ")
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

func (t *SimpleChaincode) inputCash(stub *shim.ChaincodeStub, partyID string, amount string) ([]byte, error) {
	var accounts map[string]*model.Account
	bytes, err := stub.GetState(KEY_ACCOUNTS)
	if err != nil {
		bytes = []byte("{}")
		fmt.Println("current bills:\n{}")
	}
	err = json.Unmarshal(bytes, &accounts)
	if err != nil {
		accounts = make(map[string]*model.Account)
		fmt.Println("current bills:\n Unmarshalling failed. ")
	} else {
		fmt.Println("current bills: \n Unmarshalling failed. ")
	}

	//机构编码
	accounts[partyID] = model.NewAccount(partyID, amount)
	//{partyID, amount}
	bytes, err = json.Marshal(accounts)
	if err != nil {
		fmt.Println("Bill JSON marshalling failed. ")
		return nil, errors.New("Bill JSON marshalling failed.")
	}
	stub.PutState(KEY_ACCOUNTS, bytes)

	return nil, nil
}

func (t *SimpleChaincode) tradeBill(stub *shim.ChaincodeStub, userID string, tradeID string, billID string, partyFrom string, partyTo string, price string) ([]byte, error) {
	var trades map[string]*model.Trade
	bytes, err := stub.GetState(KEY_TRADES)
	if err != nil {
		bytes = []byte("{}")
		fmt.Println("current trades:\n{}")
	}
	err = json.Unmarshal(bytes, &trades)
	if err != nil {
		fmt.Println("current trades:\n Unmarshalling failed. ")
	}

	//机构编码
	trades[tradeID] = model.NewTrade(tradeID, partyFrom, partyTo, userID)
	//{partyID, amount}
	bytes, err = json.Marshal(trades)
	if err != nil {
		fmt.Println("Bill JSON marshalling failed. ")
		return nil, errors.New("Bill JSON marshalling failed.")
	}
	stub.PutState(KEY_ACCOUNTS, bytes)

	return nil, nil
}

func (t *SimpleChaincode) doSign(stub *shim.ChaincodeStub, tradeID string, step string, userID string, comments string) ([]byte, error) {
	bytes, err := stub.GetState(KEY_TRADES)
	if err != nil {
		bytes = []byte("{}")
		fmt.Println("current trades:\n{}")
	}
	var trades map[string]*model.Trade
	err = json.Unmarshal(bytes, &trades)
	if err != nil {
		return nil, errors.New("Trade Not Found.")
	}

	trades[tradeID].AddSign(step, userID, comments)

	return nil, nil
}

//------------------------------------------------------------------------------------------------------------------------------
//	Invoke Router - Called on chaincode invoke. Takes a function name passed and calls that function. Converts some
//		  initial arguments passed to other things for use in the called function e.g. name -> ecert
//------------------------------------------------------------------------------------------------------------------------------
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "inputBill" {
		//录入票据
		if len(args) != 19 {
			return nil, errors.New("[inputBill]Parameter count is not correct.")
		}
		bytes, err := t.inputBill(stub, args)
		if err != nil {
			return nil, err
		}
		return bytes, nil
	} else if function == "inputCash" {
		//输入现金
		if len(args) != 2 {
			return nil, errors.New("[inputCash]Parameter count is not correct.")
		}
		bytes, err := t.inputCash(stub, args[0], args[1])
		if err != nil {
			return nil, err
		}
		return bytes, nil
	} else if function == "tradeBill" {
		//申请卖出票据
		if len(args) != 6 {
			return nil, errors.New("[tradeBill]Parameter count is not correct.")
		}
		bytes, err := t.tradeBill(stub, args[0], args[1], args[2], args[3], args[4], args[5])
		if err != nil {
			return nil, err
		}
		return bytes, nil
	} else if function == "doSign" {
		//执行审批
		if len(args) != 4 {
			return nil, errors.New("[doSign]Parameter count is not correct.")
		}
		bytes, err := t.doSign(stub, args[0], args[1], args[2], args[3])
		if err != nil {
			return nil, err
		}
		return bytes, nil
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
func (t *SimpleChaincode) viewAccount(stub *shim.ChaincodeStub, partyID string) ([]byte, error) {
	 var accounts map[string] model.Account
	 bytes, err := stub.GetState(KEY_ACCOUNTS)
	 if err != nil {
			bytes = []byte("{}")
	 }
	 if bytes != nil {
			fmt.Printf("current trades:\n %s \n", string(bytes))
	 }
	 err = json.Unmarshal(bytes, &accounts)
	 if err != nil {
			fmt.Println("Account not found.")
			return nil, errors.New("Trade JSON marshalling failed.")
	 }

	 bytes, err = json.Marshal(accounts[partyID])
	 if err != nil {
			fmt.Println("Account not found.")
			return nil, errors.New("Account JSON marshalling failed.")
	 }
	 fmt.Printf("view Account:\n %s \n", string(bytes))
	 return bytes, nil
}
//=================================================================================================================================
//	Query - Called on chaincode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
//=================================================================================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "viewBill" {
		return t.viewBill(stub, args[0], args[1])
	} else if function == "viewTrades" {
		return t.viewTrades(stub, args[0], args[1])
	} else if function == "viewBills" {
		return t.viewBills(stub, args[0])
	} else if function == "viewCashes" {
		return t.viewAccount(stub, args[0])
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
