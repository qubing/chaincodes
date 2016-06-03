package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
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
//	Hello - Defines the structure for a car object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON make -> Struct Make.
//==============================================================================================================================
type Hello struct {
	Name     string `json:"name"`
	Greeting string `json:"greeting"`
}

func ToHello(name string, greeting string) (Hello, error) {
	var hello Hello
	name_json := "\"name\":\"" + name + "\", "
	greeting_json := "\"greeting\":\"" + greeting + "\""
	_json := "{" + name_json + greeting_json + "}"
	fmt.Println(_json)
	err := json.Unmarshal([]byte(_json), &hello)
	if err != nil {
		return hello, errors.New("Invalid JSON object")
	}
	return hello, nil
}

func ToJSON(hello Hello) ([]byte, error) {
	bytes, err := json.Marshal(hello)
	if err != nil {
		return nil, errors.New("Invalid Hello object")
	}
	return bytes, nil
}

//==============================================================================================================================
//	Init Function - Called when the user deploys the chaincode
//==============================================================================================================================
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	return nil, nil
}

//==============================================================================================================================
//	 retrieve_v5c - Gets the state of the data at v5cID in the ledger then converts it from the stored
//					JSON into the Vehicle struct for use in the contract. Returns the Vehcile struct.
//					Returns empty v if it errors.
//==============================================================================================================================
func (t *SimpleChaincode) getString(stub *shim.ChaincodeStub, name string) ([]byte, error) {

	//var v Hello

	bytes, err := stub.GetState("STRING:" + name)

	if err != nil {
		fmt.Printf("Hello: Failed to invoke vehicle_code: %s", err)
		return nil, errors.New("Hello: Error retrieving vehicle with name = " + name)
	}

	return bytes, nil
}

//==============================================================================================================================
//	 retrieve_v5c - Gets the state of the data at v5cID in the ledger then converts it from the stored
//					JSON into the Vehicle struct for use in the contract. Returns the Vehcile struct.
//					Returns empty v if it errors.
//==============================================================================================================================
func (t *SimpleChaincode) getJSON(stub *shim.ChaincodeStub, name string) ([]byte, error) {

	//var v Hello

	bytes, err := stub.GetState("JSON:" + name)

	if err != nil {
		fmt.Printf("Hello: Failed to invoke vehicle_code: %s", err)
		return nil, errors.New("Hello: Error retrieving vehicle with name = " + name)
	}

	return bytes, nil
}

//=================================================================================================================================
//	 createX
//=================================================================================================================================
func (t *SimpleChaincode) create(stub *shim.ChaincodeStub, name string, greeting string) (bool, error) {
	var hello Hello

	// Variables to define the JSON
	name_json := "\"name\":\"" + name + "\", "
	greeting_json := "\"greeting\":\"" + greeting + "\""

	_json := "{" + name_json + greeting_json + "}"

	err := json.Unmarshal([]byte(_json), &hello)
	if err != nil {
		return false, errors.New("Invalid JSON object")
	}

	bytes, err := json.Marshal(hello)
	if err != nil {
		fmt.Printf("SAVE_CHANGES: Error converting vehicle record: %s", err)
		return false, errors.New("Error converting vehicle record")
	}

	err = stub.PutState(hello.Name, bytes)

	if err != nil {
		fmt.Printf("SAVE_CHANGES: Error storing vehicle record: %s", err)
		return false, errors.New("Error storing vehicle record")
	}

	return true, nil
}

//==============================================================================================================================
//	 Router Functions
//==============================================================================================================================
//	Invoke - Called on chaincode invoke. Takes a function name passed and calls that function. Converts some
//		  initial arguments passed to other things for use in the called function e.g. name -> ecert
//==============================================================================================================================
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "putString" {
		if len(args) != 2 {
			fmt.Printf("Incorrect number of arguments passed.")
			return nil, errors.New("create@Invoke: Incorrect number of arguments passed.")
		}
		err := stub.PutState("STRING:" + args[0], []byte(args[1]))

		if err != nil {
			return nil, errors.New("Error storing string record")
		}
		return nil, nil
	} else if function == "putJson" {
		if len(args) != 2 {
			fmt.Printf("Incorrect number of arguments passed.")
			return nil, errors.New("create@Invoke: Incorrect number of arguments passed.")
		}

		hello, err := ToHello(args[0], args[1])
		if err != nil {
			return nil, errors.New("Error storing string record")
		}
		bytes, err := ToJSON(hello)
		err = stub.PutState("JSON:" + args[0], bytes)
		if err != nil {
			return nil, errors.New("Error storing string record")
		}
		return nil, nil
	}
	return nil, errors.New("not valid invoke method")
}

//=================================================================================================================================
//	Query - Called on chaincode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
//=================================================================================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "getString" {
		bytes, err := t.getString(stub, args[0])
		if err != nil {
			return nil, errors.New("no data found.")
		}
		return bytes, nil
	} else if function == "getJson" {
		bytes, err := t.getJSON(stub, args[0])
		if err != nil {
			return nil, errors.New("no data found.")
		}
		return bytes, nil
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
}
