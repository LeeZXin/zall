package util

import (
	"fmt"
	"github.com/LeeZXin/zsf-utils/idutil"
	"time"
)

func RandomIdWithTime() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixMilli(), idutil.RandomUuid()[:8])
}
