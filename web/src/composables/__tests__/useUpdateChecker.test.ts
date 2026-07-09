import { afterEach, describe, expect, it, vi } from 'vitest'

describe('useUpdateChecker', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
    vi.restoreAllMocks()
    vi.resetModules()
  })

  it('fetches the latest release and clears loading state', async () => {
    const release = {
      tag_name: 'v9.9.9',
      html_url: 'https://github.com/terarush/ping-uptime/releases/tag/v9.9.9',
      published_at: '2026-01-01T00:00:00Z',
      prerelease: false,
    }
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: true,
        json: vi.fn().mockResolvedValue(release),
      }),
    )
    vi.spyOn(console, 'warn').mockImplementation(() => {})

    const { useUpdateChecker } = await import('@/composables/useUpdateChecker')
    const checker = useUpdateChecker()

    await checker.checkForUpdate()

    expect(fetch).toHaveBeenCalledWith(
      'https://api.github.com/repos/terarush/ping-uptime/releases/latest',
      { headers: { Accept: 'application/vnd.github.v3+json' } },
    )
    expect(checker.latestRelease.value).toEqual(release)
    expect(checker.loading.value).toBe(false)
    expect(checker.error.value).toBe(false)
    expect(checker.isUpdateAvailable.value).toBe(true)
  })

  it('sets error state when the release request fails', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: false,
      }),
    )
    vi.spyOn(console, 'warn').mockImplementation(() => {})

    const { useUpdateChecker } = await import('@/composables/useUpdateChecker')
    const checker = useUpdateChecker()

    await checker.checkForUpdate()

    expect(checker.error.value).toBe(true)
    expect(checker.loading.value).toBe(false)
  })
})
