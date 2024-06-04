package tasksrv

import (
	"github.com/LeeZXin/zsf-utils/httputil"
)

var (
	httpClient = httputil.NewRetryableHttpClient()
)
