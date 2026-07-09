import { afterEach, describe, expect, it, vi } from 'vitest'

describe('useAppTitle', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
    vi.resetModules()
  })

  it('starts with the configured site name', async () => {
    const { useAppTitle } = await import('@/composables/useAppTitle')
    const { siteConfig } = await import('@/content/config')

    expect(useAppTitle().appTitle.value).toBe(siteConfig.name)
  })

  it('updates the shared title from the public settings endpoint', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: true,
        json: vi.fn().mockResolvedValue({ data: { value: 'Custom Uptime' } }),
      }),
    )

    const { useAppTitle } = await import('@/composables/useAppTitle')
    const firstConsumer = useAppTitle()
    const secondConsumer = useAppTitle()

    await firstConsumer.fetchAppTitle()

    expect(fetch).toHaveBeenCalledWith('/api/settings/public/system-name')
    expect(firstConsumer.appTitle.value).toBe('Custom Uptime')
    expect(secondConsumer.appTitle.value).toBe('Custom Uptime')
  })

  it('keeps the default title when the fetch fails', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('network down')))

    const { useAppTitle } = await import('@/composables/useAppTitle')
    const { siteConfig } = await import('@/content/config')
    const { appTitle, fetchAppTitle } = useAppTitle()

    await fetchAppTitle()

    expect(appTitle.value).toBe(siteConfig.name)
  })
})
