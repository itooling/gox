package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"os"
)

//GenRsaKey
//openssl genrsa -out private.pem 2048
//openssl rsa -in private.pem -pubout -out public.pem
func GenRsaKey(bit int) error {
	//private
	privateKey, err := rsa.GenerateKey(rand.Reader, bit)
	if err != nil {
		return err
	}

	derStream, _ := x509.MarshalPKCS8PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}

	file, err := os.Create("private.pem")
	if err != nil {
		return err
	}

	err = pem.Encode(file, block)
	if err != nil {
		return err
	}

	//public
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	block = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPkix,
	}

	file, err = os.Create("public.pem")
	if err != nil {
		return err
	}

	err = pem.Encode(file, block)
	if err != nil {
		return err
	}

	return nil
}

func Encrypt(src, public []byte) ([]byte, error) {
	if block, _ := pem.Decode(public); block != nil {
		public = block.Bytes
	}

	key, err := x509.ParsePKIXPublicKey(public)
	if err != nil {
		return nil, err
	}

	pub := key.(*rsa.PublicKey)
	res, err := rsa.EncryptPKCS1v15(rand.Reader, pub, src)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func Decrypt(src, private []byte) ([]byte, error) {
	if block, _ := pem.Decode(private); block != nil {
		private = block.Bytes
	}

	key, err := x509.ParsePKCS8PrivateKey(private)
	if err != nil {
		return nil, err
	}

	pri := key.(*rsa.PrivateKey)
	res, err := rsa.DecryptPKCS1v15(rand.Reader, pri, src)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func EncryptOAEP(src, public []byte) ([]byte, error) {
	if block, _ := pem.Decode(public); block != nil {
		public = block.Bytes
	}

	key, err := x509.ParsePKIXPublicKey(public)
	if err != nil {
		return nil, err
	}

	pub := key.(*rsa.PublicKey)
	res, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, pub, src, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func DecryptOAEP(src, private []byte) ([]byte, error) {
	if block, _ := pem.Decode(private); block != nil {
		private = block.Bytes
	}

	key, err := x509.ParsePKCS8PrivateKey(private)
	if err != nil {
		return nil, err
	}

	pri := key.(*rsa.PrivateKey)
	res, err := rsa.DecryptOAEP(sha1.New(), rand.Reader, pri, src, nil)
	if err != nil {
		return nil, err
	}

	return res, nil
}
