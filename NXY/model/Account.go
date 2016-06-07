package model

//=====================================
//	Account - 账户信息
//=====================================
//{"no": "%机构编号%", "amount":"1,000,000", "bills": ["%票据编号%", ...]}
type Account struct {
	No     string `json:"no"`
	Amount string `json:"amount"`
	Bills  []string `json:"bills"`
}

func NewAccount(no string, amount string) *Account {
	account := new(Account)
	account.No = no
	account.Amount = amount
	account.Bills = make([]string, 0)
	return account
}
