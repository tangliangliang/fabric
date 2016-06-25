/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	mathRand "math/rand"
	"strings"
	"time"

	gp "google/protobuf"

	"golang.org/x/crypto/sha3"
)

// ComputeCryptoHash should be used in openchain code so that we can change the actual algo used for crypto-hash at one place
func ComputeCryptoHash(data []byte) (hash []byte) {
	hash = make([]byte, 64)
	sha3.ShakeSum256(hash, data)
	return
}

// GenerateBytesUUID returns a UUID based on RFC 4122 returning the generated bytes
func GenerateBytesUUID() []byte {
	uuid := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, uuid)
	if err != nil {
		panic(fmt.Sprintf("Error generating UUID: %s", err))
	}

	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80

	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return uuid
}

// GenerateIntUUID returns a UUID based on RFC 4122 returning a big.Int
func GenerateIntUUID() *big.Int {
	uuid := GenerateBytesUUID()
	z := big.NewInt(0)
	return z.SetBytes(uuid)
}

// GenerateUUID returns a UUID based on RFC 4122
func GenerateUUID() string {
	uuid := GenerateBytesUUID()
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

// CreateUtcTimestamp returns a google/protobuf/Timestamp in UTC
func CreateUtcTimestamp() *gp.Timestamp {
	now := time.Now().UTC()
	secs := now.Unix()
	nanos := int32(now.UnixNano() - (secs * 1000000000))
	return &(gp.Timestamp{Seconds: secs, Nanos: nanos})
}

//GenerateHashFromSignature returns a hash of the combined parameters
func GenerateHashFromSignature(path string, ctor string, args []string) []byte {
	fargs := ctor
	if args != nil {
		for _, str := range args {
			fargs = fargs + str
		}
	}
	cbytes := []byte(path + fargs)

	b := make([]byte, len(cbytes))
	copy(b, cbytes)
	hash := ComputeCryptoHash(b)
	return hash
}

// FindMissingElements identifies the elements of the first slice that are not present in the second
// The second slice is expected to be a subset of the first slice
func FindMissingElements(all []string, some []string) (delta []string) {
all:
	for _, v1 := range all {
		for _, v2 := range some {
			if strings.Compare(v1, v2) == 0 {
				continue all
			}
		}
		delta = append(delta, v1)
	}
	return
}

// Shuffle shuffles a given string slice
func Shuffle(a []string) []string {
	b := make([]string, len(a))
	copy(b, a)
	for i := range b {
		j := mathRand.Intn(i + 1)
		b[i], b[j] = b[j], b[i]
	}
	return b
}
