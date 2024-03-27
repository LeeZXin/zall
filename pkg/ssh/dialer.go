package ssh

import (
	"errors"
	"github.com/LeeZXin/zall/util"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"os"
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
	return target.Run(src.RawCommand())
}

func NewDialer(opts *DialerOpts) (*Dialer, error) {
	if opts.UserName == "" {
		return nil, errors.New("empty username")
	}
	hostKey, err := util.ReadOrGenRsaKey(opts.HostKey)
	if err != nil {
		return nil, err
	}
	privateKey, err := os.ReadFile(hostKey)
	if err != nil {
		return nil, err
	}
	keySigner, err := gossh.ParsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	cfg := NewCommonClientConfig(opts.UserName, keySigner)
	return &Dialer{
		cfg: cfg,
	}, nil
}

func NewCommonClientConfig(username string, signer gossh.Signer) *gossh.ClientConfig {
	return &gossh.ClientConfig{
		Config: gossh.Config{
			KeyExchanges: keyExchanges,
			Ciphers:      ciphers,
			MACs:         macs,
		},
		User: username,
		Auth: []gossh.AuthMethod{
			gossh.PublicKeys(signer),
		},
		HostKeyCallback: gossh.InsecureIgnoreHostKey(),
	}
}
