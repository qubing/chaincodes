package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// BalanceManager Smart Contract(Chaincode) implementation
type BalanceManager struct {
}

func (t *BalanceManager) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("<BalanceMgr Init>")

	return shim.Success(nil)
}

func (t *BalanceManager) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("<BalanceMgr Invoke>")

	funcName, args := stub.GetFunctionAndParameters()
	if funcName == "create" {
		// Create account with balance of '0'
		return t.create(stub, args)
	} else if funcName == "charge" {
		// Charge account with amount
		return t.charge(stub, args)
	} else if funcName == "transfer" {
		// Transfer A to B with some money
		return t.transfer(stub, args)
	} else if funcName == "query" {
		// Query account balance
		return t.query(stub, args)
	}

	return shim.Error(fmt.Sprintf(`Invalid invoke function name. Expecting "create" "transfer" "query". Actual: "%s"`, funcName))
}

// create: create account initialized with 0
func (t *BalanceManager) create(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("create account with initial balance of '0'")

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	accountName := args[0]

	valBytes,err := stub.GetState(accountName)
	if err != nil {
		return shim.Error("Account create failed with unknown reason.")
	}

	if valBytes != nil && len(valBytes) > 0 {
		return shim.Error(fmt.Sprintf(`Account already existed. (Account: "%s")`, accountName))
	}

	stub.PutState(accountName, []byte(strconv.Itoa(0)))

	byts, err := stub.GetCreator()

	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf(`Printing current user is "%s".`, string(byts))

	return shim.Success(nil)
}

func (t *BalanceManager) charge(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("charge account with amount")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	accountFrom := args[0]

	bytesFrom, err := stub.GetState(accountFrom)
	if err != nil {
		return shim.Error("Account charge failed with unknown reason.")
	}
	if bytesFrom == nil {
		return shim.Error("Account not found")
	}
	valFrom, _ := strconv.Atoi(string(bytesFrom))

	// Perform the execution
	amountCharge, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	valFrom = valFrom + amountCharge

	stub.PutState(accountFrom, []byte(strconv.Itoa(valFrom)))

	stub.SetEvent("hello", []byte(fmt.Sprintf(`{"status":"ok", "account":"%s", "amount":"%d"}`, accountFrom, amountCharge)))

	return shim.Success(nil)
}

// transfer: Transfer balance from account to another
func (t *BalanceManager) transfer(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("transfer account")

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	accountFrom := args[0]
	accountTo := args[1]

	// Get the state from the ledger
	bytesFrom, err := stub.GetState(accountFrom)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if bytesFrom == nil {
		return shim.Error("Entity not found")
	}
	valFrom, _ := strconv.Atoi(string(bytesFrom))

	bytesTo, err := stub.GetState(accountTo)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if bytesTo == nil {
		return shim.Error("Entity not found")
	}
	valTo, _ := strconv.Atoi(string(bytesTo))

	// Perform the execution
	amountTransfer, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	valFrom = valFrom - amountTransfer
	valTo = valTo + amountTransfer
	fmt.Printf("valFrom = %d, valTo = %d\n", valFrom, valTo)

	// Write the state back to the ledger
	err = stub.PutState(accountFrom, []byte(strconv.Itoa(valFrom)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(accountTo, []byte(strconv.Itoa(valTo)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// query: query account balance
func (t *BalanceManager) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	fmt.Println()

	byts, err := stub.GetCreator()

	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf(`Printing current user is "%s".`, string(byts))
	fmt.Println()

	id, err := cid.GetID(stub)

	if err != nil {
		fmt.Printf(`GetID failed. Error: "%s".`, err.Error())
	} else {
		fmt.Printf(`GetID="%s".`, id)
	}
	fmt.Println()

	mspId, err := cid.GetMSPID(stub)
	if err != nil {
		fmt.Printf(`GetMSPID failed. Error: "%s".`, err.Error())
	} else {
		fmt.Printf(`GetMSPID="%s".`, mspId)
	}
	fmt.Println()

	email_val, email_found, err := cid.GetAttributeValue(stub, "email")
	if err != nil {
		fmt.Printf(`GetAttributeValue('email') failed. Error: "%s".`, err.Error())
	} else if !email_found {
		fmt.Printf(`GetAttributeValue('email') not found. Value: "%s".`, email_val)
	} else {
		fmt.Printf(`GetAttributeValue('email') successfully. Value: "%s".`, email_val)
	}
	fmt.Println()

	admin_val, admin_found, err := cid.GetAttributeValue(stub, "admin")
	if err != nil {
		fmt.Printf(`GetAttributeValue('admin') failed. Error: "%s".`, err.Error())
	} else if !admin_found {
		fmt.Printf(`GetAttributeValue('admin') not found. Value: "%s".`, admin_val)
	} else {
		fmt.Printf(`GetAttributeValue('admin') successfully. Value: "%s".`, admin_val)
	}
	fmt.Println()

	err = cid.AssertAttributeValue(stub, "admin", "true")
	if err != nil {
		fmt.Printf(`AssertAttributeValue('admin', 'true') failed. Error: "%s".`, err.Error())
	} else {
		fmt.Printf(`AssertAttributeValue('admin', 'true') successfully.`)
	}
	fmt.Println()

	cert, err := cid.GetX509Certificate(stub)
	if err != nil {
		fmt.Printf(`GetX509Certificate failed. Error: "%s".`, err.Error())
	} else {
		fmt.Printf(`GetX509Certificate="%s".`, cert.Raw)
	}

	fmt.Println()

	return shim.Success(Avalbytes)
}

func main() {
	err := shim.Start(new(BalanceManager))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
