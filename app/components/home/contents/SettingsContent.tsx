import { useMainStore } from '@/app/store/useMainStore'
import GeneralSettingsContent from './settings/GeneralSettingsContent'
import AudioSettingsContent from './settings/AudioSettingsContent'
import AccountSettingsContent from './settings/AccountSettingsContent'
import KeyboardSettingsContent from './settings/KeyboardSettingsContent'
import AdvancedSettingsContent from './settings/AdvancedSettingsContent'
import { useTranslation } from 'react-i18next'

export default function SettingsContent() {
  const { settingsPage, setSettingsPage } = useMainStore()
  const { t } = useTranslation()

  const settingsMenuItems = [
    { id: 'general', label: t('settings.general'), active: settingsPage === 'general' },
    { id: 'keyboard', label: t('settings.keyboard'), active: settingsPage === 'keyboard' },
    { id: 'audio', label: t('settings.audioAndMic'), active: settingsPage === 'audio' },
    { id: 'account', label: t('settings.account'), active: settingsPage === 'account' },
    { id: 'advanced', label: 'Advanced', active: settingsPage === 'advanced' },
  ]

  const renderSettingsContent = () => {
    switch (settingsPage) {
      case 'general':
        return <GeneralSettingsContent />
      case 'keyboard':
        return <KeyboardSettingsContent />
      case 'audio':
        return <AudioSettingsContent />
      case 'account':
        return <AccountSettingsContent />
      case 'advanced':
        return <AdvancedSettingsContent />
      default:
        return <GeneralSettingsContent />
    }
  }

  return (
    <div className="w-full px-32">
      <div className="space-y-6">
        {/* Horizontal Tab/Pill Selector */}
        <div className="flex gap-1 p-1 bg-slate-100 rounded-lg w-fit mx-auto">
          {settingsMenuItems.map(item => (
            <button
              key={item.id}
              onClick={() => setSettingsPage(item.id as any)}
              className={`px-4 py-2 rounded-md text-sm font-medium transition-all duration-200 ${
                item.active
                  ? 'bg-white text-slate-900 shadow-sm'
                  : 'text-slate-600 hover:text-slate-900 hover:bg-slate-50'
              }`}
            >
              {item.label}
            </button>
          ))}
        </div>

        {/* Content Area */}
        <div className="w-full pt-8">{renderSettingsContent()}</div>
      </div>
    </div>
  )
}
