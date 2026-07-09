import { beforeEach, describe, expect, it } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

import { useApiTokens } from '@/composables/useApiTokens'

describe('useApiTokens', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('returns an em dash for missing or empty dates', () => {
    const { formatDate } = useApiTokens()

    expect(formatDate()).toBe('—')
    expect(formatDate(null)).toBe('—')
    expect(formatDate('')).toBe('—')
  })
})
