import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

import { useIncidents } from '@/composables/useIncidents'
import type { Incident } from '@/stores/incidents'

describe('useIncidents', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.useRealTimers()
  })

  it('formats resolved incident durations', () => {
    const { getDuration } = useIncidents()

    expect(
      getDuration(createIncident('2026-01-01T00:00:00.000Z', '2026-01-01T00:00:45.000Z')),
    ).toBe('45s')
    expect(
      getDuration(createIncident('2026-01-01T00:00:00.000Z', '2026-01-01T00:03:05.000Z')),
    ).toBe('3m 5s')
    expect(
      getDuration(createIncident('2026-01-01T00:00:00.000Z', '2026-01-01T02:05:10.000Z')),
    ).toBe('2h 5m')
  })

  it('uses the current time for active incident durations', () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date('2026-01-01T00:01:10.000Z'))

    const { getDuration } = useIncidents()

    expect(getDuration(createIncident('2026-01-01T00:00:00.000Z'))).toBe('1m 10s')
  })

  it('returns an em dash for missing dates', () => {
    const { formatDate } = useIncidents()

    expect(formatDate()).toBe('—')
  })
})

function createIncident(createdAt: string, resolvedAt?: string): Incident {
  return {
    id: 1,
    monitor_id: 1,
    user_id: 1,
    status: resolvedAt ? 'resolved' : 'active',
    error_message: '',
    latency: 0,
    created_at: createdAt,
    resolved_at: resolvedAt,
  }
}
