package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ExitWithErrMsg(session ssh.Session, msg string) {
	fmt.Fprintf(session.Stderr(), "[%s] %s\n", time.Now().Format(time.DateTime), msg)
	session.Exit(1)
}

func GenRsaKeyPair(keyPath string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	f, err := os.OpenFile(keyPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = pem.Encode(f, privateKeyPEM); err != nil {
		return err
	}
	pub, err := gossh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	public := gossh.MarshalAuthorizedKey(pub)
	p, err := os.OpenFile(keyPath+".pub", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer p.Close()
	_, err = p.Write(public)
	return err
}

func CutEnv(envs []string) map[string]string {
	ret := make(map[string]string, len(envs))
	for _, env := range envs {
		k, v, b := strings.Cut(env, "=")
		if b {
			ret[k] = v
		}
	}
	return ret
}

func ReadOrGenRsaKey(hostKey string) (string, error) {
	if hostKey == "" {
		return "", errors.New("empty hostKey")
	}
	var err error
	if !filepath.IsAbs(hostKey) {
		hostKey, err = filepath.Abs(hostKey)
		if err != nil {
			return "", err
		}
	}
	if err = os.MkdirAll(filepath.Dir(hostKey), os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create dir %s: %v", filepath.Dir(hostKey), err)
	}
	exist, err := IsExist(hostKey)
	if err != nil {
		return "", fmt.Errorf("check host key failed %s: %v", hostKey, err)
	}
	if !exist {
		err = GenRsaKeyPair(hostKey)
		if err != nil {
			return "", fmt.Errorf("gen host key pair failed %s: %v", hostKey, err)
		}
	}
	return hostKey, nil
}
