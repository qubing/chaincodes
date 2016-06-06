package model

import (
	"encoding/json"
	"fmt"
)

//======================================
//	Trade - 交易信息
//======================================
//交易总账
//{
//	"no": "%交易编号%",
//	"from": "%卖方机构编号%",
//	"to": "买方机构编号",
//	"bills": [
//		{%TradeBill%},
//    {...},
//    ...
//  ]
//}
type Trade struct {
	ID    string `json:"id"`
	From  string `json:"from"`
	To    string `json:"to"`
	Signs []Sign `json:"signs"`
	Bills map[string]TradeBill `json:"bills"`
}

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
