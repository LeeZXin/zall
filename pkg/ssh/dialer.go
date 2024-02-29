package ssh

import (
	"errors"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"os"
	"path/filepath"
	"strings"
)

type DialerOpts struct {
	UserName string
	HostKey  string
}

type ProxyOpts struct {
	Envs map[string]string

	SrcFingerprint string
}

type Dialer struct {
	cfg *gossh.ClientConfig
}

func (d *Dialer) ProxySession(sshHost string, src ssh.Session, opts *ProxyOpts) error {
	// 建立SSH连接
	client, err := gossh.Dial("tcp", sshHost, d.cfg)
	if err != nil {
		return err
	}
	defer client.Close()
	target, err := client.NewSession()
	if err != nil {
		return err
	}
	for _, env := range src.Environ() {
		k, v, ok := strings.Cut(env, "=")
		if ok {
			if err = target.Setenv(k, v); err != nil {
				return err
			}
		}
	}
	if opts != nil {
		for k, v := range opts.Envs {
			if err = target.Setenv(k, v); err != nil {
				return err
			}
		}
	}
	err = target.Setenv("ZGIT_SRC_FINGERPRINT", opts.SrcFingerprint)
	target.Stdin = src
	target.Stdout = src
	target.Stderr = src.Stderr()
	target.Run(src.RawCommand())
	return nil
}

func NewDialer(opts *DialerOpts) (*Dialer, error) {
	if opts.UserName == "" {
		return nil, errors.New("empty username")
	}
	if opts.HostKey == "" || !filepath.IsAbs(opts.HostKey) {
		return nil, errors.New("wrong hostKey")
	}
	privateKey, err := os.ReadFile(opts.HostKey)
	if err != nil {
		return nil, err
	}
	keySigner, err := gossh.ParsePrivateKey(privateKey)
	cfg := &gossh.ClientConfig{
		Config: gossh.Config{
			KeyExchanges: keyExchanges,
			Ciphers:      ciphers,
			MACs:         macs,
		},
		User: opts.UserName,
		Auth: []gossh.AuthMethod{
			gossh.PublicKeys(keySigner),
		},
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
	}
	return &Dialer{
		cfg: cfg,
	}, nil
}
