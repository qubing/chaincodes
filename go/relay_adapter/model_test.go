package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SenderExist(t *testing.T) {
	topic := new(Topic)
	org1 := new(Sender)
	org1.OrgID = "org1"
	org1.PublicKey = "PK"
	topic.Senders = append(topic.Senders, org1)
	assert.True(t, topic.SenderExist("org1"))
}
