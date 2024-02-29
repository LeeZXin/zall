package i18n

import (
	"github.com/LeeZXin/zsf-utils/i18n"
	"github.com/LeeZXin/zsf/common"
	"github.com/LeeZXin/zsf/property/static"
	"path/filepath"
)

const (
	ZH_CN = "zh-CN"
	EN_US = "en-US"
)

func init() {
	lang := static.GetString("app.lang")
	if lang != "" {
		iniPath := filepath.Join(common.ResourcesDir, "i18n", lang+".ini")
		locale, err := i18n.NewImmutableLocaleFromIniFile(iniPath, lang)
		if err == nil {
			i18n.AddLocale(locale)
			i18n.SetDefaultLocale(lang)
			return
		}
	}
	iniPath := filepath.Join(common.ResourcesDir, "i18n", ZH_CN+".ini")
	locale, err := i18n.NewImmutableLocaleFromIniFile(iniPath, ZH_CN)
	if err == nil {
		i18n.AddLocale(locale)
		i18n.SetDefaultLocale(ZH_CN)
	}
}

func SupportedLangeList() []string {
	return []string{
		ZH_CN, EN_US,
	}
}

func GetByKey(item KeyItem) string {
	return i18n.GetOrDefault(item.Id, item.DefaultRet)
}
