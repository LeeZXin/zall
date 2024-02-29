package signature

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/keybase/go-crypto/openpgp"
	"github.com/keybase/go-crypto/openpgp/armor"
	"github.com/keybase/go-crypto/openpgp/packet"
	"io"
	"strings"
	"time"
)

const (
	StartGPGSigLineTag = "-----BEGIN PGP SIGNATURE-----"
	EndGPGSigLineTag   = "-----END PGP SIGNATURE-----"
)

type CommitSig string

func (s CommitSig) IsGPGSig() bool {
	return strings.HasPrefix(string(s), StartGPGSigLineTag)
}

func (s CommitSig) IsSSHSig() bool {
	return strings.HasPrefix(string(s), StartSSHSigLineTag)
}

func (s CommitSig) String() string {
	return string(s)
}

func ConvertArmoredGPGKeyString(content string) (openpgp.EntityList, error) {
	list, err := openpgp.ReadArmoredKeyRing(strings.NewReader(content))
	if err != nil {
		return nil, err
	}
	return CoalesceGpgEntityList(list), nil
}

func CheckArmoredDetachedSignature(ekeys openpgp.EntityList, token, signature string) (*openpgp.Entity, error) {
	signer, err := openpgp.CheckArmoredDetachedSignature(
		ekeys,
		strings.NewReader(token),
		strings.NewReader(signature),
	)
	if err != nil {
		signer, err = openpgp.CheckArmoredDetachedSignature(
			ekeys,
			strings.NewReader(token+"\n"),
			strings.NewReader(signature),
		)
		if err != nil {
			signer, err = openpgp.CheckArmoredDetachedSignature(
				ekeys,
				strings.NewReader(token+"\r\n"),
				strings.NewReader(signature),
			)
		}
	}
	return signer, err
}

func CoalesceGpgEntityList(ekeys openpgp.EntityList) openpgp.EntityList {
	id2key := map[string]*openpgp.Entity{}
	newEKeys := make([]*openpgp.Entity, 0, len(ekeys))
	for _, ekey := range ekeys {
		id := ekey.PrimaryKey.KeyIdString()
		if original, has := id2key[id]; has {
			// Coalesce this with the other one
			for _, subKey := range ekey.Subkeys {
				if subKey.PublicKey == nil {
					continue
				}
				found := false
				for _, originalSubKey := range original.Subkeys {
					if originalSubKey.PublicKey == nil {
						continue
					}
					if originalSubKey.PublicKey.KeyId == subKey.PublicKey.KeyId {
						found = true
						break
					}
				}
				if !found {
					original.Subkeys = append(original.Subkeys, subKey)
				}
			}
			for name, identity := range ekey.Identities {
				if _, has = original.Identities[name]; has {
					continue
				}
				original.Identities[name] = identity
			}
			continue
		}
		id2key[id] = ekey
		newEKeys = append(newEKeys, ekey)
	}
	return newEKeys
}

// GetGPGKeyExpiryTime extract the expiry time of primary key based on sig
func GetGPGKeyExpiryTime(e *openpgp.Entity) time.Time {
	expiry := time.Time{}
	// Extract self-sign for expire date based on : https://github.com/golang/crypto/blob/master/openpgp/keys.go#L165
	var selfSig *packet.Signature
	for _, ident := range e.Identities {
		if selfSig == nil {
			selfSig = ident.SelfSignature
		} else if ident.SelfSignature.IsPrimaryId != nil && *ident.SelfSignature.IsPrimaryId {
			selfSig = ident.SelfSignature
			break
		}
	}
	if selfSig.KeyLifetimeSecs != nil {
		expiry = e.PrimaryKey.CreationTime.Add(time.Duration(*selfSig.KeyLifetimeSecs) * time.Second)
	}
	return expiry
}

func ParseGPGPublicKey(content string) (*packet.PublicKey, error) {
	block, err := armor.Decode(strings.NewReader(content))
	if err != nil {
		return nil, err
	}
	if block.Type != "PGP PUBLIC KEY BLOCK" {
		return nil, fmt.Errorf("expected '" + openpgp.SignatureType + "', got: " + block.Type)
	}
	p, err := packet.Read(block.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key packet:%v", err)
	}
	sig, ok := p.(*packet.PublicKey)
	if !ok {
		return nil, fmt.Errorf("packet is not a public key")
	}
	return sig, nil
}

// Base64EncGPGPubKey encode public key content to base 64
func Base64EncGPGPubKey(pubkey *packet.PublicKey) (string, error) {
	var w bytes.Buffer
	err := pubkey.Serialize(&w)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(w.Bytes()), nil
}

func readerFromBase64(s string) (io.Reader, error) {
	bs, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(bs), nil
}

// Base64DecGPGPubKey decode public key content from base 64
func Base64DecGPGPubKey(content string) (*packet.PublicKey, error) {
	b, err := readerFromBase64(content)
	if err != nil {
		return nil, err
	}
	// Read key
	p, err := packet.Read(b)
	if err != nil {
		return nil, err
	}
	// Check type
	pkey, ok := p.(*packet.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not a public key")
	}
	return pkey, nil
}

// readArmoredSign
func readArmoredSign(r io.Reader) (body io.Reader, err error) {
	block, err := armor.Decode(r)
	if err != nil {
		return nil, err
	}
	if block.Type != openpgp.SignatureType {
		return nil, fmt.Errorf("expected '" + openpgp.SignatureType + "', got: " + block.Type)
	}
	return block.Body, nil
}

func ExtractGpgSignature(content string) (*packet.Signature, error) {
	r, err := readArmoredSign(strings.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("failed to read signature armor")
	}
	p, err := packet.Read(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read signature packet")
	}
	sig, ok := p.(*packet.Signature)
	if !ok {
		return nil, fmt.Errorf("packet is not a signature")
	}
	return sig, nil
}
