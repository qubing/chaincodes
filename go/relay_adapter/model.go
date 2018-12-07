package main

import (
	"encoding/json"
	"fmt"

	"github.com/chaincodes/common/crypto"
)

type DocumentType string

const (
	DOC_TOPIC_IN    DocumentType = "TOPIC_IN"
	DOC_SESSION_IN  DocumentType = "SESSION_IN"
	DOC_TOPIC_OUT   DocumentType = "TOPIC_OUT"
	DOC_SESSION_OUT DocumentType = "SESSION_OUT"
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
		fmt.Println()
		return "{}", err
	}
	return string(bytes), nil
}

// ParseJSON - parse json text to object
func ParseJSON(dataObject JSONModel, dataJSON string) error {
	err := json.Unmarshal([]byte(dataJSON), &dataObject)
	if err != nil {
		fmt.Printf("[ParseJSON]%s", err.Error())
		fmt.Println()
		return err
	}
	return nil
}

// ParseJSON - parse json text to object
func ParseJSONs(dataObjects []JSONModel, dataJSON string) error {
	err := json.Unmarshal([]byte(dataJSON), &dataObjects)
	if err != nil {
		fmt.Printf("[ParseJSONs]%s", err.Error())
		fmt.Println()
		return err
	}
	return nil
}

type AbstractDoc struct {
	DocType DocumentType `json:"doc_type"`
	Id      string       `json:"id"`
}

type Topic struct {
	AbstractDoc
	Name    string    `json:"name"`
	Senders []*Sender `json:"senders"`
	Readers []*Reader `json:"readers"`
}

type Orgnization struct {
	OrgID string `json:"org_id"`
}

type Sender struct {
	Orgnization
	PublicKey string `json:"public_key"`
}

type Reader struct {
	Orgnization
	PublicKey   string `json:"public_key"`
	PrivateHash string `json:"private_hash"`
}

type Session struct {
	AbstractDoc
	TopicName string   `json:"topic_name"`
	Message   string   `json:"message"`
	Histories []*Route `json:"histories"`
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

func (t *Session) ParseJSON(dataJSON string) error {
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

func (t *Topic) NewSession() *Session {
	session := new(Session)
	if t.DocType == DOC_TOPIC_IN {
		session.DocType = DOC_SESSION_IN
	} else {
		session.DocType = DOC_SESSION_OUT
	}

	return session
}

func NewSender(orgID string, publicKey string) *Sender {
	org := Sender{PublicKey: publicKey}
	org.OrgID = orgID
	return &org
}

func NewReader(orgID string, publicKey string, privateHash string) *Reader {
	org := Reader{PublicKey: publicKey, PrivateHash: privateHash}
	org.OrgID = orgID
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

func (t *Reader) DecryptMessage(orgID string, message string, privateKey []byte) (string, error) {
	helper, errs := crypto.NewRSAHelper([]byte(t.PublicKey), privateKey)
	if errs != nil && len(errs) > 0 {
		return "", errs[0]
	}
	decoded, err := helper.Decrypt(message)
	if err != nil {
		return "", fmt.Errorf(`message decryption failed. cause: %s`, err.Error())
	}

	return decoded, nil
}

func (t *Topic) AddSender(orgID, publicKey string) {
	sender := NewSender(orgID, publicKey)
	t.Senders = append(t.Senders, sender)
}

func (t *Topic) AddReader(orgID, publicKey, privateKey string) {
	reader := NewReader(orgID, publicKey, privateKey)
	t.Readers = append(t.Readers, reader)
}
