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
  Price  string `json:"price"`
  Detail *Bill `json:"detail"`
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
  Signs []Sign `json:"signs"`
  Bills map[string]*TradeBill `json:"bills"`
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
  t.Bills = make(map[string]*TradeBill, 0)
  t.Signs = make([]Sign, 0)
  t.addSign(STEP_INIT, userID)
  return t
}

func newBill(params []string) *Bill {
	bill := new(Bill)
	bill.No = params[0]
	bill.Type = params[1]
	bill.IssuerName = params[2]
	bill.IssuerAccount = params[3]
	bill.IssuerBank = params[4]
	bill.CustodianName = params[5]
	bill.CustodianAccount = params[6]
	bill.CustodianBank = params[7]
	bill.FaceAmount = params[8]
	bill.AcceptorName = params[9]
	bill.AcceptorAccount = params[10]
	bill.AcceptorBank = params[11]
	bill.IssueDate = params[12]
	bill.DueDate = params[13]
	bill.AcceptDate = params[14]
	bill.PayBank = params[15]
	bill.TransEnable = params[16]
	return bill
}

//=======================================
//	[Trade]addBill - 加入票据
//  params: 票据参数
//=======================================
func (t *Trade) addBill(no string, price string) {
  if t.Bills == nil {
    t.Bills = make(map[string]*TradeBill, 0)
  }
	//TODO: get Bill
	bill := new(Bill)
	bill.No = no
  tradeBill := new(TradeBill)
	tradeBill.Price = price
	tradeBill.Detail = bill
  t.Bills[bill.No]= tradeBill
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
	return nil, nil
}

//==============================================================================================================================
//	 Router Functions
//==============================================================================================================================
//	Invoke - Called on chaincode invoke. Takes a function name passed and calls that function. Converts some
//		  initial arguments passed to other things for use in the called function e.g. name -> ecert
//==============================================================================================================================
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "inputBill" {
		//录入票据
	} else if function == "inputCash" {
		//输入现金
	} else if function == "tradeBill" {
		//申请卖出票据
	} else if function == "doSign" {
		//执行审批
	}
	return nil, errors.New("not valid invoke method")
}

//=================================================================================================================================
//	Query - Called on chaincode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
//=================================================================================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "getBill" {

	} else if function == "getTrade" {

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
  trade := newTrade("0000000001", "user01")
  fmt.Println(trade.toJSON())
  trade.addSign("02", "user02")
  fmt.Println(trade.toJSON())
  trade.addBill("01", "")
  fmt.Println(trade.toJSON())
}
