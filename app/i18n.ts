import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'
import LanguageDetector from 'i18next-browser-languagedetector'

import enTranslations from './locales/en.json'
import bnTranslations from './locales/bn.json'

const resources = {
  en: {
    translation: enTranslations,
  },
  bn: {
    translation: bnTranslations,
  },
}

// Initialize language from electron store if available
const getInitialLanguage = (): 'en' | 'bn' => {
  try {
    if (typeof window !== 'undefined' && window.electron?.store) {
      const storedSettings = window.electron.store.get('settings')
      return (storedSettings?.language as 'en' | 'bn') || 'bn'
    }
  } catch (error) {
    console.warn('Failed to get language from electron store:', error)
  }
  return 'bn' // Default to Bangla
}

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    lng: getInitialLanguage(), // Set initial language from settings
    fallbackLng: 'bn', // Set Bangla as the default fallback language
    debug: process.env.NODE_ENV === 'development',

    interpolation: {
      escapeValue: false, // React already escapes values
    },

    detection: {
      order: ['localStorage', 'navigator', 'htmlTag'],
      lookupLocalStorage: 'kotha-language',
      caches: ['localStorage'],
    },
  })

// Make i18n available globally
if (typeof window !== 'undefined') {
  window.i18n = i18n
}

export default i18n
