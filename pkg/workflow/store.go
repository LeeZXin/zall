package workflow

import (
	"bufio"
	"fmt"
	"github.com/LeeZXin/zall/util"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	RunningStatus = "running"
	SuccessStatus = "success"
	FailStatus    = "fail"
	TimeoutStatus = "timeout"
	CancelStatus  = "cancel"
	QueueStatus   = "queue"
	UnknownStatus = "unknown"
)

const (
	originFileName = "origin"
	statusFileName = "status"
	beginFileName  = "begin"
	errLogFileName = "error.log"
	logFileName    = "log"
)

func toStatusMsg(status string, duration time.Duration) string {
	return fmt.Sprintf("%s %d", status, duration.Milliseconds())
}

func toStatusMsgBytes(status string, duration time.Duration) []byte {
	return []byte(toStatusMsg(status, duration))
}

type Store interface {
	StoreStatus(string, time.Duration) error
	StoreBeginTime(time.Time) error
	StoreErrLog(error) error
	StoreOrigin([]byte) error
	StoreLog(io.Reader) error
}

type fileStore struct {
	BaseDir string
}

func newFileStore(dir string) Store {
	return &fileStore{
		BaseDir: dir,
	}
}

func (s *fileStore) StoreStatus(status string, duration time.Duration) error {
	return util.WriteFile(filepath.Join(s.BaseDir, statusFileName),
		toStatusMsgBytes(status, duration))
}

func (s *fileStore) StoreBeginTime(beginTime time.Time) error {
	return util.WriteFile(filepath.Join(s.BaseDir, beginFileName),
		[]byte(strconv.FormatInt(beginTime.UnixMilli(), 10)))
}
func (s *fileStore) StoreErrLog(err error) error {
	return util.WriteFile(filepath.Join(s.BaseDir, errLogFileName),
		[]byte(err.Error()))
}
func (s *fileStore) StoreOrigin(input []byte) error {
	return util.WriteFile(filepath.Join(s.BaseDir, originFileName),
		input)
}
func (s *fileStore) StoreLog(reader io.Reader) error {
	var logFile *os.File
	// 记录日志
	logFile, err := os.OpenFile(filepath.Join(s.BaseDir, logFileName), os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err == nil {
		defer logFile.Close()
		// 增加缓存
		writer := bufio.NewWriter(logFile)
		defer writer.Flush()
		_, err = io.Copy(writer, reader)
	}
	return err
}
