package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	//xx
	EMPTY = ""
	//AES-128: 16 * 8 = 128
	KEY_AES_128 = "1234567890abcdef"
	//AES-192: 24 * 8 = 192
	KEY_AES_192 = "1234567890abcd1234567890"
	//AES-256: 32 * 8 = 256
	KEY_AES_256 = "1234567890a1234567890b1234567890"
	//AES-128: 16 * 8 = 128
	ALL_SPACE = "                "
	//
	STR_NORMAL = "ABCD1234"
)

func Test_NewAESHelper(t *testing.T) {
	xaes, err := NewAESHelper(EMPTY)
	t.Log("check key is empty.")
	t.Logf(`XAESHelper: %v`, xaes)
	t.Log(err)
	assert.Nil(t, xaes, ``)
	assert.NotNil(t, err, "err should be nil.")

	xaes, err = NewAESHelper(KEY_AES_128)
	t.Log("check key length is 16.")
	t.Logf(`XAESHelper: %v`, xaes)
	t.Log(err)
	assert.NotNil(t, xaes, ``)
	assert.Nil(t, err, "err should be nil.")

	xaes, err = NewAESHelper(KEY_AES_192)
	t.Log("check key length is 24.")
	t.Logf(`XAESHelper: %v`, xaes)
	t.Log(err)
	assert.NotNil(t, xaes, ``)
	assert.Nil(t, err, "err should be nil.")

	xaes, err = NewAESHelper(KEY_AES_256)
	t.Log("check key length is 32.")
	t.Logf(`XAESHelper: %v`, xaes)
	t.Log(err)
	assert.NotNil(t, xaes, ``)
	assert.Nil(t, err, "err should be nil.")

	xaes, err = NewAESHelper(ALL_SPACE)
	t.Log("check key length is 16(all spaces).")
	t.Logf(`XAESHelper: %v`, xaes)
	t.Log(err)
	assert.NotNil(t, xaes, ``)
	assert.Nil(t, err, "err should be nil.")
}

func Test_AESEncryt(t *testing.T) {
	xaes, _ := NewAESHelper(KEY_AES_256)

	encoded, err := xaes.Encrypt(EMPTY)
	t.Log("check data is empty.")
	t.Logf(`encoded data is : %s`, encoded)
	t.Log(err)
	assert.NotNil(t, encoded, "encrypt result should not be nil.")
	assert.Nil(t, err, "err should be nil.")

	encoded, err = xaes.Encrypt(ALL_SPACE)
	t.Log("check data is all space.")
	t.Logf(`encoded data is : %s`, encoded)
	t.Log(err)
	assert.NotNil(t, encoded, ``)
	assert.Nil(t, err, "err should be nil.")

	encoded, err = xaes.Encrypt(STR_NORMAL)
	t.Log("check data is all space.")
	t.Logf(`encoded data is : %s`, encoded)
	t.Log(err)
	assert.NotNil(t, encoded, ``)
	assert.Nil(t, err, "err should be nil.")

	xaes = &AESHelper{}
	encoded, err = xaes.Encrypt(STR_NORMAL)
	t.Log("check aes helper generation without init.")
	t.Logf(`encoded data is : %s`, encoded)
	t.Log(err)
	assert.Empty(t, encoded, ``)
	assert.NotNil(t, err, "err should be nil.")

	xaes = &AESHelper{key: []byte("123")}
	encoded, err = xaes.Encrypt(STR_NORMAL)
	t.Log("check aes helper generation by wrong step.")
	t.Logf(`encoded data is : %s`, encoded)
	t.Log(err)
	assert.Empty(t, encoded, ``)
	assert.NotNil(t, err, "err should be nil.")
}

func Test_AESDecryt(t *testing.T) {
	xaes, _ := NewAESHelper(KEY_AES_256)
	encoded := "asde41"
	decoded, err := xaes.Decrypt(encoded)
	t.Log("check encoded data is not comfortable.")
	t.Logf(`decoded data is : %s`, decoded)
	assert.Empty(t, decoded, `decoded data should be nil`)
	t.Log(err)
	assert.NotNil(t, err, "err should not be nil.")

	encoded, _ = xaes.Encrypt(STR_NORMAL)
	//xaes, _ = NewAESHelper(KEY_AES_256)
	xaes = &AESHelper{key: []byte("123")}
	decoded, err = xaes.Decrypt(encoded)
	t.Log("check encoded data is not comfortable.")
	t.Logf(`decoded data is : %s`, decoded)
	assert.NotEqualf(t, decoded, STR_NORMAL, `decoded data should be '%s'`, STR_NORMAL)
	t.Log(err)
	assert.NotNil(t, err, "err should not be nil.")

	encoded, _ = xaes.Encrypt(STR_NORMAL)
	xaes, _ = NewAESHelper(KEY_AES_128)
	decoded, err = xaes.Decrypt(encoded)
	t.Log("check helper is not fake.")
	t.Logf(`decoded data is : %s`, decoded)
	assert.NotEqualf(t, decoded, STR_NORMAL, `decoded data should be '%s'`, STR_NORMAL)
	t.Log(err)
	assert.NotNil(t, err, "err should not be nil.")

	xaes, _ = NewAESHelper(KEY_AES_256)
	encoded, _ = xaes.Encrypt(STR_NORMAL)
	decoded, err = xaes.Decrypt(encoded)
	t.Log("check encoded data is not comfortable.")
	t.Logf(`encoded data is : %s`, encoded)
	assert.Equalf(t, decoded, STR_NORMAL, `decoded data should be '%s'`, STR_NORMAL)
	t.Log(err)
	assert.Equal(t, err, nil, "err should be nil.")
}
