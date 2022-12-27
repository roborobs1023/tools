package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

const keyFile = "sec/aes.key"

func Encrypt(stringToEncrypt string) (string, error) {
	keyString, err := getKey(keyFile)
	if err != nil {
		return "", err
	}
	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	fmt.Sprintf("%x", ciphertext)
	return fmt.Sprintf("%x", ciphertext), nil
}

func Decrypt(encryptedString string) (string, error) {
	keyString, err := getKey(keyFile)
	if err != nil {
		return "", err
	}
	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return string(plaintext), nil
}

func getKey(keyFile string) (string, error) {
	bts, err := os.ReadFile(keyFile)
	if err != nil {
		if os.IsNotExist(err) {
			err = createKeyFile()
			if err != nil {
				return "", err
			}
		}
	}
	key := hex.EncodeToString(bts)
	//fmt.Printf("key to encrypt/decrypt: %s\n", key)

	return key, nil
}

func createKeyFile() error {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return err
	}
	err := os.MkdirAll("./sec", os.ModeDir)
	if err != nil {
		return err
	}
	err = os.WriteFile(keyFile, bytes, 0600)

	if err != nil {
		return err
	}
	return nil
}
