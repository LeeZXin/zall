package ssh

import (
	"errors"
	"fmt"
	"github.com/LeeZXin/zall/util"
	"github.com/LeeZXin/zsf/logger"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"net"
	"strconv"
)

type Server struct {
	*ssh.Server
	Port int
}

type ServerOpts struct {
	Port             int
	HostKey          string
	PublicKeyHandler ssh.PublicKeyHandler
	SessionHandler   ssh.Handler
}

func NewServer(opts *ServerOpts) (*Server, error) {
	if opts == nil {
		return nil, errors.New("nil opts")
	}
	if opts.PublicKeyHandler == nil {
		return nil, errors.New("nil publicKeyHandler")
	}
	if opts.SessionHandler == nil {
		return nil, errors.New("nil sessionHandler")
	}
	if opts.Port <= 0 {
		return nil, errors.New("wrong port")
	}
	hostKey, err := util.ReadOrGenRsaKey(opts.HostKey)
	srv := &ssh.Server{
		Addr:             net.JoinHostPort("", strconv.Itoa(opts.Port)),
		PublicKeyHandler: opts.PublicKeyHandler,
		Handler:          opts.SessionHandler,
		ServerConfigCallback: func(ctx ssh.Context) *gossh.ServerConfig {
			config := &gossh.ServerConfig{
				Config: gossh.Config{
					KeyExchanges: keyExchanges,
					MACs:         macs,
					Ciphers:      ciphers,
				},
			}
			return config
		},
		PtyCallback: func(ctx ssh.Context, pty ssh.Pty) bool {
			return false
		},
	}
	if err = srv.SetOption(ssh.HostKeyFile(hostKey)); err != nil {
		return nil, fmt.Errorf("set host key failed: %v", err)
	}
	return &Server{
		Server: srv,
		Port:   opts.Port,
	}, nil
}

func (s *Server) OnApplicationStart() {
	go func() {
		logger.Logger.Infof("start ssh server port: %d", s.Port)
		err := s.ListenAndServe()
		if err != nil && err != ssh.ErrServerClosed {
			logger.Logger.Fatalf("start ssh server err: %v", err)
		}
	}()
}

func (s *Server) Order() int {
	return 0
}

func (s *Server) AfterInitialize() {}

func (s *Server) OnApplicationShutdown() {
	logger.Logger.Infof("shutdown ssh server: %d", s.Port)
	s.Close()
}
