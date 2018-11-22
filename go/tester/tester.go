package main

import (
	"bytes"
	"fmt"

	"github.com/chaincode/common/crypto"
)

func main() {
	publicKey := *bytes.NewBufferString("")
	privateKey := *bytes.NewBufferString("")

	err := crypto.CreateKeyPair(&publicKey, &privateKey, 1024)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(publicKey.Bytes()))
	fmt.Println()
	fmt.Println(string(privateKey.Bytes()))

	// 	pubKey := `-----BEGIN PUBLIC KEY-----
	// MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAOjZcCaHIoj/N3ppJ0jg0yxepHhPbFey
	// a7NQursLMF95tTgCMoQlQJAEQfayc4PawnOPA34HBue/CdIu5BqkbxkCAwEAAQ==
	// -----END PUBLIC KEY-----`
	pubKey := "-----BEGIN PUBLIC KEY-----\n" +
		"MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAOjZcCaHIoj/N3ppJ0jg0yxepHhPbFey\n" +
		"a7NQursLMF95tTgCMoQlQJAEQfayc4PawnOPA34HBue/CdIu5BqkbxkCAwEAAQ==\n" +
		"-----END PUBLIC KEY-----\n"
	privKey := "-----BEGIN PRIVATE KEY-----\n" +
		"MIIBVAIBADANBgkqhkiG9w0BAQEFAASCAT4wggE6AgEAAkEA6NlwJociiP83emkn\n" +
		"SODTLF6keE9sV7Jrs1C6uwswX3m1OAIyhCVAkARB9rJzg9rCc48DfgcG578J0i7k\n" +
		"GqRvGQIDAQABAkEAuQ01aCk1dRL/kDVJl022BikhJMxaGkgd9+BMxqHZy8V1xaiw\n" +
		"RC4Z5u0oizAoO0ji88a47EbIu/bT6786wpDhkQIhAPzKISautItiha350Q9Q9B+2\n" +
		"9nStvDhhJbuxEuRLeu9vAiEA6856mvVRauUl83Alk7ibWC4+zdSjfZHd0qD3Uewp\n" +
		"xfcCIGsvGVdZhFwFbkESR76CyMAZx+45LDGLn4Ax2JzMFFgpAiAWH+BnC59g/TEL\n" +
		"XzlXW9nPcz9XRp00WexLJ+ksmZDtzwIgWwFu6JdVgQPnqjRymm3bPvgxJdNPa+D7\n" +
		"hJHgeD2IiE4=\n" +
		"-----END PRIVATE KEY-----\n"

	fmt.Println(pubKey)
	fmt.Println()
	fmt.Println(privKey)

	helper, _ := crypto.NewRSAHelper([]byte(pubKey), []byte(privKey))
	encoded, _ := helper.Encrypt("hello")
	fmt.Println(encoded)
	decoded, _ := helper.Decrypt(encoded)
	fmt.Println(decoded)

}
