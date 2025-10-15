import { Switch } from '@/app/components/ui/switch'
import { useSettingsStore } from '@/app/store/useSettingsStore'
import { useWindowContext } from '@/app/components/window/WindowContext'

export default function GeneralSettingsContent() {
  const {
    shareAnalytics,
    launchAtLogin,
    showKothaBarAlways,
    showAppInDock,
    language,
    setShareAnalytics,
    setLaunchAtLogin,
    setShowKothaBarAlways,
    setShowAppInDock,
    setLanguage,
  } = useSettingsStore()

  const windowContext = useWindowContext()

  return (
    <div className="space-y-8">
      <div>
        <div className="space-y-6">
          <div className="flex items-center justify-between">
            <div>
              <div className="text-sm font-medium">Language / ভাষা</div>
              <div className="text-xs text-gray-600 mt-1">
                Select your preferred language / আপনার পছন্দের ভাষা নির্বাচন করুন
              </div>
            </div>
            <select
              value={language}
              onChange={(e) => setLanguage(e.target.value as 'en' | 'bn')}
              className="px-3 py-2 border border-slate-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="en">English</option>
              <option value="bn">বাংলা</option>
            </select>
          </div>

          <div className="flex items-center justify-between">
            <div>
              <div className="text-sm font-medium">Share analytics</div>
              <div className="text-xs text-gray-600 mt-1">
                Share anonymous usage data to help us improve Kotha.
              </div>
            </div>
            <Switch
              checked={shareAnalytics}
              onCheckedChange={setShareAnalytics}
            />
          </div>

          <div className="flex items-center justify-between">
            <div>
              <div className="text-sm font-medium">Launch at Login</div>
              <div className="text-xs text-gray-600 mt-1">
                Open Kotha automatically when your computer starts.
              </div>
            </div>
            <Switch
              checked={launchAtLogin}
              onCheckedChange={setLaunchAtLogin}
            />
          </div>

          <div className="flex items-center justify-between">
            <div>
              <div className="text-sm font-medium">
                Show Kotha bar at all times
              </div>
              <div className="text-xs text-gray-600 mt-1">
                Show the Kotha bar at all times.
              </div>
            </div>
            <Switch
              checked={showKothaBarAlways}
              onCheckedChange={setShowKothaBarAlways}
            />
          </div>

          {windowContext?.window?.platform === 'darwin' && (
            <div className="flex items-center justify-between">
              <div>
                <div className="text-sm font-medium">Show app in dock</div>
                <div className="text-xs text-gray-600 mt-1">
                  Show the Kotha app in the dock for quick access.
                </div>
              </div>
              <Switch
                checked={showAppInDock}
                onCheckedChange={setShowAppInDock}
              />
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
