package main

import (
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"encoding/json"
)

// Chaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Hello struct {
	Name string `json:"name"`
	Greeting string `json:"greeting"`
}

//==============================================================================================================================
//	Init Function - Called when the user deploys the chaincode
//==============================================================================================================================
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	return t.create(stub, args[0], args[1])
}

//==============================================================================================================================
//	 Router Functions
//==============================================================================================================================
//	Invoke - Called on chaincode invoke. Takes a function name passed and calls that function. Converts some
//		  initial arguments passed to other things for use in the called function e.g. name -> ecert
//==============================================================================================================================
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "create" {
		if len(args) != 2 { fmt.Printf("Incorrect number of arguments passed."); return nil, errors.New("create@Invoke: Incorrect number of arguments passed.") }
		t.create(stub, args[0], args[1]);
	}

	return nil, nil
}

//=================================================================================================================================
//	Query - Called on chaincode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
//=================================================================================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	if len(args) != 1 { fmt.Printf("Incorrect number of arguments passed"); return nil, errors.New("QUERY: Incorrect number of arguments passed") }

	if function == "query" {
		return t.retrieve(stub, args[0]), nil
	}
	return nil, errors.New("Received unknown function invocation")
}

func (t *SimpleChaincode) retrieve(stub *shim.ChaincodeStub, name string) (Hello, error) {
	var hello Hello

	bytes, err := stub.GetState(name)
	err = json.Unmarshal(bytes, &hello)	;

	if err != nil {	fmt.Printf("Retrieve: Corrupt hello record "+string(bytes)+": %s", err); return hello, errors.New("Retrieve: Corrupt hello record"+string(bytes))	}

	return hello, nil
}

func (t *SimpleChaincode) create(stub *shim.ChaincodeStub, name string, greeting string) (bool, error) {
	var hello Hello

	// Variables to define the JSON
	name_json         := "\"name\":\""+name+"\", "
	greeting_json     := "\"greeting\":"+greeting

	_json := "{" + name_json + greeting_json + "}"

	err := json.Unmarshal([]byte(_json), &hello)
	if err != nil { return nil, errors.New("Invalid JSON object") }

	t.save(stub, hello)

	return nil, nil
}

func (t *SimpleChaincode) save(stub *shim.ChaincodeStub, hello Hello) (bool, error) {
	bytes, err := json.Marshal(hello)

	if err != nil { fmt.Printf("SAVE_CHANGES: Error converting vehicle record: %s", err); return false, errors.New("Error converting vehicle record") }

	err = stub.PutState(hello.Name, bytes)

	if err != nil { fmt.Printf("SAVE_CHANGES: Error storing vehicle record: %s", err); return false, errors.New("Error storing vehicle record") }

	return true, nil
}

//=================================================================================================================================
//	 Main - main - Starts up the chaincode
//=================================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

