package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func findTopic(stub shim.ChaincodeStubInterface, docType DocumentType, topicName string) (*Topic, error) {
	selector := fmt.Sprintf(`{
		"selector": {
			"$and": [
				{"doc_type": {"$eq": "%s"}}, 
				{"name": {"$eq": "%s"}}
			]
		}
	}`, docType, topicName)
	queryIt, err := stub.GetQueryResult(selector)
	if err != nil {
		return nil, err
	}

	defer queryIt.Close()

	if queryIt.HasNext() {
		queryResult, err := queryIt.Next()
		if err != nil {
			return nil, err
		}
		topic := Topic{}
		ParseJSON(&topic, string(queryResult.GetValue()))
		return &topic, nil
	}

	return nil, nil
}

func queryTopicsByOrg(stub shim.ChaincodeStubInterface, docType DocumentType, orgID string) ([]*Topic, error) {
	topics := make([]*Topic, 0)

	selector := fmt.Sprintf(`{
		"selector": {
			"$and": [
				{
					"doc_type":{"$eq":"%s"}
				}, {
					"senders":{
						"$elemMatch":{"org_id":"%s"}
					}
				}
			]
		}
	}`, docType, orgID)
	queryIt, err := stub.GetQueryResult(selector)
	if err != nil {
		fmt.Printf(`query data failed, error ignored. error: '%s'`, err.Error())
		fmt.Println()
		return topics, err
	}

	defer queryIt.Close()

	for queryIt.HasNext() {
		queryResult, err := queryIt.Next()
		if err != nil {
			fmt.Printf(`query data failed. error: '%s'`, err.Error())
			fmt.Println()
			return nil, err
		}
		topic := Topic{}
		ParseJSON(&topic, string(queryResult.GetValue()))
		topics = append(topics, &topic)
	}
	return topics, nil
}

func queryMessagesByTopic(stub shim.ChaincodeStubInterface, docType DocumentType, topicName string) ([]*Message, error) {
	messages := make([]*Message, 0)

	selector := fmt.Sprintf(`{
		"selector": {
			"$and": [
				{
					"doc_type":{"$eq":"%s"}
				}, {
					"name":{"$eq":"%s"}
				}
			]
		}
	}`, docType, topicName)
	queryIt, err := stub.GetQueryResult(selector)
	if err != nil {
		fmt.Printf(`query data failed, error ignored. error: '%s'`, err.Error())
		fmt.Println()
		return messages, err
	}

	defer queryIt.Close()

	for queryIt.HasNext() {
		queryResult, err := queryIt.Next()
		if err != nil {
			fmt.Printf(`query data failed. error: '%s'`, err.Error())
			fmt.Println()
			return nil, err
		}
		message := Message{}
		ParseJSON(&message, string(queryResult.GetValue()))
		messages = append(messages, &message)
	}
	return messages, nil
}

func findMessage(stub shim.ChaincodeStubInterface, docType DocumentType, messageID string) (*Message, error) {
	selector := fmt.Sprintf(`{
		"selector": {
			"$and": [
				{"doc_type": {"$eq": "%s"}}, 
				{"id": {"$eq": "%s"}}
			]
		}
	}`, docType, messageID)
	queryIt, err := stub.GetQueryResult(selector)
	if err != nil {
		return nil, err
	}

	defer queryIt.Close()

	if queryIt.HasNext() {
		queryResult, err := queryIt.Next()
		if err != nil {
			return nil, err
		}
		message := Message{}
		ParseJSON(&message, string(queryResult.GetValue()))
		return &message, nil
	}

	return nil, nil
}
