package signalstor // 	"github.com/davidwalter0/go-signalstor"

import (
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"log"
	"math"

	"github.com/btcsuite/btcutil/base58"
)

func Base32Encode(in []byte) (out string) {
	out = base32.StdEncoding.EncodeToString(in)
	return
}
func Base64Encode(in []byte) (out string) {
	out = base64.StdEncoding.EncodeToString(in)
	return
}

func Base32Decode(in string) (out []byte) {
	var err error
	out, err = base32.StdEncoding.DecodeString(in)
	CheckError(err)
	return
}

func Base64Decode(in string) (out []byte) {
	var err error
	out, err = base64.StdEncoding.DecodeString(in)
	CheckError(err)
	return
}

// Base58Encode
func Base58Encode(data []byte) string {
	return base58.Encode(data)
}

// Base58Decode
func Base58Decode(addr string) []byte {
	return base58.Decode(addr)
}

// CheckError standardize error handling
func CheckError(err error) {
	if err != nil {
		log.Println(err)
	}
}

// Float64frombytes import of a value
func Float64frombytes(bytes []byte) float64 {
	bits := binary.BigEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

// Float64bytes export to []byte
func Float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, bits)
	return bytes
}

// Int64FromBytes import of a value
func Int64FromBytes(bytes []byte) int64 {
	i := binary.BigEndian.Uint64(bytes)
	return int64(i)
}

// Int64ToBytes export to []byte
func Int64ToBytes(i int64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(i))
	return bytes
}
