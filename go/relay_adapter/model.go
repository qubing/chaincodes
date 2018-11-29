package main

import (
	"encoding/json"
	"fmt"

	"github.com/chaincodes/common/crypto"
)

type DocumentType string

const (
	DOC_TOPIC_IN    DocumentType = "TOPIC_IN"
	DOC_MESSAGE_IN  DocumentType = "MESSAGE_IN"
	DOC_TOPIC_OUT   DocumentType = "TOPIC_OUT"
	DOC_MESSAGE_OUT DocumentType = "MESSAGE_OUT"
)

type TopicType string

const (
	IN  TopicType = "IN"
	OUT TopicType = "OUT"
)

type RegisterType string

const (
	REGISTER_SENDER RegisterType = "SENDER"
	REGISTER_READER RegisterType = "READER"
)

// JSONModel - xx
type JSONModel interface {
	ParseJSON(dataJSON string) error
}

//ToJSON - generate json from object
func ToJSON(dataObject JSONModel) (string, error) {
	bytes, err := json.Marshal(dataObject)
	if err != nil {
		fmt.Printf("[ToJSON]%s", err.Error())
		return "{}", err
	}
	return string(bytes), nil
}

// ParseJSON - parse json text to object
func ParseJSON(dataObject JSONModel, dataJSON string) error {
	err := json.Unmarshal([]byte(dataJSON), &dataObject)
	if err != nil {
		fmt.Printf("[ParseJSON]%s", err.Error())
		return err
	}
	return nil
}

// ParseJSON - parse json text to object
func ParseJSONs(dataObjects []JSONModel, dataJSON string) error {
	err := json.Unmarshal([]byte(dataJSON), &dataObjects)
	if err != nil {
		fmt.Printf("[ParseJSONs]%s", err.Error())
		return err
	}
	return nil
}

type AbstractDoc struct {
	DocType DocumentType `json:"doc_type"`
	Name    string       `json:"name"`
	Id      string       `json:"id"`
}

type Topic struct {
	AbstractDoc
	Senders []*Sender `json:"senders"`
	Readers []*Reader `json:"readers"`
}

type Sender struct {
	OrgID     string `json:"org_id"`
	PublicKey string `json:"public_key"`
}

type Reader struct {
	OrgID      string `json:"org_id"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

type Message struct {
	AbstractDoc
	Body   string   `json:"body"`
	Routes []*Route `json:"routes"`
}

type Route struct {
	NetworkID  string `json:"network"`
	TrxID      string `json:"trx_id"`
	OrgID      string `json:"org_id"`
	UserID     string `json:"user_id"`
	UpdateTime string `json:"update_time"`
	Comment    string `json:"comment"`
}

func (t *Topic) ParseJSON(dataJSON string) error {
	return ParseJSON(t, dataJSON)
}

func (t *Message) ParseJSON(dataJSON string) error {
	return ParseJSON(t, dataJSON)
}

func (t *Route) ParseJSON(dataJSON string) error {
	return ParseJSON(t, dataJSON)
}

func (t *Topic) SenderExist(orgID string) bool {
	for _, sender := range t.Senders {
		if sender.OrgID == orgID {
			return true
		}
	}
	return false
}

func (t *Topic) ReaderExist(orgID string) bool {
	for _, reader := range t.Readers {
		if reader.OrgID == orgID {
			return true
		}
	}
	return false
}

func NewTopic(topicType DocumentType, topicName string) *Topic {
	topic := Topic{}
	topic.DocType = topicType
	topic.Name = topicName
	topic.Senders = make([]*Sender, 0)
	topic.Readers = make([]*Reader, 0)
	return &topic
}

func (t *Topic) NewMessage() *Message {
	message := new(Message)
	if t.DocType == DOC_TOPIC_IN {
		message.DocType = DOC_MESSAGE_IN
	} else {
		message.DocType = DOC_MESSAGE_OUT
	}

	return message
}

func NewSender(orgID string, publicKey string) *Sender {
	org := Sender{OrgID: orgID, PublicKey: publicKey}
	return &org
}

func NewReader(orgID string, publicKey string, privateKey string) *Reader {
	org := Reader{OrgID: orgID, PublicKey: publicKey, PrivateKey: privateKey}
	return &org
}

func (t *Topic) GetSender(orgID string) (*Sender, error) {
	for _, sender := range t.Senders {
		if sender.OrgID == orgID {
			return sender, nil
		}
	}
	return nil, fmt.Errorf(`sender not found. (orgID:%s)`, orgID)
}

func (t *Topic) GetReader(orgID string) (*Reader, error) {
	for _, reader := range t.Readers {
		if reader.OrgID == orgID {
			return reader, nil
		}
	}
	return nil, fmt.Errorf(`reader not found. (orgID:%s)`, orgID)
}

// func (t *Topic) NewMessage(topicName, messageBody, networkID, orgID, userID, transactionID, updateTime, comment string, histories []*Route) (*Message, error) {
// 	message := Message{}
// 	message.DocType = "MESSAGE"
// 	message.Name = topicName
// 	sender, err := t.GetSender(orgID)
// 	if err != nil {
// 		shim.Error(err.Error())
// 	}
// 	helper, errs := crypto.NewRSAHelper([]byte(sender.PublicKey), nil)
// 	if errs != nil && len(errs) > 0 {
// 		return nil, errs[0]
// 	}
// 	encoded, err := helper.Encrypt(messageBody)
// 	if err != nil {
// 		return nil, fmt.Errorf(`message encryption failed. cause: %s`, err.Error())
// 	}

// 	message.Body = encoded
// 	message.AddHistories(histories)
// 	message.AddRoute(networkID, orgID, userID, transactionID, updateTime, comment)
// 	return &message, nil
// }

func (t *Topic) EncryptMessage(orgID string, message string) (string, error) {
	sender, err := t.GetSender(orgID)
	if err != nil {
		return "", err
	}

	helper, errs := crypto.NewRSAHelper([]byte(sender.PublicKey), nil)
	if errs != nil && len(errs) > 0 {
		return "", errs[0]
	}
	encoded, err := helper.Encrypt(message)
	if err != nil {
		return "", fmt.Errorf(`message encryption failed. cause: %s`, err.Error())
	}

	return encoded, nil
}

func (t *Topic) DecryptMessage(orgID string, message string) (string, error) {
	receiver, err := t.GetReader(orgID)
	if err != nil {
		return "", err
	}

	helper, errs := crypto.NewRSAHelper([]byte(receiver.PublicKey), []byte(receiver.PrivateKey))
	if errs != nil && len(errs) > 0 {
		return "", errs[0]
	}
	decoded, err := helper.Decrypt(message)
	if err != nil {
		return "", fmt.Errorf(`message decryption failed. cause: %s`, err.Error())
	}

	return decoded, nil
}

// func (t *Message) AddHistories(histories []*Route) error {
// 	if t.Routes == nil {
// 		t.Routes = make([]*Route, 0)
// 	}

// 	for _, history := range histories {
// 		t.Routes = append(t.Routes, history)
// 	}
// 	return nil
// }

func (t *Topic) AddSender(orgID, publicKey string) {
	sender := NewSender(orgID, publicKey)
	t.Senders = append(t.Senders, sender)
}

func (t *Topic) AddReader(orgID, publicKey, privateKey string) {
	reader := NewReader(orgID, publicKey, privateKey)
	t.Readers = append(t.Readers, reader)
}

func (t *Message) AddRoute(networkID, orgID, userID, transactionID, updateTime, comment string) error {
	route := Route{}
	route.NetworkID = networkID
	route.OrgID = orgID
	route.UserID = userID
	route.TrxID = transactionID
	route.UpdateTime = updateTime
	route.Comment = comment
	if t.Routes == nil {
		t.Routes = make([]*Route, 0)
	}
	t.Routes = append(t.Routes, &route)
	return nil
}
