package common

import (
	"fmt"

	"github.com/iost-official/Go-IOS-Protocol/log"
	"github.com/golang/protobuf/proto"
)

//go:generate gencode go -schema=structs.schema -package=common

type SignAlgorithm uint8

const (
	Secp256k1 SignAlgorithm = iota
)

type Signature struct {
	Algorithm SignAlgorithm

	Sig    []byte
	Pubkey []byte
}

func Sign(algo SignAlgorithm, info, privkey []byte) (Signature, error) {
	s := Signature{}
	s.Algorithm = algo
	switch algo {
	case Secp256k1:
		s.Pubkey = CalcPubkeyInSecp256k1(privkey)
		s.Sig = SignInSecp256k1(info, privkey)
		return s, nil
	}
	return s, fmt.Errorf("algorithm not exist")
}

func VerifySignature(info []byte, s Signature) bool {
	switch s.Algorithm {
	case Secp256k1:
		return VerifySignInSecp256k1(info, s.Pubkey, s.Sig)
	}
	return false
}

func (s *Signature) Encode() []byte {
	sr := &SignatureRaw{
		Algorithm:int32(s.Algorithm),
		Sig:s.Sig,
		PubKey:s.Pubkey,
	}
	b, err := proto.Marshal(sr)
	if err != nil {
		log.Log.E("Error in Encode of signature ", s.Pubkey, err.Error())
		return nil
	}
	return b
}

func (s *Signature) Decode(b []byte) error {
	var sr *SignatureRaw
	err := proto.Unmarshal(b, sr)
	if err != nil {
		log.Log.E("Error in Decode of signature ", b, err.Error())
		return err
	}
	s.Algorithm = SignAlgorithm(sr.Algorithm)
	s.Sig = sr.Sig
	s.Pubkey = sr.PubKey
	return err
}

func (s *Signature) Hash() []byte {
	return Sha256(s.Encode())
}
