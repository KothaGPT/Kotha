// Define the native binaries that are shared across platforms
// Only include binaries that actually exist under native/
const nativeBinaries = [
  'global-key-listener',
  'audio-recorder',
  'active-application',
]

const getMacResources = () =>
  nativeBinaries.map(binary => ({
    from: `native/${binary}/target/\${arch}-apple-darwin/release/${binary}`,
    to: `binaries/${binary}`,
  }))

const getWindowsResources = () =>
  nativeBinaries.map(binary => ({
    from: `native/${binary}/target/x86_64-pc-windows-gnu/release/${binary}.exe`,
    to: `binaries/${binary}.exe`,
  }))

module.exports = {
  appId: 'ai.kotha.kotha',
  productName: 'Kotha',
  copyright: 'Copyright © 2025 Demox Labs',
  directories: {
    buildResources: 'resources',
    output: 'dist',
  },
  files: [
    '!**/.vscode/*',
    '!src/*',
    '!electron.vite.config.{js,ts,mjs,cjs}',
    '!.eslintignore',
    '!.eslintrc.cjs',
    '!.prettierignore',
    '!.prettierrc.yaml',
    '!README.md',
    '!.env',
    '!.env.*',
    '!.npmrc',
    '!pnpm-lock.yaml',
    '!tsconfig.json',
    '!tsconfig.node.json',
    '!tsconfig.web.json',
    '!native/**',
    '!build-*.sh',
    {
      from: 'out',
      filter: ['**/*'],
    },
  ],
  asar: true,
  asarUnpack: ['resources/**'],
  extraMetadata: {
    version: process.env.VITE_KOTHA_VERSION || '0.0.0-dev',
  },
  protocols: {
    name: 'kotha',
    schemes: ['kotha'],
  },
  mac: {
    target: 'default',
    icon: 'resources/build/icon.icns',
    darkModeSupport: true,
    hardenedRuntime: true,
    gatekeeperAssess: false,
    identity: 'Demox Labs, Inc. (294ZSTM7UB)',
    notarize: true,
    entitlements: 'build/entitlements.mac.plist',
    entitlementsInherit: 'build/entitlements.mac.inherit.plist',
    extendInfo: {
      NSMicrophoneUsageDescription:
        'Kotha requires microphone access to transcribe your speech.',
    },
    extraResources: [
      ...getMacResources(),
      { from: 'resources/build/kotha-logo.png', to: 'build/kotha-logo.png' },
    ],
  },
  dmg: {
    artifactName: 'Kotha-Installer.${ext}',
  },
  win: {
    target: [
      {
        target: 'zip',
        arch: ['x64'],
      },
      {
        target: 'nsis',
        arch: ['x64'],
      },
    ],
    artifactName: '${productName}-${version}.${ext}',
    icon: 'resources/build/icon.ico',
    executableName: 'Kotha',
    requestedExecutionLevel: 'asInvoker',
    extraResources: [
      ...getWindowsResources(),
      { from: 'resources/build/kotha-logo.png', to: 'build/kotha-logo.png' },
    ],
    forceCodeSigning: false,
    asarUnpack: [
      'resources/**',
      '**/node_modules/@sentry/**',
      '**/node_modules/sqlite3/**',
    ],
  },
  nodeGypRebuild: false,
  buildDependenciesFromSource: false,
  nsis: {
    shortcutName: '${productName}',
    uninstallDisplayName: '${productName}-uninstaller',
    createDesktopShortcut: false,
    createStartMenuShortcut: true,
    oneClick: false,
    perMachine: false,
    allowToChangeInstallationDirectory: false,
    deleteAppDataOnUninstall: true,
  },
}
