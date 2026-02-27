package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

type PrivateKeyExt struct {
	Key       *secp256k1.PrivateKey
	ChainCode [32]byte
}

type PublicKeyExt struct {
	Key       *secp256k1.PublicKey
	ChainCode [32]byte
}

func Uint32ToBytes(i uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, i)
	return b
}

func PrivateToPublic(p PrivateKeyExt) PublicKeyExt {
	return PublicKeyExt{
		Key:       p.Key.PubKey(),
		ChainCode: p.ChainCode,
	}
}

func DeriveNormalPrivate(parent PrivateKeyExt, index uint32) PrivateKeyExt {

	if index >= 0x80000000 {
		panic("index is hardened")
	}

	parentPub := parent.Key.PubKey().SerializeCompressed()
	data := append(parentPub, Uint32ToBytes(index)...)

	mac := hmac.New(sha512.New, parent.ChainCode[:])
	mac.Write(data)
	I := mac.Sum(nil)

	var tweak secp256k1.ModNScalar
	tweak.SetByteSlice(I[:32])

	var parentScalar secp256k1.ModNScalar
	parentScalar.SetByteSlice(parent.Key.Serialize())

	tweak.Add(&parentScalar)
	childKey := secp256k1.NewPrivateKey(&tweak)

	var chain [32]byte
	copy(chain[:], I[32:])

	return PrivateKeyExt{
		Key:       childKey,
		ChainCode: chain,
	}
}

func DeriveHardened(parent PrivateKeyExt, index uint32) PrivateKeyExt {

	index += 0x80000000

	data := append([]byte{0x00}, parent.Key.Serialize()...)
	data = append(data, Uint32ToBytes(index)...)

	mac := hmac.New(sha512.New, parent.ChainCode[:])
	mac.Write(data)
	I := mac.Sum(nil)

	var tweak secp256k1.ModNScalar
	tweak.SetByteSlice(I[:32])

	var parentScalar secp256k1.ModNScalar
	parentScalar.SetByteSlice(parent.Key.Serialize())

	tweak.Add(&parentScalar)

	childKey := secp256k1.NewPrivateKey(&tweak)

	var chain [32]byte
	copy(chain[:], I[32:])

	return PrivateKeyExt{
		Key:       childKey,
		ChainCode: chain,
	}
}

func CreateMasterPrivate(seed []byte) PrivateKeyExt {

	mac := hmac.New(sha512.New, []byte("Bitcoin seed"))
	mac.Write(seed)
	I := mac.Sum(nil)

	privKey := secp256k1.PrivKeyFromBytes(I[:32])

	var chain [32]byte
	copy(chain[:], I[32:])

	return PrivateKeyExt{
		Key:       privKey,
		ChainCode: chain,
	}
}
