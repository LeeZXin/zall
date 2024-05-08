package git

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/LeeZXin/zall/pkg/git/process"
	"github.com/LeeZXin/zsf/logger"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	gitExecutablePath string
)

func initCommand() {
	var err error
	gitExecutablePath, err = exec.LookPath("git")
	if err != nil {
		logger.Logger.Fatalf("could not LookPath err: %v", err)
	}
}

var (
	globalCmdArgs      = make([]string, 0)
	passThroughEnvKeys = []string{
		"GNUPGHOME",
	}
)

// addGlobalCmdArgs not thread safe
func addGlobalCmdArgs(args ...string) {
	globalCmdArgs = append(globalCmdArgs, args...)
}

type Command struct {
	args        []string
	invalidArgs []string
}

func (c *Command) AddArgs(args ...string) *Command {
	c.args = append(c.args, args...)
	return c
}

func (c *Command) AddDynamicArgs(args ...string) *Command {
	if c.checkSafeDynamicArg(args...) {
		c.args = append(c.args, args...)
	}
	return c
}

func (c *Command) isSafeDynamicArg(arg string) bool {
	return len(strings.Fields(arg)) == 1 && (arg == "" || arg[0] != '-')
}

func (c *Command) checkSafeDynamicArg(args ...string) bool {
	for _, arg := range args {
		if !c.isSafeDynamicArg(arg) {
			c.invalidArgs = append(c.invalidArgs, arg)
			return false
		}
	}
	return true
}

func (c *Command) Run(ctx context.Context, ros ...RunOpts) (*Result, error) {
	stdOut := new(bytes.Buffer)
	if err := c.run(ctx, append(ros, withStdOut(stdOut))...); err != nil {
		return nil, err
	}
	return &Result{
		stdOut: stdOut,
	}, nil
}

func (c *Command) RunWithStdout(ctx context.Context, output io.Writer, ros ...RunOpts) error {
	return c.run(ctx, append(ros, withStdOut(output))...)
}

func (c *Command) RunWithReadPipe(ctx context.Context, ros ...RunOpts) *ReadPipeResult {
	reader, writer := io.Pipe()
	go func() {
		if err := c.run(ctx, append(ros, withStdOut(writer))...); err != nil {
			writer.CloseWithError(err)
		} else {
			writer.Close()
		}
	}()
	return &ReadPipeResult{
		reader: reader,
	}
}

func (c *Command) RunWithStdinPipe(ctx context.Context, ros ...RunOpts) *ReadWritePipeResult {
	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()
	go func() {
		if err := c.run(ctx, append(ros, WithStdin(stdinReader), withStdOut(stdoutWriter))...); err != nil {
			stdoutWriter.CloseWithError(err)
			stdinReader.CloseWithError(err)
		} else {
			stdoutWriter.Close()
			stdinReader.Close()
		}
	}()
	return &ReadWritePipeResult{
		reader: stdoutReader,
		writer: stdinWriter,
	}
}

func (c *Command) run(ctx context.Context, ros ...RunOpts) error {
	if len(c.invalidArgs) > 0 {
		return fmt.Errorf("invalid arguments: %v", c.invalidArgs)
	}
	opts := new(runOpts)
	for _, o := range ros {
		o(opts)
	}
	if ctx == nil {
		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithTimeout(context.Background(), 360*time.Second)
		defer cancelFunc()
	}
	cmd := exec.CommandContext(ctx, gitExecutablePath, c.args...)
	if opts.Env == nil {
		cmd.Env = os.Environ()
	} else {
		cmd.Env = append(os.Environ(), opts.Env...)
	}
	stdErr := new(bytes.Buffer)
	process.SetSysProcAttribute(cmd)
	cmd.Env = append(cmd.Env, CommonGitCmdEnvs()...)
	cmd.Dir = opts.Dir
	cmd.Stdout = opts.StdOut
	cmd.Stderr = stdErr
	cmd.Stdin = opts.Stdin
	if err := cmd.Start(); err != nil {
		return stdErrorResult(err, bytesToString(stdErr.Bytes()))
	}
	err := cmd.Wait()
	if err != nil && ctx.Err() != context.DeadlineExceeded {
		return stdErrorResult(err, bytesToString(stdErr.Bytes()))
	}
	if ctx.Err() != nil {
		return stdErrorResult(ctx.Err(), bytesToString(stdErr.Bytes()))
	}
	return nil
}

func NewCommand(args ...string) *Command {
	if args == nil {
		args = []string{}
	}
	return &Command{
		args:        append(globalCmdArgs, args...),
		invalidArgs: make([]string, 0),
	}
}

func CommonEnvs() []string {
	envs := []string{
		"HOME=" + HomeDir(),
		"GIT_NO_REPLACE_OBJECTS=1",
	}
	for _, key := range passThroughEnvKeys {
		if val, ok := os.LookupEnv(key); ok {
			envs = append(envs, key+"="+val)
		}
	}
	return envs
}

func CommonGitCmdEnvs() []string {
	return append(CommonEnvs(), "LC_ALL=C", "GIT_TERMINAL_PROMPT=0")
}

func IsExitCode(err error, code int) bool {
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		return exitError.ExitCode() == code
	}
	return false
}

func ExecutablePath() string {
	return gitExecutablePath
}
