import { beforeEach, describe, expect, it } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

import { useSslMonitors } from '@/composables/useSslMonitors'

describe('useSslMonitors', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('returns the expected style classes for days remaining', () => {
    const { daysClass } = useSslMonitors()

    expect(daysClass(-1)).toBe('text-red-500 font-bold')
    expect(daysClass(0)).toBe('text-red-500 font-bold')
    expect(daysClass(13)).toBe('text-red-400 font-semibold')
    expect(daysClass(14)).toBe('text-amber-400 font-semibold')
    expect(daysClass(30)).toBe('text-amber-400 font-semibold')
    expect(daysClass(31)).toBe('text-emerald-500')
  })

  it('maps SSL statuses to the correct badge variants', () => {
    const { statusVariant } = useSslMonitors()

    expect(statusVariant('valid')).toBe('bg-emerald-500/10 text-emerald-500 border-emerald-500/20')
    expect(statusVariant('expiring_soon')).toBe(
      'bg-amber-500/10 text-amber-500 border-amber-500/20',
    )
    expect(statusVariant('expired')).toBe('bg-red-500/10 text-red-500 border-red-500/20')
    expect(statusVariant('error')).toBe('bg-destructive/10 text-destructive border-destructive/20')
    expect(statusVariant('unknown')).toBe('bg-muted text-muted-foreground border-border')
  })

  it('returns an em dash for missing dates', () => {
    const { formatDate } = useSslMonitors()

    expect(formatDate()).toBe('—')
  })
})
