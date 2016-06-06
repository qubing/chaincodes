package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)
//==============================================================================
//	 Data Model Definitions
//==============================================================================

//=====================================
//	Bill - 票据信息
//=====================================
type Bill struct {
	No               string `json:"no"`
	Attr             string `json:"attr"`
	Type             string `json:"type"`
	IssuerName       string `json:"issuer_name"`
	IssuerAccount    string `json:"issuer_account"`
	IssuerBank       string `json:"issuer_bank"`
	CustodianName    string `json:"custodian_name"`
	CustodianAccount string `json:"custodian_account"`
	CustodianBank    string `json:"custodian_bank"`
	FaceAmount       string `json:"face_amount"`
	AcceptorName     string `json:"acceptor_name"`
	AcceptorAccount  string `json:"cceptor_account"`
	AcceptorBank     string `json:"cceptor_bank"`
	IssueDate        string `json:"issue_date"`
	DueDate          string `json:"due_date"`
	AcceptDate       string `json:"accept_date"`
	PayBank          string `json:"pay_bank"`
	TransEnable      string `json:"trans_enable"`
}
//=====================================
//	TradeBill - 交易票据信息
//=====================================
type TradeBill struct {
	No    string `json:"no"`
	Price string `json:"price"`
}

//=====================================
//	Sign - 审批信息
//=====================================
type Sign struct {
	StepID   string `json:"step"`
	Operator string `json:"operator"`
}
//======================================
//	Trade - 交易信息
//======================================
type Trade struct {
	ID    string `json:"id"`
	From  string `json:"from"`
	To    string `json:"to"`
	Signs []Sign `json:"signs"`
	Bills map[string]TradeBill `json:"bills"`
}

const STEP_INIT = "01"

//======================================
//	newTrade - 交易初始化
//  tradeID: 交易ID
//  userID: 用户ID
//======================================
func newTrade(tradeID string, userID string) *Trade {
	t := new(Trade)
	t.ID = tradeID
	t.Bills = make(map[string]TradeBill, 0)
	t.Signs = make([]Sign, 0)
	t.addSign(STEP_INIT, userID)
	return t
}

func newBill(params []string) *Bill {
	bill := new(Bill)
	bill.No = params[1]
	bill.Attr = params[2]
	bill.Type = params[3]
	bill.IssuerName = params[4]
	bill.IssuerAccount = params[5]
	bill.IssuerBank = params[6]
	bill.CustodianName = params[7]
	bill.CustodianAccount = params[8]
	bill.CustodianBank = params[9]
	bill.FaceAmount = params[10]
	bill.AcceptorName = params[11]
	bill.AcceptorAccount = params[12]
	bill.AcceptorBank = params[13]
	bill.IssueDate = params[14]
	bill.DueDate = params[15]
	bill.AcceptDate = params[16]
	bill.PayBank = params[17]
	bill.TransEnable = params[18]
	return bill
}

//=======================================
//	[Trade]addBill - 加入票据
//  params: 票据参数
//=======================================
func (t *Trade) addBill(no string, price string) {
	if t.Bills == nil {
		t.Bills = make(map[string]TradeBill, 0)
	}

	t.Bills[no] = TradeBill{no, price}
}
//========================================
//	[Trade]addSign - 添加审批
//  stepID: 审批步骤
//  userID: 用户ID
//========================================
func (t *Trade) addSign(stepID string, userID string) {
	if t.Signs == nil {
		t.Signs = make([]Sign, 0)
	}

	sign := Sign{stepID, userID}
	t.Signs = append(t.Signs, sign)
}
//========================================
//	[Bill]toJSON - JSON格式转换
//========================================
func (t *Bill) toJSON() string {
	bytes, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return "{}"
	}
	return string(bytes)
}
//========================================
//	[Trade]toJSON - JSON格式转换
//========================================
func (t *Trade) toJSON() string {
	bytes, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return "{}"
	}
	return string(bytes)
}

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

	var bills map[string] map[string] Bill
	bytes, err := stub.GetState(KEY_BILLS)
	if err != nil {
		bytes = []byte("{}")
	}
	err = json.Unmarshal(bytes, &bills)
	if err != nil {
		bills = make(map[string]map[string]Bill)
	}

	bill := newBill(args)
	//机构编码
	bills[args[0]][bill.No] = *bill
	bytes, err = json.Marshal(bills)
	if err != nil {
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
	var bills map[string] map[string] Bill
	bytes, err := stub.GetState(KEY_BILLS)
	if err != nil {
		bytes = []byte("{}")
	}
	err = json.Unmarshal(bytes, &bills)
	if err != nil {
		bills = make(map[string]map[string]Bill)
	}

	bytes, err = json.Marshal(bills[party][no])
	if err != nil {
		return nil, errors.New("Bill JSON marshalling failed.")
	}

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

	} else if function == "viewBills" {

	} else if function == "viewCashes" {

	}
	return nil, errors.New("Received unknown function invocation")
}

//=================================================================================================================================
//	 Main - main - Starts up the chaincode
//=================================================================================================================================
// func main() {
// 	err := shim.Start(new(SimpleChaincode))
// 	if err != nil {
// 		fmt.Printf("Error starting Chaincode: %s", err)
// 	}
// }
func main() {
	//trade := newTrade("0000000001", "user01")
	//fmt.Println(trade.toJSON())
	//trade.addSign("02", "user02")
	//fmt.Println(trade.toJSON())
	//trade.addBill("01", "")
	//fmt.Println(trade.toJSON())
	//init bills
	bills := make(map[string]Bill, 0)
	bill := newBill([]string{"01","","","","","","","","","","","","","","","","","","","",""})
	bills["01"] = *bill
	bytes, err := json.Marshal(bills)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(bytes))

	err = json.Unmarshal(bytes, &bills)
	if err != nil {
		fmt.Println(err)
		return
	}

	bytes, err = json.Marshal(bills["01"])

	fmt.Println(string(bytes))
}
