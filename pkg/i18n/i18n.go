package i18n

import (
	"github.com/LeeZXin/zsf-utils/i18n"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/property/static"
	"path/filepath"
)

const (
	ZH_CN = "zh-cn"
	EN_US = "en-us"
)

func init() {
	iniPath := filepath.Join(common.ResourcesDir, "i18n", EN_US+".ini")
	locale, err := i18n.NewImmutableLocaleFromIniFile(iniPath, EN_US)
	if err == nil {
		i18n.AddLocale(locale)
	}
	iniPath = filepath.Join(common.ResourcesDir, "i18n", ZH_CN+".ini")
	locale, err = i18n.NewImmutableLocaleFromIniFile(iniPath, ZH_CN)
	if err == nil {
		i18n.AddLocale(locale)
	}
	i18n.SetDefaultLocale(static.GetString("app.lang"))
}

func GetByKey(item KeyItem) string {
	return i18n.GetOrDefault(item.Id, item.DefaultRet)
}

func GetByValue(val string) string {
	return i18n.Get(val)
}

func GetByLangAndValue(lang, val string) string {
	locale, b := i18n.GetLocale(lang)
	if !b {
		return ""
	}
	return locale.Get(val)
}
