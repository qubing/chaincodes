package main

import (
    "errors"
    "fmt"
    "time"
    "github.com/hyperledger/fabric/core/chaincode/shim"
)

//==============================================================================
//	 Hyperledger Chaincode Method Definitions
//==============================================================================
type SimpleChaincode struct {
}

const TABLE_XX = "T_XX"
const COL_XX_1 = "COL1"
const COL_XX_2 = "COL2"
const COL_XX_3 = "COL3"

func createTable_XX(stub *shim.ChaincodeStub) error {
    var columnDefs []*shim.ColumnDefinition
	  col1 := shim.ColumnDefinition{
        Name: COL_XX_1,
		    Type: shim.ColumnDefinition_STRING,
        Key: true}
	  col2 := shim.ColumnDefinition{
        Name: COL_XX_2,
		    Type: shim.ColumnDefinition_STRING,
        Key: false}
    col3 := shim.ColumnDefinition{
        Name: COL_XX_3,
		    Type: shim.ColumnDefinition_STRING,
        Key: false}
    columnDefs = append(columnDefs, &col1)
    columnDefs = append(columnDefs, &col2)
    columnDefs = append(columnDefs, &col3)

    return stub.CreateTable(TABLE_XX, columnDefs)
}

func insertTable_XX(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    if len(args) < 3 {
        return nil, errors.New("insertTableOne failed. Must include 3 column values")
    }

    var cols []*shim.Column
		col1 := shim.Column{Value: &shim.Column_String_{String_: args[0]}}
		col2 := shim.Column{Value: &shim.Column_String_{String_: args[1]}}
		col3 := shim.Column{Value: &shim.Column_String_{String_: args[2]}}

		cols = append(cols, &col1)
		cols = append(cols, &col2)
		cols = append(cols, &col3)

		row := shim.Row{Columns: cols}
		ok, err := stub.InsertRow(TABLE_XX, row)
		if err != nil {
			   return nil, fmt.Errorf("insertTableOne operation failed. %s", err)
		}
		if !ok {
			   return nil, errors.New("insertTableOne operation failed. Row with given key already exists")
		}

    return nil, nil
}

func searchTable_XX(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var cols []shim.Column
		col1 := shim.Column{Value: &shim.Column_String_{String_: args[0]}}
		cols = append(cols, col1)

		row, err := stub.GetRow(TABLE_XX, cols)
		if err != nil {
			return nil, fmt.Errorf("getRow operation failed. %s", err)
		}

    rowString := fmt.Sprintf("%s", row)
		return []byte(rowString), nil
}

//==============================================================================================================================
//	Init Function - Called when the user deploys the chaincode
//==============================================================================================================================
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	  //init BPM templates
    //err := createTable_XX(stub)
    //if (err != nil) {
    //    return nil, err
    //}
    //init User & Roles mapping

    //init Org mapping
    //init Invest Target validations

    stub.PutState("a", []byte("aaaaaaaaaa"))

	return nil, nil
}

func (t *SimpleChaincode) createInvest(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    id := time.Now().UnixNano()
    bytes := []byte(string(id))
    return bytes, nil
}

func (t *SimpleChaincode) processInvest(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    return nil, nil
}

func (t *SimpleChaincode) createHost(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    id := time.Now().UnixNano()
    bytes := []byte(string(id))
    return bytes, nil
}

func (t *SimpleChaincode) processHost(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    return nil, nil
}

//------------------------------------------------------------------------------------------------------------------------------
//	Invoke Router - Called on chaincode invoke. Takes a function name passed and calls that function. Converts some
//		  initial arguments passed to other things for use in the called function e.g. name -> ecert
//------------------------------------------------------------------------------------------------------------------------------
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    if function == "createInvest" {
		    return t.createInvest(stub, args)
    } else if function == "processInvest" {
        return t.processInvest(stub, args)
    } else if function == "createHost" {
        return t.createHost(stub, args)
    } else if function == "processHost" {
        return t.processHost(stub, args)
    } else if function == "pushA" {
		stub.PutState("a", []byte(args[0]))
		return []byte("OK"), nil
	}

    return nil, errors.New("Received unknown function invocation.")
}

func (t *SimpleChaincode) listInvest(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    return nil, nil
}

func (t *SimpleChaincode) viewInvest(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    return nil, nil
}

func (t *SimpleChaincode) listHost(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    return nil, nil
}

func (t *SimpleChaincode) viewHost(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    return nil, nil
}

//=================================================================================================================================
//	Query - Called on chaincode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
//=================================================================================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
  	if function == "listInvest" {
        return t.listInvest(stub, args)
  	} else if function == "viewInvest" {
  		  return t.viewInvest(stub, args)
  	} else if function == "listHost" {
  		  return t.listHost(stub, args)
  	} else if function == "viewHost" {
  		  return t.viewHost(stub, args)
  	} else if function == "pushA" {
		return stub.GetState("a")
	}
  	return nil, errors.New("Received unknown function query.")
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
