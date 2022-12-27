package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func GenerateRSAKeys() error {
	bitSize := 4096

	// Generate RSA key.
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return err
	}

	err = key.Validate()
	if err != nil {
		return err
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
		return err
	}
	//fmt.Println("Private Key:", string(keyPEM))
	// Write public key to file.
	if err := os.WriteFile("id_rsa.pub", pubPEM, 0755); err != nil {
		return err
	}

	return nil
}

// CompareKeys takes in a filename of private key and compares it with the
func CompareKeys(keyid string) (bool, error) {
	private, err := os.ReadFile(keyid)
	if err != nil {
		return false, err
	}

	//fmt.Println("Private Key:", rsaPrivateKey)
	pub := keyid + ".pub"
	public, err := os.ReadFile(pub)
	if err != nil {
		return false, err
	}

	block, _ := pem.Decode(private)

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return false, err
	}
	pubBlock, _ := pem.Decode(public)

	pubKey, err := x509.ParsePKCS1PublicKey(pubBlock.Bytes)
	if err != nil {
		return false, err
	}

	return key.PublicKey.Equal(pubKey), nil
}
