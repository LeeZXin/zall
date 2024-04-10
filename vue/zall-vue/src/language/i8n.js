import { createI18n } from 'vue-i18n'
import Zh from './zh'
import En from './en'

const i18n = createI18n({
    legacy: false,
    locale: 'zh',
    messages: {
        en: { ...En },
        zh: { ...Zh }
    },
    fallbackLocale: 'zh',
    globalInjection: true
})

export default i18n;

