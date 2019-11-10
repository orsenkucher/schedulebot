package hash

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/gob"
)

// EncodeAsBase64 calculates sha1 hash and returns it as base64 string
func EncodeAsBase64(e interface{}) string {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(e); err != nil {
		panic(err)
	}
	bytes := b.Bytes()
	sha := sha1.New()
	sha.Write(bytes)
	bytes = sha.Sum(nil)
	b64 := base64.StdEncoding.EncodeToString(bytes)
	return b64
}
