import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

const fetchMock = vi.hoisted(() => ({
  get: vi.fn(),
}))

const cookiesMock = vi.hoisted(() => ({
  get: vi.fn(),
  set: vi.fn(),
  remove: vi.fn(),
}))

vi.mock('@/lib/fetch', () => ({
  default: fetchMock,
}))

vi.mock('js-cookie', () => ({
  default: cookiesMock,
}))

import { useAuthStore } from '@/stores/auth'

describe('auth store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('returns false without verifying when no access token cookie exists', async () => {
    const store = useAuthStore()
    cookiesMock.get.mockReturnValueOnce(undefined)

    await expect(store.verifyToken()).resolves.toBe(false)

    expect(fetchMock.get).not.toHaveBeenCalled()
    expect(store.currentUser).toBeNull()
    expect(store.isAuthenticated).toBe(false)
  })

  it('verifies the token and stores the returned user', async () => {
    const store = useAuthStore()
    const user = { id: 1, email: 'ops@example.com' }

    cookiesMock.get.mockReturnValueOnce('token')
    fetchMock.get.mockResolvedValueOnce({ data: { data: user } })

    await expect(store.verifyToken()).resolves.toBe(true)

    expect(fetchMock.get).toHaveBeenCalledWith('/users/verify')
    expect(store.currentUser).toEqual(user)
    expect(store.isAuthenticated).toBe(true)
    expect(store.isVerifying).toBe(false)
  })

  it('clears the session and cookies when token verification fails', async () => {
    const store = useAuthStore()
    const warnSpy = vi.spyOn(console, 'warn').mockImplementation(() => {})

    cookiesMock.get.mockReturnValueOnce('bad-token')
    fetchMock.get.mockRejectedValueOnce(new Error('unauthorized'))

    try {
      await expect(store.verifyToken()).resolves.toBe(false)

      expect(warnSpy).toHaveBeenCalledWith(
        'Token verification failed on server, clearing session:',
        expect.any(Error),
      )
      expect(cookiesMock.remove).toHaveBeenCalledWith('accessToken')
      expect(cookiesMock.remove).toHaveBeenCalledWith('refreshToken')
      expect(store.currentUser).toBeNull()
      expect(store.isAuthenticated).toBe(false)
    } finally {
      warnSpy.mockRestore()
    }
  })

  it('sets and clears the session through cookies', () => {
    const store = useAuthStore()
    const user = { id: 2, email: 'dev@example.com' } as any

    store.setSession('access', 'refresh', user)

    expect(cookiesMock.set).toHaveBeenNthCalledWith(1, 'accessToken', 'access', { expires: 7 })
    expect(cookiesMock.set).toHaveBeenNthCalledWith(2, 'refreshToken', 'refresh', { expires: 7 })
    expect(store.currentUser).toEqual(user)
    expect(store.isAuthenticated).toBe(true)

    store.logout()

    expect(cookiesMock.remove).toHaveBeenCalledWith('accessToken')
    expect(cookiesMock.remove).toHaveBeenCalledWith('refreshToken')
    expect(store.currentUser).toBeNull()
    expect(store.isAuthenticated).toBe(false)
  })
})
