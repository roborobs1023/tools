package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func GenerateRSAKeys() {
	bitSize := 4096

	// Generate RSA key.
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		panic(err)
	}

	err = key.Validate()
	if err != nil {
		panic(err)
	}
	// Extract public component.
	pub := key.Public()

	// Encode private key to PKCS#1 ASN.1 PEM.
	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	// Encode public key to PKCS#1 ASN.1 PEM.
	pubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pub.(*rsa.PublicKey)),
		},
	)

	// Write private key to file.
	if err := os.WriteFile("id_rsa", keyPEM, 0700); err != nil {
		panic(err)
	}
	//fmt.Println("Private Key:", string(keyPEM))
	// Write public key to file.
	if err := os.WriteFile("id_rsa.pub", pubPEM, 0755); err != nil {
		panic(err)
	}

	//fmt.Println("Public Key:", string(pubPEM))
}


//CompareKeys takes in a filename of private key and compares it with the 
func CompareKeys(keyid string) bool {
	private, err := os.ReadFile(keyid)
	if err != nil {
		log.Fatalln(err)
	}

	//fmt.Println("Private Key:", rsaPrivateKey)
	pub := keyid + ".pub"
	public, err := os.ReadFile(pub)
	if err != nil {
		log.Fatalln(err)
	}
	//rsaPublicKey := string(public)
	//fmt.Println("Public Key:", rsaPublicKey)

	block, _ := pem.Decode(private)
	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	pubBlock, _ := pem.Decode(public)
	pubKey, _ := x509.ParsePKCS1PublicKey(pubBlock.Bytes)
	return key.PublicKey.Equal(pubKey)
}
