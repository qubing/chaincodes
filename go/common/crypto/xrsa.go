package crypto

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
)

const (
	PUB_KEY                = "PUBLIC KEY"
	PRIV_KEY               = "PRIVATE KEY"
	CHAR_SET               = "UTF-8"
	BASE_64_FORMAT         = "UrlSafeNoPadding"
	RSA_ALGORITHM_KEY_TYPE = "PKCS8"
	RSA_ALGORITHM_SIGN     = crypto.SHA256
)

// RSAHelper - RSAHelper structure definition.
type RSAHelper struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

// CreateKeyPair - generate key-pair.
func CreateKeyPair(pubKeyWriter, prvKeyWriter io.Writer, keyLength int) error {
	if pubKeyWriter == nil || prvKeyWriter == nil {
		return errors.New("public key writer or private key writer is nil")
	}

	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, keyLength)
	if err != nil {
		return err
	}
	derStream := marshalPKCS8PrivateKey(privateKey)
	block := &pem.Block{
		Type:  PRIV_KEY,
		Bytes: derStream,
	}
	err = pem.Encode(prvKeyWriter, block)
	if err != nil {
		return err
	}

	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  PUB_KEY,
		Bytes: derPkix,
	}
	err = pem.Encode(pubKeyWriter, block)

	return err
}

// NewRSAHelper - generate an instance of RSAHelper strcture.
func NewRSAHelper(publicKey []byte, privateKey []byte) (*RSAHelper, []error) {
	xrsa := &RSAHelper{}
	errs := []error{}
	// parse public key
	if len(publicKey) > 0 {
		block, _ := pem.Decode(publicKey)
		if block != nil {
			pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err == nil {
				pub := pubInterface.(*rsa.PublicKey)
				xrsa.publicKey = pub
			} else {
				errs = append(errs, err)
			}
		} else {
			errs = append(errs, errors.New("public key error"))
		}
	}

	//parse private key
	if len(privateKey) > 0 {
		block, _ := pem.Decode(privateKey)
		if block != nil {
			privInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err == nil {
				pri, _ := privInterface.(*rsa.PrivateKey)
				xrsa.privateKey = pri
			} else {
				errs = append(errs, errors.New("private key error"))
			}
		} else {
			errs = append(errs, errors.New("private key error"))
		}
	}

	if xrsa.publicKey == nil && xrsa.privateKey == nil {
		errs = append(errs, errors.New("both public key and private key not provided"))
	}

	if len(errs) > 0 {
		return xrsa, errs
	}

	return xrsa, nil
}

// Encrypt - encrypt data with public key.
func (xrsa *RSAHelper) Encrypt(data string) (string, error) {
	if xrsa == nil || xrsa.publicKey == nil {
		return "", errors.New("can not encrypt, because public key not provided")
	}

	partLen := xrsa.publicKey.N.BitLen()/8 - 11
	chunks := split([]byte(data), partLen)

	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		bytes, err := rsa.EncryptPKCS1v15(rand.Reader, xrsa.publicKey, chunk)
		if err != nil {
			return "", err
		}
		buffer.Write(bytes)
	}

	return base64.RawURLEncoding.EncodeToString(buffer.Bytes()), nil
}

// Decrypt - decrypt data with public key and private key.
func (xrsa *RSAHelper) Decrypt(encrypted string) (string, error) {
	if xrsa == nil || xrsa.publicKey == nil || xrsa.privateKey == nil {
		return "", errors.New("can not decrypt, because public key or private key not provided")
	}

	partLen := xrsa.publicKey.N.BitLen() / 8
	raw, err := base64.RawURLEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	chunks := split(raw, partLen)

	buffer := bytes.NewBufferString("")
	for _, chunk := range chunks {
		decrypted, err1 := rsa.DecryptPKCS1v15(rand.Reader, xrsa.privateKey, chunk)
		if err1 != nil {
			return "", err1
		}
		buffer.Write(decrypted)
	}

	return buffer.String(), err
}

// Sign - sign data with private key.
func (xrsa *RSAHelper) Sign(data string) (string, error) {
	if xrsa == nil || xrsa.privateKey == nil {
		return "", errors.New("can not sign, because private key not provided")
	}

	h := RSA_ALGORITHM_SIGN.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)

	sign, err := rsa.SignPKCS1v15(rand.Reader, xrsa.privateKey, RSA_ALGORITHM_SIGN, hashed)
	if err == nil {
		return base64.RawURLEncoding.EncodeToString(sign), nil
	}

	return "", err
}

//Verify - verify data with sign.
func (xrsa *RSAHelper) Verify(data string, sign string) error {
	if xrsa == nil || xrsa.publicKey == nil {
		return errors.New("can not verify, because public key not provided")
	}

	h := RSA_ALGORITHM_SIGN.New()
	h.Write([]byte(data))
	hashed := h.Sum(nil)

	decodedSign, err := base64.RawURLEncoding.DecodeString(sign)
	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(xrsa.publicKey, RSA_ALGORITHM_SIGN, hashed, decodedSign)
}

func marshalPKCS8PrivateKey(key *rsa.PrivateKey) []byte {
	info := struct {
		Version             int
		PrivateKeyAlgorithm []asn1.ObjectIdentifier
		PrivateKey          []byte
	}{}
	info.Version = 0
	info.PrivateKeyAlgorithm = make([]asn1.ObjectIdentifier, 1)
	info.PrivateKeyAlgorithm[0] = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	info.PrivateKey = x509.MarshalPKCS1PrivateKey(key)

	k, _ := asn1.Marshal(info)
	return k
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
