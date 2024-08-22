package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"math/big"
)

const (
	checksumLength = 4
	version        = byte(0x00)
)

type Wallet struct {
	PrivateKey []byte
	PublicKey  []byte
}

// Serialize Wallet to gob format
func (w *Wallet) GobEncode() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	
	// Encode the PrivateKey and PublicKey fields manually
	err := encoder.Encode(w.PrivateKey)
	if err != nil {
		return nil, err
	}
	
	err = encoder.Encode(w.PublicKey)
	if err != nil {
		return nil, err
	}
	
	return buf.Bytes(), nil
}

// Deserialize Wallet from gob format
func (w *Wallet) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	
	err := decoder.Decode(&w.PrivateKey)
	if err != nil {
		return err
	}
	
	err = decoder.Decode(&w.PublicKey)
	if err != nil {
		return err
	}
	
	return nil
}

func (w Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)

	versionedHash := append([]byte{version}, pubHash...)
	checksum := Checksum(versionedHash)

	fullHash := append(versionedHash, checksum...)
	address := Base58Encode(fullHash)
	return address
}

// Create new key pair and return the Wallet
func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pub
}

func MakeWallet() *Wallet {
	private, public := NewKeyPair()
	wallet := Wallet{
		PrivateKey: private.D.Bytes(),
		PublicKey:  public,
	}

	return &wallet
}

// Retrieve the ecdsa.PrivateKey from the serialized data
func (w *Wallet) GetPrivateKey() *ecdsa.PrivateKey {
	priv := new(ecdsa.PrivateKey)
	priv.D = new(big.Int).SetBytes(w.PrivateKey)
	priv.PublicKey.Curve = elliptic.P256()
	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(w.PrivateKey)
	return priv
}

func PublicKeyHash(pubKey []byte) []byte {
	// First SHA-256 hash
	firstHash := sha256.Sum256(pubKey)

	// Second SHA-256 hash
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:]
}

func Checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}
