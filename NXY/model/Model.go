package model

import (
	"encoding/json"
	"fmt"
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
func NewTrade(tradeID string, userID string) *Trade {
	t := new(Trade)
	t.ID = tradeID
	t.Bills = make(map[string]TradeBill, 0)
	t.Signs = make([]Sign, 0)
	t.AddSign(STEP_INIT, userID)
	return t
}

func NewBill(params []string) *Bill {
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
func (t *Trade) AddBill(no string, price string) {
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
func (t *Trade) AddSign(stepID string, userID string) {
	if t.Signs == nil {
		t.Signs = make([]Sign, 0)
	}

	sign := Sign{stepID, userID}
	t.Signs = append(t.Signs, sign)
}
//========================================
//	[Bill]toJSON - JSON格式转换
//========================================
func (t *Bill) ToJSON() string {
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
func (t *Trade) ToJSON() string {
	bytes, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return "{}"
	}
	return string(bytes)
}