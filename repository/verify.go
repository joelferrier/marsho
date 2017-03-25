package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type gpgKeyring struct {
	version        string
	defaultKeyring string
	keys           *openpgp.EntityList
}

type gpgKey struct {
	key *openpgp.Entity
}

func getDefaultKeyring() (*gpgKeyring, error) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(usr.HomeDir, ".gnupg", "pubring.gpg")

	var keyring openpgp.EntityList
	// check if the default gpg keyring exists
	if _, err := os.Stat(path); err == nil {
		keyringFile, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer keyringFile.Close()

		keyring, err = openpgp.ReadKeyRing(keyringFile)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return &gpgKeyring{
		"gpg",
		path,
		&keyring,
	}, nil
}

//TODO: handle armored and unarmored keys
func readKey(reader io.Reader) (*gpgKey, error) {
	block, err := armor.Decode(reader)
	if err != nil {
		return nil, err
	}
	preader := packet.NewReader(block.Body)
	keyEntity, err := openpgp.ReadEntity(preader)
	if err != nil {
		return nil, err
	}
	return &gpgKey{
		keyEntity,
	}, err
}

func (k *gpgKey) fingerprint() string {
	return strings.ToUpper(hex.EncodeToString(k.key.PrimaryKey.Fingerprint[:20]))
}

func (kr *gpgKeyring) contains(k *gpgKey) bool {
	for _, entity := range *kr.keys {
		if entity.PrimaryKey.Fingerprint == k.key.PrimaryKey.Fingerprint {
			log.Debug(fmt.Sprintf(
				"found key with fingerprint %s in keychain",
				hex.EncodeToString(k.key.PrimaryKey.Fingerprint[:20]),
			))
			return true
		}
	}

	return false
}

func (kr *gpgKeyring) verifyDetachedSig(data io.Reader, sig io.Reader) (*gpgKey, error) {
	signer, err := openpgp.CheckDetachedSignature(*kr.keys, data, sig)
	return &gpgKey{
		signer,
	}, err
}

func sha256sum(data []byte, checksum string) (bool, string) {
	calcSum := sha256.Sum256(data)
	calcSumString := hex.EncodeToString(calcSum[:32])
	if calcSumString == checksum {
		return true, calcSumString
	} else {
		return false, calcSumString
	}
}
