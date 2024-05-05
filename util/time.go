package util

import (
	"fmt"
	"github.com/LeeZXin/zall/pkg/i18n"
	"github.com/LeeZXin/zsf-utils/idutil"
	"strconv"
	"time"
)

var (
	day   = 24 * time.Hour
	month = 30 * day
	year  = 12 * month
)

func ReadableTimeComparingNow(t time.Time) string {
	d := time.Since(t)
	if d < time.Second {
		return "0" + i18n.GetByKey(i18n.TimeBeforeSecondUnit)
	}
	if d < time.Minute {
		return strconv.FormatInt(int64(d.Seconds()), 10) + i18n.GetByKey(i18n.TimeBeforeSecondUnit)
	}
	if d < time.Hour {
		return strconv.FormatInt(int64(d.Minutes()), 10) + i18n.GetByKey(i18n.TimeBeforeMinuteUnit)
	}
	if d < day {
		return strconv.FormatInt(int64(d.Hours()), 10) + i18n.GetByKey(i18n.TimeBeforeHourUnit)
	}
	if d < month {
		return strconv.FormatInt(int64(d/day), 10) + i18n.GetByKey(i18n.TimeBeforeDayUnit)
	}
	if d < year {
		return strconv.FormatInt(int64(d/month), 10) + i18n.GetByKey(i18n.TimeBeforeMonthUnit)
	}
	return strconv.FormatInt(int64(d/year), 10) + i18n.GetByKey(i18n.TimeBeforeYearUnit)
}

func RandomIdWithTime() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixMilli(), idutil.RandomUuid()[:8])
}
