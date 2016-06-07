package model
import (
	"encoding/json"
	"fmt"
)

//=====================================
//	Bill - 票据信息
//=====================================
type Bill struct {
	ID               string `json:"id"`
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

func NewBill(params []string) *Bill {
	bill := new(Bill)
	bill.ID = params[1]
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
