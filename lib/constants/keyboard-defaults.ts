import { KothaMode } from '@/app/generated/kotha_pb'

// Platform-specific keyboard shortcut defaults
export const KOTHA_MODE_SHORTCUT_DEFAULTS_MAC = {
  [KothaMode.TRANSCRIBE]: ['fn'],
  [KothaMode.EDIT]: ['control-left', 'fn'],
}

export const KOTHA_MODE_SHORTCUT_DEFAULTS_WIN = {
  [KothaMode.TRANSCRIBE]: ['option-left', 'control-left'],
  [KothaMode.EDIT]: ['option-left', 'shift-left'],
}

// Helper to detect platform - works in both main and renderer process
export function getPlatform(): 'darwin' | 'win32' {
  if (typeof process !== 'undefined' && process.platform) {
    return process.platform as 'darwin' | 'win32'
  }
  // Fallback if process is not available
  return 'darwin'
}

// Get platform-specific defaults
export function getKothaModeShortcutDefaults(
  platform?: 'darwin' | 'win32',
): Record<KothaMode, string[]> {
  const currentPlatform = platform || getPlatform()

  if (currentPlatform === 'darwin') {
    return KOTHA_MODE_SHORTCUT_DEFAULTS_MAC
  } else {
    return KOTHA_MODE_SHORTCUT_DEFAULTS_WIN
  }
}

// For backward compatibility, export the defaults for the current platform
export const KOTHA_MODE_SHORTCUT_DEFAULTS = getKothaModeShortcutDefaults()
