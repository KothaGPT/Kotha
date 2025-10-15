import { DEFAULT_ADVANCED_SETTINGS } from '../../constants/generated-defaults.js'
import { KothaMode } from '../../generated/kotha_pb.js'

export const KOTHA_MODE_PROMPT: { [key in KothaMode]: string } = {
  [KothaMode.TRANSCRIBE]: DEFAULT_ADVANCED_SETTINGS.transcriptionPrompt,
  [KothaMode.EDIT]: DEFAULT_ADVANCED_SETTINGS.editingPrompt,
}

export const KOTHA_MODE_SYSTEM_PROMPT: { [key in KothaMode]: string } = {
  [KothaMode.TRANSCRIBE]: 'You are a helpful AI transcription assistant.',
  [KothaMode.EDIT]: 'You are an AI assistant helping to edit documents.',
}
