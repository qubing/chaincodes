package model

//=====================================
//	Sign - 审批信息
//=====================================
type Sign struct {
	StepID   string `json:"step"`
	Operator string `json:"operator"`
  Comments string `json:"comments"`
}
