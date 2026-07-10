import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

const fetchMock = vi.hoisted(() => ({
  get: vi.fn(),
  delete: vi.fn(),
  uploadForm: vi.fn(),
}))

vi.mock('@/lib/fetch', () => ({
  default: fetchMock,
}))

import { useBackupStore } from '@/stores/backup'

describe('backup store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('exports backups and refreshes history after importing a file', async () => {
    const store = useBackupStore()
    const file = new File(['backup'], 'backup.json', { type: 'application/json' })

    fetchMock.get.mockResolvedValueOnce({ data: { data: { url: '/download' } } })
    fetchMock.uploadForm.mockResolvedValueOnce({ success: true })
    fetchMock.get.mockResolvedValueOnce({ data: { data: [{ id: 1, file_name: 'backup.json' }] } })

    await expect(store.exportBackup()).resolves.toEqual({ url: '/download' })
    await expect(store.importBackup(file)).resolves.toEqual({ success: true })

    expect(fetchMock.uploadForm).toHaveBeenCalledWith('/backup/import', expect.any(FormData))
    expect(store.records).toEqual([{ id: 1, file_name: 'backup.json' }])
    expect(store.loading).toBe(false)
  })

  it('downloads backups using the filename from the content disposition header', async () => {
    const store = useBackupStore()
    const click = vi.fn()
    const revokeObjectURL = vi.fn()
    const createObjectURL = vi.fn().mockReturnValue('blob:download')
    const createElement = vi.spyOn(document, 'createElement').mockReturnValue({
      click,
      href: '',
      download: '',
    } as any)
    vi.stubGlobal('URL', {
      createObjectURL,
      revokeObjectURL,
    })

    fetchMock.get.mockResolvedValueOnce({
      data: new Blob(['backup']),
      headers: { 'content-disposition': 'attachment; filename=custom-backup.json' },
    })

    await store.downloadBackup()

    expect(fetchMock.get).toHaveBeenCalledWith('/backup/download', { responseType: 'blob' })
    expect(createElement).toHaveBeenCalledWith('a')
    expect(click).toHaveBeenCalled()
    expect(createObjectURL).toHaveBeenCalled()
    expect(revokeObjectURL).toHaveBeenCalledWith('blob:download')
  })

  it('loads and deletes backup history records', async () => {
    const store = useBackupStore()
    store.records = [{ id: 1 }, { id: 2 }] as any

    fetchMock.get.mockResolvedValueOnce({ data: { data: [{ id: 1 }, { id: 2 }] } })

    await store.getHistory()
    await store.deleteRecord(1)

    expect(fetchMock.delete).toHaveBeenCalledWith('/backup/history/1')
    expect(store.records).toEqual([{ id: 2 }])
  })
})
