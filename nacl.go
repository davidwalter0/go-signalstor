package signalstor // 	"github.com/davidwalter0/go-signalstor"
import (
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

var (
	// ErrEncrypt package error message
	ErrEncrypt = errors.New("secret: NaclEncryption failed")
	// ErrDecrypt package error message
	ErrDecrypt = errors.New("secret: NaclDecryption failed")
)

const (
	// KeySize recommendation from reference material
	KeySize = 32
	// NonceSize recommendation from reference material
	NonceSize = 24
)

// Base58Key
type Base58Key string

// Base58Nonce
type Base58Nonce string

// Key for NACL encryption
type Key [KeySize]byte

// Nonce for NACL encryption
type Nonce [NonceSize]byte

// Encode base58
func (key *Key) Encode() Base58Key {
	return Base58Key(Base58Encode((*key)[:]))
}

// DecodeKey base58
func DecodeKey(k string) (out Key) {
	var key = Base58Decode(k)
	out = Key{}
	copy(out[:], key)
	return
}

// DecodeKey base58
func (k Base58Key) Decode() (out Key) {
	var key = Base58Decode(string(k))
	out = Key{}
	copy(out[:], key)
	return
}

// RawPtr (*[KeySize]byte) from key type
func (key *Key) RawPtr() *[KeySize]byte {
	return (*[KeySize]byte)(key)
}

// RawPtr (*[NonceSize]byte) from nonce type
func (nonce *Nonce) RawPtr() *[NonceSize]byte {
	return (*[NonceSize]byte)(nonce)
}

// Encode base58
func (nonce *Nonce) Encode() Base58Nonce {
	return Base58Nonce(Base58Encode((*nonce)[:]))
}

// DecodeNonce base58
func DecodeNonce(n string) (out Nonce) {
	var nonce = Base58Decode(n)
	out = Nonce{}
	copy(out[:], nonce)
	return
}

// DecodeNonce base58
func (n Base58Nonce) Decode() (out Nonce) {
	var nonce = Base58Decode(string(n))
	out = Nonce{}
	copy(out[:], nonce)
	return
}

// NaclTool to provide arguments and inputs and outputs
type NaclTool struct {
	Key           *Key
	Nonce         *Nonce
	ID            *[]byte
	PlainMessage  *[]byte
	CypherMessage *[]byte
}

// NewEncryptionNaclTool generates a key
func NewEncryptionNaclTool(plain []byte) *NaclTool {
	var err error
	var nt = &NaclTool{
		Key:   &Key{},
		Nonce: &Nonce{},
	}
	if nt.Key, err = GenerateKey(); err != nil {
		panic(err)
	}
	if nt.Nonce, err = GenerateNonce(); err != nil {
		panic(err)
	}
	nt.PlainMessage = &plain

	return nt
}

// NewDecryptionNaclTool requires key
func NewDecryptionNaclTool(key *Key, encrypted []byte) *NaclTool {
	var nt = &NaclTool{
		Key:   &Key{},
		Nonce: &Nonce{},
	}
	copy(nt.Key[:], key[:])

	nt.CypherMessage = &encrypted
	return nt
}

// func (nt NaclTool) Encrypt() {
// 	var err error
// 	nt.Buffer, err = NaclEncrypt(nt.Key[:KeySize], nt.Message[:])
// }

// GenerateKey creates a new random secret key.
func GenerateKey() (*Key, error) {
	key := new([KeySize]byte)
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		return nil, err
	}

	return (*Key)(key), nil
}

// GenerateNonce creates a new random nonce.
func GenerateNonce() (*Nonce, error) {
	nonce := new([NonceSize]byte)
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return nil, err
	}

	return (*Nonce)(nonce), nil
}

// NaclEncrypt generates a random nonce and NaclEncrypts the input using
// NaCl's secretbox package. The nonce is prepended to the ciphertext.
// A sealed message will the same size as the original message plus
// secretbox.Overhead bytes long.
func (nt *NaclTool) NaclEncrypt() error {
	var out = make([]byte, len(nt.Nonce))
	copy(out, nt.Nonce[:])
	out = secretbox.Seal(out, *nt.PlainMessage, nt.Nonce.RawPtr(), nt.Key.RawPtr())
	nt.CypherMessage = &out
	return nil
}

// NaclEncrypt generates a random nonce and encrypts the input using
// NaCl's secretbox package. The nonce is prepended to the ciphertext.
// A sealed message will the same size as the original message plus
// secretbox.Overhead bytes long.
func NaclEncrypt(key *[KeySize]byte, message []byte) (out []byte, err error) {
	var nonce *Nonce
	if nonce, err = GenerateNonce(); err != nil {
		return nil, ErrEncrypt
	}

	out = make([]byte, len(nonce))
	copy(out, nonce[:])
	out = secretbox.Seal(out, message, nonce.RawPtr(), key)
	return out, nil
}

// NaclDecrypt extracts the nonce from the ciphertext, and attempts to
// NaclDecrypt with NaCl's secretbox.
func (nt *NaclTool) NaclDecrypt() error {

	if nt.CypherMessage == nil || len(*nt.CypherMessage) < (NonceSize+secretbox.Overhead) {
		return ErrDecrypt
	}

	var ok bool
	copy(nt.Nonce[:], (*nt.CypherMessage)[:NonceSize])
	var out []byte
	if out, ok = secretbox.Open(nil, (*nt.CypherMessage)[NonceSize:], nt.Nonce.RawPtr(), nt.Key.RawPtr()); !ok {
		return ErrDecrypt
	}
	nt.PlainMessage = &out
	return nil
}

// NaclDecrypt extracts the nonce from the ciphertext, and attempts to
// NaclDecrypt with NaCl's secretbox.
func NaclDecrypt(key *[KeySize]byte, message []byte) ([]byte, error) {
	if len(message) < (NonceSize + secretbox.Overhead) {
		return nil, ErrDecrypt
	}

	var nonce [NonceSize]byte
	copy(nonce[:], message[:NonceSize])
	out, ok := secretbox.Open(nil, message[NonceSize:], &nonce, key)
	if !ok {
		return nil, ErrDecrypt
	}
	return out, nil
}
