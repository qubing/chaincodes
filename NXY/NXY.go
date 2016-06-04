package main

import (
	"encoding/json"
	// "errors"
	"fmt"
	// "github.com/hyperledger/fabric/core/chaincode/shim"
)

//==============================================================================================================================
//	 Structure Definitions
//==============================================================================================================================
//	Chaincode - A blank struct for use with Shim (A HyperLedger included go file used for get/put state
//				and other HyperLedger functions)
//==============================================================================================================================
type SimpleChaincode struct {
}

//==============================================================================================================================
//	Bill - Defines the structure for a car object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON make -> Struct Make.
//==============================================================================================================================
type Bill struct {
	No       string `json:"no"`
	Greeting string `json:"greeting"`
}

type Sign struct {
  StepID   string `json:"step"`
  Operator string `json:"operator"`
}

type Contract struct {
  ID    string `json:"id"`
  Signs []Sign `json:"signs"`
  Bills map[string]Bill `json:"bills"`
}

const STATE_INIT = "01"


func newContract(contractID string, userID string) *Contract {
  t := new(Contract)
  t.ID = contractID
  t.Bills = make(map[string]Bill, 0)
  t.Signs = make([]Sign, 0)
  t.addSign(STATE_INIT, userID)
  return t
}
func (t *Contract) addBill(params []string) {
  if t.Bills == nil {
    t.Bills = make(map[string]Bill, 0)
  }
  bill := Bill{params[0], params[1]}
  t.Bills[bill.No]= bill
}
func (t *Contract) addSign(state string, userID string) {
  if t.Signs == nil {
    t.Signs = make([]Sign, 0)
  }
  sign := Sign{state, userID}
  t.Signs = append(t.Signs, sign)
}

func (t *Contract) toJSON() string {
  bytes, err := json.Marshal(t)
  if err != nil {
    fmt.Printf("Error: %s\n", err)
    return "{}"
  }
  return string(bytes)
}

func main() {
  contract := newContract("0000000001", "user01")
  //contract.init("0000000001", "user01")
  fmt.Println(contract.toJSON())
  contract.addSign("02", "user02")
  fmt.Println(contract.toJSON())
  contract.addBill([]string {"01", "AB"})
  fmt.Println(contract.toJSON())
}
