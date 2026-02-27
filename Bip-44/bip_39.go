package main

//let's suppose that our 12 words are sufficient entropy-wise

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"strings"
)

func DeriveSeedFromMnemonic(mnemonic []string, passphrase string) []byte {
	mnemonicSentence := strings.Join(mnemonic, " ")
	salt := "mnemonic" + passphrase

	return pbkdf2SHA512([]byte(mnemonicSentence), []byte(salt), 2048, 64)
}

func pbkdf2SHA512(password, salt []byte, iterations, keyLen int) []byte {
	hLen := 64
	numBlocks := (keyLen + hLen - 1) / hLen

	var derivedKey []byte

	for block := 1; block <= numBlocks; block++ {

		blockIndex := make([]byte, 4)
		binary.BigEndian.PutUint32(blockIndex, uint32(block))

		mac := hmac.New(sha512.New, password)
		mac.Write(salt)
		mac.Write(blockIndex)
		u := mac.Sum(nil)

		t := make([]byte, hLen)
		copy(t, u)

		for i := 1; i < iterations; i++ {
			mac = hmac.New(sha512.New, password)
			mac.Write(u)
			u = mac.Sum(nil)

			for j := 0; j < hLen; j++ {
				t[j] ^= u[j]
			}
		}

		derivedKey = append(derivedKey, t...)
	}

	return derivedKey[:keyLen]
}
