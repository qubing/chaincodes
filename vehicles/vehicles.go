package main

import (
"errors"
"fmt"
"strconv"
"strings"
"github.com/hyperledger/fabric/core/chaincode/shim"
"encoding/json"
"crypto/x509"
"encoding/pem"
"net/url"

)

//==============================================================================================================================
//	 Participant types - Each participant type is mapped to an integer which we use to compare to the value stored in a
//						 user's eCert
//==============================================================================================================================
const   AUTHORITY      =  1
const   MANUFACTURER   =  2
const   PRIVATE_ENTITY =  3
const   LEASE_COMPANY  =  4
const   SCRAP_MERCHANT =  5


//==============================================================================================================================
//	 Status types - Asset lifecycle is broken down into 5 statuses, this is part of the business logic to determine what can
//					be done to the vehicle at points in it's lifecycle
//==============================================================================================================================
const   STATE_TEMPLATE  			=  0
const   STATE_MANUFACTURE  			=  1
const   STATE_PRIVATE_OWNERSHIP 	=  2
const   STATE_LEASED_OUT 			=  3
const   STATE_BEING_SCRAPPED  		=  4

//==============================================================================================================================
//	 Structure Definitions
//==============================================================================================================================
//	Chaincode - A blank struct for use with Shim (A HyperLedger included go file used for get/put state
//				and other HyperLedger functions)
//==============================================================================================================================
type  SimpleChaincode struct {
}

//==============================================================================================================================
//	Vehicle - Defines the structure for a car object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON make -> Struct Make.
//==============================================================================================================================
type Vehicle struct {
Make            string `json:"make"`
Model           string `json:"model"`
Reg             string `json:"reg"`
VIN             int    `json:"VIN"`
Owner           string `json:"owner"`
Scrapped        bool   `json:"scrapped"`
Status          int    `json:"status"`
Colour          string `json:"colour"`
V5cID           string `json:"v5cID"`
LeaseContractID string `json:"leaseContractID"`
}

//==============================================================================================================================
//	ECertResponse - Struct for storing the JSON response of retrieving an ECert. JSON OK -> Struct OK
//==============================================================================================================================
type ECertResponse struct {
OK string `json:"OK"`
}

//==============================================================================================================================
//	Init Function - Called when the user deploys the chaincode
//==============================================================================================================================
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	return nil, nil
}

//==============================================================================================================================
//	 General Functions
//==============================================================================================================================
//	 get_ecert - Takes the name passed and calls out to the REST API for HyperLedger to retrieve the ecert
//				 for that user. Returns the ecert as retrived including html encoding.
//==============================================================================================================================
func (t *SimpleChaincode) get_ecert(stub *shim.ChaincodeStub, name string) ([]byte, error) {
	return nil, errors.New("Error retrieving peer address")
}

//==============================================================================================================================
//	 get_caller - Retrieves the username of the user who invoked the chaincode.
//				  Returns the username as a string.
//==============================================================================================================================


func (t *SimpleChaincode) get_username(stub *shim.ChaincodeStub) (string, error) {

bytes, err := stub.GetCallerCertificate();
if err != nil { return "", errors.New("Couldn't retrieve caller certificate") }
x509Cert, err := x509.ParseCertificate(bytes);				// Extract Certificate from result of GetCallerCertificate
if err != nil { return "", errors.New("Couldn't parse certificate")	}

return x509Cert.Subject.CommonName, nil
}

//==============================================================================================================================
//	 check_affiliation - Takes an ecert as a string, decodes it to remove html encoding then parses it and checks the
// 				  		certificates common name. The affiliation is stored as part of the common name.
//==============================================================================================================================

func (t *SimpleChaincode) check_affiliation(stub *shim.ChaincodeStub, cert string) (int, error) {

decodedCert, err := url.QueryUnescape(cert);    				// make % etc normal //

if err != nil { return -1, errors.New("Could not decode certificate") }

pem, _ := pem.Decode([]byte(decodedCert))           				// Make Plain text   //

x509Cert, err := x509.ParseCertificate(pem.Bytes);				// Extract Certificate from argument //

if err != nil { return -1, errors.New("Couldn't parse certificate")	}

cn := x509Cert.Subject.CommonName

res := strings.Split(cn,"\\")

affiliation, _ := strconv.Atoi(res[2])

return affiliation, nil
}

//==============================================================================================================================
//	 get_caller_data - Calls the get_ecert and check_role functions and returns the ecert and role for the
//					 name passed.
//==============================================================================================================================

func (t *SimpleChaincode) get_caller_data(stub *shim.ChaincodeStub) (string, int, error){

user, err := t.get_username(stub)
if err != nil { return "", -1, err }

ecert, err := t.get_ecert(stub, user);
if err != nil { return "", -1, err }

affiliation, err := t.check_affiliation(stub,string(ecert));
if err != nil { return "", -1, err }

return user, affiliation, nil
}

//==============================================================================================================================
//	 retrieve_v5c - Gets the state of the data at v5cID in the ledger then converts it from the stored
//					JSON into the Vehicle struct for use in the contract. Returns the Vehcile struct.
//					Returns empty v if it errors.
//==============================================================================================================================
func (t *SimpleChaincode) retrieve_v5c(stub *shim.ChaincodeStub, v5cID string) (Vehicle, error) {

var v Vehicle

bytes, err := stub.GetState(v5cID)	;

if err != nil {	fmt.Printf("RETRIEVE_V5C: Failed to invoke vehicle_code: %s", err); return v, errors.New("RETRIEVE_V5C: Error retrieving vehicle with v5cID = " + v5cID) }

err = json.Unmarshal(bytes, &v)	;

if err != nil {	fmt.Printf("RETRIEVE_V5C: Corrupt vehicle record "+string(bytes)+": %s", err); return v, errors.New("RETRIEVE_V5C: Corrupt vehicle record"+string(bytes))	}

return v, nil
}

//==============================================================================================================================
// save_changes - Writes to the ledger the Vehicle struct passed in a JSON format. Uses the shim file's
//				  method 'PutState'.
//==============================================================================================================================
func (t *SimpleChaincode) save_changes(stub *shim.ChaincodeStub, v Vehicle) (bool, error) {

bytes, err := json.Marshal(v)

if err != nil { fmt.Printf("SAVE_CHANGES: Error converting vehicle record: %s", err); return false, errors.New("Error converting vehicle record") }

err = stub.PutState(v.V5cID, bytes)

if err != nil { fmt.Printf("SAVE_CHANGES: Error storing vehicle record: %s", err); return false, errors.New("Error storing vehicle record") }

return true, nil
}

//==============================================================================================================================
//	 Router Functions
//==============================================================================================================================
//	Invoke - Called on chaincode invoke. Takes a function name passed and calls that function. Converts some
//		  initial arguments passed to other things for use in the called function e.g. name -> ecert
//==============================================================================================================================
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	return nil, errors.New("Error retrieving caller information")
}
//=================================================================================================================================
//	Query - Called on chaincode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
//=================================================================================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	return nil, errors.New("Received unknown function invocation")
}

//=================================================================================================================================
//	 Create Function
//=================================================================================================================================
//	 Create Vehicle - Creates the initial JSON for the vehcile and then saves it to the ledger.
//=================================================================================================================================
func (t *SimpleChaincode) create_vehicle(stub *shim.ChaincodeStub, caller string, caller_affiliation int, v5cID string) ([]byte, error) {
	return nil, nil
}

//=================================================================================================================================
//	 Transfer Functions
//=================================================================================================================================
//	 authority_to_manufacturer
//=================================================================================================================================
func (t *SimpleChaincode) authority_to_manufacturer(stub *shim.ChaincodeStub, v Vehicle, caller string, caller_affiliation int, recipient_name string, recipient_affiliation int) ([]byte, error) {
	return nil, nil									// We are Done
}

//=================================================================================================================================
//	 manufacturer_to_private
//=================================================================================================================================
func (t *SimpleChaincode) manufacturer_to_private(stub *shim.ChaincodeStub, v Vehicle, caller string, caller_affiliation int, recipient_name string, recipient_affiliation int) ([]byte, error) {
	return nil, nil
}

//=================================================================================================================================
//	 private_to_private
//=================================================================================================================================
func (t *SimpleChaincode) private_to_private(stub *shim.ChaincodeStub, v Vehicle, caller string, caller_affiliation int, recipient_name string, recipient_affiliation int) ([]byte, error) {
	return nil, nil
}

//=================================================================================================================================
//	 private_to_lease_company
//=================================================================================================================================
func (t *SimpleChaincode) private_to_lease_company(stub *shim.ChaincodeStub, v Vehicle, caller string, caller_affiliation int, recipient_name string, recipient_affiliation int) ([]byte, error) {
	return nil, nil
}

//=================================================================================================================================
//	 lease_company_to_private
//=================================================================================================================================
func (t *SimpleChaincode) lease_company_to_private(stub *shim.ChaincodeStub, v Vehicle, caller string, caller_affiliation int, recipient_name string, recipient_affiliation int) ([]byte, error) {
	return nil, nil
}

//=================================================================================================================================
//	 private_to_scrap_merchant
//=================================================================================================================================
func (t *SimpleChaincode) private_to_scrap_merchant(stub *shim.ChaincodeStub, v Vehicle, caller string, caller_affiliation int, recipient_name string, recipient_affiliation int) ([]byte, error) {
	return nil, nil
}

//=================================================================================================================================
//	 Update Functions
//=================================================================================================================================
//	 update_vin
//=================================================================================================================================
func (t *SimpleChaincode) update_vin(stub *shim.ChaincodeStub, v Vehicle, caller string, caller_affiliation int, new_value string) ([]byte, error) {
	return nil, nil
}


//=================================================================================================================================
//	 update_registration
//=================================================================================================================================
func (t *SimpleChaincode) update_registration(stub *shim.ChaincodeStub, v Vehicle, caller string, caller_affiliation int, new_value string) ([]byte, error) {
	return nil, nil
}

//=================================================================================================================================
//	 update_make
//=================================================================================================================================
func (t *SimpleChaincode) update_make(stub *shim.ChaincodeStub, v Vehicle, caller string, caller_affiliation int, new_value string) ([]byte, error) {
	return nil, nil
}

//=================================================================================================================================
//	 update_model
//=================================================================================================================================
func (t *SimpleChaincode) update_model(stub *shim.ChaincodeStub, v Vehicle, caller string, caller_affiliation int, new_value string) ([]byte, error) {
	return nil, nil
}

//=================================================================================================================================
//	 scrap_vehicle
//=================================================================================================================================
func (t *SimpleChaincode) scrap_vehicle(stub *shim.ChaincodeStub, v Vehicle, caller string, caller_affiliation int) ([]byte, error) {
	return nil, nil
}

//=================================================================================================================================
//	 Main - main - Starts up the chaincode
//=================================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))

	if err != nil { fmt.Printf("Error starting Chaincode: %s", err) }
}

