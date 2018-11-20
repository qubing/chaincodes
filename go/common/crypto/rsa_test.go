package crypto

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateKeyPair(t *testing.T) {
	publicKey := *bytes.NewBufferString("")
	privateKey := *bytes.NewBufferString("")

	err := CreateKeyPair(nil, &privateKey, 12)
	assert.NotNil(t, err)

	err = CreateKeyPair(&publicKey, nil, 12)
	assert.NotNil(t, err)

	err = CreateKeyPair(nil, nil, 12)
	assert.NotNil(t, err)

	err = CreateKeyPair(&publicKey, &privateKey, 0)
	assert.NotNil(t, err)
	err = CreateKeyPair(&publicKey, &privateKey, 11)
	assert.NotNil(t, err)
	err = CreateKeyPair(&publicKey, &privateKey, 12)
	assert.Nil(t, err)
	err = CreateKeyPair(&publicKey, &privateKey, 2048)
	assert.Nil(t, err)
}

func Test_NewRSAHelper(t *testing.T) {
	publicKey := *bytes.NewBufferString("")
	privateKey := *bytes.NewBufferString("")

	xrsa, errs := NewRSAHelper([]byte("a"), []byte("bbb"))
	t.Log("check nil pub key and priv key.")
	t.Logf(`XRSA: %v`, xrsa)
	t.Log(errs)
	assert.NotNil(t, xrsa, ``)
	assert.NotEqual(t, len(errs), 0, "xxxx")

	CreateKeyPair(&publicKey, &privateKey, 2048)

	xrsa, errs = NewRSAHelper(nil, nil)
	t.Log("check both pub key and priv key are nil at the same time.")
	t.Logf(`XRSA: %v`, xrsa)
	t.Log(errs)
	assert.NotNil(t, xrsa, ``)
	assert.NotEqual(t, len(errs), 0, "xxxx")

	xrsa, errs = NewRSAHelper(publicKey.Bytes(), nil)
	t.Log("check pub key is not nil and priv key is nil.")
	t.Log(xrsa)
	t.Log(errs)
	assert.NotNil(t, xrsa)
	assert.Equal(t, len(errs), 0)

	xrsa, errs = NewRSAHelper(nil, privateKey.Bytes())
	t.Log("check pub key is nil and priv key is not nil.")
	t.Log(xrsa)
	t.Log(errs)
	assert.NotNil(t, xrsa)
	assert.Equal(t, len(errs), 0)

	xrsa, errs = NewRSAHelper(publicKey.Bytes(), privateKey.Bytes())
	t.Log("check both pub key and priv key are not nil.")
	t.Log(xrsa)
	t.Log(errs)
	assert.NotNil(t, xrsa)
	assert.Equal(t, len(errs), 0)
}

func Test_Encryt(t *testing.T) {
	publicKey := *bytes.NewBufferString("")
	privateKey := *bytes.NewBufferString("")

	err := CreateKeyPair(&publicKey, &privateKey, 2048)
	if err != nil {
		return
	}

	data := "hello"
	dataEmpty := ""

	t.Log(">>check both pub key and priv key nil at the same time.")
	xrsa, _ := NewRSAHelper(nil, nil)
	t.Logf(">>>data is '%s'", data)
	encoded, err := xrsa.Encrypt(data)
	assert.NotNil(t, err)
	assert.Equal(t, encoded, "")

	t.Logf(">>>data is '%s'", dataEmpty)
	encoded, err = xrsa.Encrypt(dataEmpty)
	assert.NotNil(t, err)
	assert.Equal(t, encoded, "")

	t.Log(">>check pub key is not nil and priv key is nil.")
	xrsa, _ = NewRSAHelper(publicKey.Bytes(), nil)
	t.Logf(">>>data is '%s'", data)
	encoded, err = xrsa.Encrypt(data)
	assert.Nil(t, err)
	assert.NotEqual(t, encoded, "")

	t.Logf(">>>data is '%s'", dataEmpty)
	encoded, err = xrsa.Encrypt(dataEmpty)
	assert.Nil(t, err)
	assert.Equal(t, encoded, "")

	t.Log(">>check pub key is nil and priv key is not nil.")
	xrsa, _ = NewRSAHelper(nil, privateKey.Bytes())
	t.Logf(">>>data is '%s'", data)
	encoded, err = xrsa.Encrypt(data)
	assert.NotNil(t, err)
	assert.Equal(t, encoded, "")

	t.Logf(">>>data is '%s'", dataEmpty)
	encoded, err = xrsa.Encrypt(dataEmpty)
	assert.NotNil(t, err)
	assert.Equal(t, encoded, "")

	t.Log(">>check both pub key and priv key are not nil.")
	xrsa, _ = NewRSAHelper(publicKey.Bytes(), privateKey.Bytes())
	t.Logf(">data is '%s'", data)
	encoded, err = xrsa.Encrypt(data)
	assert.Nil(t, err)
	assert.NotEqual(t, encoded, "")
	t.Logf(">>>data is '%s'", dataEmpty)
	encoded, err = xrsa.Encrypt(dataEmpty)
	assert.Nil(t, err)
	assert.Equal(t, encoded, "")
}

func Test_Decryt(t *testing.T) {
	publicKey := *bytes.NewBufferString("")
	privateKey := *bytes.NewBufferString("")

	err := CreateKeyPair(&publicKey, &privateKey, 2048)
	if err != nil {
		return
	}

	data := "hello"
	xrsa, _ := NewRSAHelper(publicKey.Bytes(), nil)
	t.Logf(">>>data is '%s'", data)
	encoded, err := xrsa.Encrypt(data)

	t.Log(">>check both pub key and priv key nil at the same time.")
	xrsa, _ = NewRSAHelper(nil, nil)
	decoded, err := xrsa.Decrypt(encoded)
	assert.NotNil(t, err)
	assert.Equal(t, decoded, "")

	t.Log(">>check pub key is not nil and priv key is nil.")
	xrsa, _ = NewRSAHelper(publicKey.Bytes(), nil)
	t.Logf(">>>data is '%s'", data)
	decoded, err = xrsa.Decrypt(encoded)
	assert.NotNil(t, err)
	assert.Equal(t, decoded, "")

	t.Log(">>check pub key is nil and priv key is not nil.")
	xrsa, _ = NewRSAHelper(nil, privateKey.Bytes())
	t.Logf(">>>data is '%s'", data)
	decoded, err = xrsa.Decrypt(encoded)
	assert.NotNil(t, err)
	assert.Equal(t, decoded, "")

	t.Log(">>check both pub key and priv key are not nil.")
	xrsa, _ = NewRSAHelper(publicKey.Bytes(), privateKey.Bytes())
	t.Logf(">>>data is '%s'", data)
	decoded, err = xrsa.Decrypt(encoded)
	assert.Nil(t, err)
	assert.NotEqual(t, decoded, "")
}
func Test_Sign(t *testing.T) {
	publicKey := *bytes.NewBufferString("")
	privateKey := *bytes.NewBufferString("")

	err := CreateKeyPair(&publicKey, &privateKey, 2048)
	if err != nil {
		return
	}

	data := "hello"
	dataEmpty := ""

	t.Log(">>check both pub key and priv key nil at the same time.")
	xrsa, _ := NewRSAHelper(nil, nil)
	t.Logf(">>>data is '%s'", data)
	sign, err := xrsa.Sign(data)
	assert.NotNil(t, err)
	assert.Equal(t, sign, "")
	t.Logf(">>>data is '%s'", dataEmpty)
	sign, err = xrsa.Sign(dataEmpty)
	assert.NotNil(t, err)
	assert.Equal(t, sign, "")

	t.Log(">>check pub key is not nil and priv key is nil.")
	xrsa, _ = NewRSAHelper(publicKey.Bytes(), nil)
	t.Logf(">>>data is '%s'", data)
	sign, err = xrsa.Sign(data)
	assert.NotNil(t, err)
	assert.Equal(t, sign, "")
	t.Logf(">>>data is '%s'", dataEmpty)
	sign, err = xrsa.Sign(dataEmpty)
	assert.NotNil(t, err)
	assert.Equal(t, sign, "")

	t.Log(">>check pub key is nil and priv key is not nil.")
	xrsa, _ = NewRSAHelper(nil, privateKey.Bytes())
	t.Logf(">>>data is '%s'", data)
	sign, err = xrsa.Sign(data)
	assert.Nil(t, err)
	assert.NotEqual(t, sign, "")
	t.Logf(">>>data is '%s'", dataEmpty)
	sign, err = xrsa.Sign(dataEmpty)
	assert.Nil(t, err)
	assert.NotEqual(t, sign, "")

	t.Log(">>check both pub key and priv key are not nil.")
	xrsa, _ = NewRSAHelper(publicKey.Bytes(), privateKey.Bytes())
	t.Logf(">data is '%s'", data)
	sign, err = xrsa.Sign(data)
	assert.Nil(t, err)
	assert.NotEqual(t, sign, "")
	t.Logf(">>>data is '%s'", dataEmpty)
	sign, err = xrsa.Sign(dataEmpty)
	assert.Nil(t, err)
	assert.NotEqual(t, sign, "")
}

func Test_Verify(t *testing.T) {
	publicKey := *bytes.NewBufferString("")
	privateKey := *bytes.NewBufferString("")

	err := CreateKeyPair(&publicKey, &privateKey, 2048)
	if err != nil {
		return
	}

	data := "hello"
	dataEmpty := ""
	xrsa, _ := NewRSAHelper(publicKey.Bytes(), privateKey.Bytes())
	t.Logf(">>>data is '%s'", data)
	sign, _ := xrsa.Sign(data)
	t.Log(sign)
	signEmpty, _ := xrsa.Sign(dataEmpty)
	t.Log(signEmpty)

	t.Log(">>check both pub key and priv key nil at the same time.")
	xrsa, _ = NewRSAHelper(nil, nil)
	t.Logf(">>>data is '%s'", data)
	err = xrsa.Verify(data, sign)
	assert.NotNil(t, err)
	t.Logf(">>>data is '%s'", dataEmpty)
	err = xrsa.Verify(dataEmpty, signEmpty)
	assert.NotNil(t, err)

	t.Log(">>check pub key is not nil and priv key is nil.")
	xrsa, _ = NewRSAHelper(publicKey.Bytes(), nil)
	t.Logf(">>>data is '%s'", data)
	err = xrsa.Verify(data, sign)
	assert.Nil(t, err)
	t.Logf(">>>data is '%s'", dataEmpty)
	err = xrsa.Verify(dataEmpty, signEmpty)
	assert.Nil(t, err)

	t.Log(">>check pub key is nil and priv key is not nil.")
	xrsa, _ = NewRSAHelper(nil, privateKey.Bytes())
	t.Logf(">>>data is '%s'", data)
	err = xrsa.Verify(data, sign)
	assert.NotNil(t, err)
	t.Logf(">>>data is '%s'", dataEmpty)
	err = xrsa.Verify(dataEmpty, signEmpty)
	assert.NotNil(t, err)

	t.Log(">>check both pub key and priv key are not nil.")
	xrsa, _ = NewRSAHelper(publicKey.Bytes(), privateKey.Bytes())
	t.Logf(">>>data is '%s'", data)
	err = xrsa.Verify(data, sign)
	assert.Nil(t, err)
	t.Logf(">>>data is '%s'", dataEmpty)
	err = xrsa.Verify(dataEmpty, signEmpty)
	assert.Nil(t, err)

	t.Log(">>check wrong pub key.")
	xrsa, _ = NewRSAHelper([]byte("aaaa"), nil)
	t.Logf(">>>data is '%s'", data)
	err = xrsa.Verify(data, sign)
	assert.NotNil(t, err)
	t.Logf(">>>data is '%s'", dataEmpty)
	err = xrsa.Verify(dataEmpty, signEmpty)
	assert.NotNil(t, err)

	t.Log(">>check wrong priv key.")
	xrsa, _ = NewRSAHelper([]byte("aaaa"), []byte("bbbb"))
	t.Logf(">>>data is '%s'", data)
	err = xrsa.Verify(data, sign)
	assert.NotNil(t, err)
	t.Logf(">>>data is '%s'", dataEmpty)
	err = xrsa.Verify(dataEmpty, signEmpty)
	assert.NotNil(t, err)
}
