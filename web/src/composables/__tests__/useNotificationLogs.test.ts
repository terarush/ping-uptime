import { beforeEach, describe, expect, it } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

import { useNotificationLogs } from '@/composables/useNotificationLogs'

describe('useNotificationLogs', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('maps known channel types to icons', () => {
    const { channelIcon } = useNotificationLogs()

    expect(channelIcon('email')).toBe('📧')
    expect(channelIcon('discord')).toBe('💬')
    expect(channelIcon('discord_bot')).toBe('🤖')
    expect(channelIcon('telegram')).toBe('📱')
    expect(channelIcon('slack')).toBe('💼')
    expect(channelIcon('webhook')).toBe('🔗')
  })

  it('returns the fallback icon for unknown channel types', () => {
    const { channelIcon } = useNotificationLogs()

    expect(channelIcon('sms')).toBe('📨')
  })

  it('returns an em dash for missing dates', () => {
    const { formatDate } = useNotificationLogs()

    expect(formatDate()).toBe('—')
  })
})
