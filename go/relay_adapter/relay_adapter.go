package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const (
	PREFIX_TOPIC_SENDER   = "TOPIC_SEND_"
	PREFIX_TOPIC_RECEIVER = "TOPIC_RECV_"
)

// RelayAdapter - relay adapter for communication cross blockchain network
type RelayAdapter struct {
}

func main() {
	err := shim.Start(new(RelayAdapter))
	if err != nil {
		fmt.Printf("Error starting RelayAdapter: %s", err)
		fmt.Println()
	}
}

// Init - Initializing smart contract
func (t *RelayAdapter) Init(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, _ := stub.GetFunctionAndParameters()
	if funcName == "init" {
		// init deployment of smart contract
		return t.init(stub)
	} else if funcName == "upgrade" {
		// upgrade deployment of smart contract
		return t.upgrade(stub)
	}

	return shim.Error(fmt.Sprintf(`Invalid function name for deployment. 
		(expecting 'init' or 'upgrade', actual: '%s')`, funcName))
}

func (t *RelayAdapter) init(stub shim.ChaincodeStubInterface) pb.Response {
	// TODO:
	return shim.Success(nil)
}

func (t *RelayAdapter) upgrade(stub shim.ChaincodeStubInterface) pb.Response {
	// TODO:
	return shim.Success(nil)
}

// Invoke - accessing smart contract interface(query or invoke)
func (t *RelayAdapter) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, params := stub.GetFunctionAndParameters()
	if funcName == "topic" {
		// query installed topics (IN or OUT)
		return t.findTopic(stub, params)
	} else if funcName == "topics" {
		// list installed topics (IN or OUT)
		return t.listTopics(stub, params)
	} else if funcName == "messages" {
		// list message (IN or OUT)
		return t.listMessages(stub, params)
	} else if funcName == "install" {
		// install relay topic (IN or OUT)
		return t.install(stub, params)
	} else if funcName == "uninstall" {
		// uninstall relay topic (IN or OUT)
		return t.uninstall(stub, params)
	} else if funcName == "register" {
		// register relay topic (IN or OUT)
		return t.register(stub, params)
	} else if funcName == "send" {
		// send message to topic (IN or OUT)
		return t.send(stub, params)
	} else if funcName == "read" {
		return t.readMessage(stub, params)
	}
	return shim.Success(nil)
}

func (t *RelayAdapter) findTopic(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	if len(params) != 2 {
		return shim.Error(fmt.Sprintf(`Incorrect number of arguments. (expecting 2, actual: "%d")`, len(params)))
	}

	topicType := params[0]
	topicName := params[1]

	var topic *Topic
	var err error
	if topicType == string(OUT) {
		topic, err = findTopic(stub, DOC_TOPIC_OUT, topicName)
	} else if topicType == string(IN) {
		topic, err = findTopic(stub, DOC_TOPIC_IN, topicName)
	} else {
		return shim.Error("topic type is wrong")
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	bytes, err := json.Marshal(topic)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(bytes)
}

func (t *RelayAdapter) listTopics(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	if len(params) != 2 {
		return shim.Error(fmt.Sprintf(`Incorrect number of arguments. (expecting 2, actual: "%d")`, len(params)))
	}

	topicType := params[0]
	orgID := params[1]

	var topics []*Topic
	var err error
	if topicType == string(OUT) {
		topics, err = queryTopicsByOrg(stub, DOC_TOPIC_OUT, orgID)
	} else if topicType == string(IN) {
		topics, err = queryTopicsByOrg(stub, DOC_TOPIC_IN, orgID)
	} else {
		return shim.Error("topic type is wrong")
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	bytes, err := json.Marshal(topics)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(bytes)
}

func (t *RelayAdapter) readMessage(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	if len(params) != 3 {
		return shim.Error(fmt.Sprintf(`Incorrect number of arguments. (expecting 3, actual: "%d")`, len(params)))
	}

	messageType := params[0]
	messageID := params[1]
	orgID := params[2]

	var message *Message
	var err error
	var topicDocType, messageDocType DocumentType
	if messageType == string(OUT) {
		messageDocType = DOC_MESSAGE_OUT
		topicDocType = DOC_TOPIC_OUT
	} else if messageType == string(IN) {
		messageDocType = DOC_MESSAGE_IN
		topicDocType = DOC_TOPIC_IN
	} else {
		return shim.Error("message type is wrong")
	}

	message, err = findMessage(stub, messageDocType, messageID)

	if err != nil {
		return shim.Error(err.Error())
	}

	if message == nil {
		return shim.Error(fmt.Sprintf(`message not found.(id: %s)`, messageID))
	}

	topic, err := findTopic(stub, topicDocType, message.Name)

	if err != nil {
		return shim.Error(err.Error())
	}

	decoded, err := topic.DecryptMessage(orgID, message.Body)

	if err != nil {
		return shim.Error(err.Error())
	}
	message.Body = decoded

	bytes, err := json.Marshal(message)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(bytes)
}

func (t *RelayAdapter) listMessages(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	if len(params) != 2 {
		return shim.Error(fmt.Sprintf(`Incorrect number of arguments. (expecting 2, actual: "%d")`, len(params)))
	}

	messageType := params[0]
	topicName := params[1]

	var messages []*Message
	var err error
	if messageType == string(OUT) {
		messages, err = queryMessagesByTopic(stub, DOC_MESSAGE_OUT, topicName)
	} else if messageType == string(IN) {
		messages, err = queryMessagesByTopic(stub, DOC_MESSAGE_IN, topicName)
	} else {
		return shim.Error("topic type is wrong")
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	bytes, err := json.Marshal(messages)

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(bytes)
}

// install topic for sender|receiver
func (t *RelayAdapter) install(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	if len(params) != 2 {
		return shim.Error(fmt.Sprintf(`Incorrect number of arguments. (expecting 2, actual: "%d")`, len(params)))
	}

	topicType := params[0]
	topicName := params[1]

	var topic *Topic
	var err error
	if topicType == string(IN) {
		topic, err = findTopic(stub, DOC_TOPIC_IN, topicName)
	} else if topicType == string(OUT) {
		topic, err = findTopic(stub, DOC_TOPIC_OUT, topicName)
	} else {
		return shim.Error("topic type is wrong")
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	if topic == nil {
		if topicType == string(IN) {
			topic = NewTopic(DOC_TOPIC_IN, topicName)
		} else if topicType == string(OUT) {
			topic = NewTopic(DOC_TOPIC_OUT, topicName)
		} else {
			return shim.Error("topic type is wrong")
		}
		topic.Id = stub.GetTxID()
		data, err := ToJSON(topic)
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Printf(`Putting state '%s'`, data)
		fmt.Println()
		err = stub.PutState(topic.Id, []byte(data))
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(nil)
	}

	return shim.Error(fmt.Sprintf(`topic already existing. (topic: %s)`, topic.Name))
}

func (t *RelayAdapter) register(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	if len(params) != 6 {
		return shim.Error(fmt.Sprintf(`Incorrect number of arguments. (expecting 6, actual: "%d")`, len(params)))
	}

	topicType := params[0]
	topicName := params[1]
	registerType := params[2]
	orgID := params[3]
	publicKey := params[4]
	privateKey := params[5]

	var topic *Topic
	var err error
	if topicType == string(IN) {
		topic, err = findTopic(stub, DOC_TOPIC_IN, topicName)
	} else if topicType == string(OUT) {
		topic, err = findTopic(stub, DOC_TOPIC_OUT, topicName)
	} else {
		return shim.Error("topic type is wrong")
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	if topic == nil {
		return shim.Error(fmt.Sprintf(`topic not existing. (topic: %s)`, topicName))
	}

	if topicType == string(IN) && topic.SenderExist(orgID) {
		return shim.Error(fmt.Sprintf(`topic already registered. (topic: %s ,org: %s)`, topic.Name, orgID))
	} else if topicType == string(OUT) && topic.SenderExist(orgID) {
		return shim.Error(fmt.Sprintf(`topic already registered. (topic: %s ,org: %s)`, topic.Name, orgID))
	}

	if registerType == string(REGISTER_SENDER) {
		if topicType == string(IN) && topic.SenderExist(orgID) {
			return shim.Error(fmt.Sprintf(`topic already registered. (topic: %s ,org: %s)`, topic.Name, orgID))
		} else if topicType == string(OUT) && topic.SenderExist(orgID) {
			return shim.Error(fmt.Sprintf(`topic already registered. (topic: %s ,org: %s)`, topic.Name, orgID))
		}
		topic.AddSender(orgID, publicKey)
	} else if registerType == string(REGISTER_READER) {
		if topicType == string(IN) && topic.ReaderExist(orgID) {
			return shim.Error(fmt.Sprintf(`topic already registered. (topic: %s ,org: %s)`, topic.Name, orgID))
		} else if topicType == string(OUT) && topic.ReaderExist(orgID) {
			return shim.Error(fmt.Sprintf(`topic already registered. (topic: %s ,org: %s)`, topic.Name, orgID))
		}
		topic.AddReader(orgID, publicKey, privateKey)
	} else {
		return shim.Error("register type is wrong")
	}

	data, err := ToJSON(topic)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf(`Putting state '%s'`, data)
	fmt.Println()
	err = stub.PutState(topic.Id, []byte(data))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *RelayAdapter) uninstall(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	// TODO:
	return shim.Success(nil)
}

func (t *RelayAdapter) send(stub shim.ChaincodeStubInterface, params []string) pb.Response {
	if len(params) != 3 {
		return shim.Error(fmt.Sprintf(`Incorrect number of arguments. (expecting 3, actual: "%d")`, len(params)))
	}

	topicType := params[0]
	messageJSON := params[1]
	message := new(Message)
	err := json.Unmarshal([]byte(messageJSON), &message)
	if err != nil {
		return shim.Error(err.Error())
	}
	//fmt.Printf(`message:> %s`, messageJSON)
	//data, _ := ToJSON(message)
	//return shim.Success([]byte(data))

	routeJSON := params[2]
	route := new(Route)
	err = json.Unmarshal([]byte(routeJSON), route)
	if err != nil {
		return shim.Error(err.Error())
	}
	// fmt.Printf(`route:> %s`, routeJSON)
	// data, _ := ToJSON(route)
	// return shim.Success([]byte(data))

	var topic *Topic
	if topicType == string(OUT) {
		topic, err = findTopic(stub, DOC_TOPIC_OUT, message.Name)
		if err != nil {
			return shim.Error(err.Error())
		} else if topic == nil {
			return shim.Error("send message(OUT) failed, cause: topic not found")
		} else if topic.SenderExist(route.OrgID) == false {
			return shim.Error(fmt.Sprintf(`send message(OUT) failed, cause: sender not registered.(org:%s)`, route.OrgID))
		}
		message.Body, err = topic.EncryptMessage(route.OrgID, message.Body)
		if err != nil {
			return shim.Error(fmt.Sprintf(`send message failed, cause: %s`, err.Error()))
		}
		message.DocType = DOC_MESSAGE_OUT
	} else if topicType == string(IN) {
		topic, err = findTopic(stub, DOC_TOPIC_IN, message.Name)

		if err != nil {
			return shim.Error(err.Error())
		} else if topic == nil {
			return shim.Error("send message(IN) failed, cause: topic not found")
		} else if topic.SenderExist(route.OrgID) == false {
			return shim.Error(fmt.Sprintf(`send message(IN) failed, cause: sender not registered.(topic:%s, org:%s)`, topic.Name, route.OrgID))
		}
		message.DocType = DOC_MESSAGE_IN
	} else {
		return shim.Error("topic type is wrong")
	}
	message.Id = stub.GetTxID()
	route.TrxID = stub.GetTxID()
	message.Routes = append(message.Routes, route)

	//encryption failed.
	if err != nil {
		return shim.Error(err.Error())
	}

	data, err := ToJSON(message)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Printf(`Putting state '%s'`, data)
	fmt.Println()
	err = stub.PutState(message.Id, []byte(data))
	if err != nil {
		return shim.Error(err.Error())
	}

	if topicType == string(OUT) {
		//send message to event hub
		fmt.Printf(`send message: |%s|`, data)
		err = stub.SetEvent(message.Name, []byte(data))
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	return shim.Success(nil)
}
