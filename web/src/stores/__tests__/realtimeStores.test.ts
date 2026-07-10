import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

const fetchMock = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  delete: vi.fn(),
}))

vi.mock('@/lib/fetch', () => ({
  default: fetchMock,
}))

import { useIncidentsStore } from '@/stores/incidents'
import { useMonitorsStore } from '@/stores/monitors'

describe('realtime stores', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('incident store fetches monitors and incidents, then starts monitor SSE', async () => {
    const incidentsStore = useIncidentsStore()
    const monitorsStore = useMonitorsStore()
    const initializeSSE = vi.spyOn(monitorsStore, 'initializeSSE').mockImplementation(() => {})

    fetchMock.get
      .mockResolvedValueOnce({
        data: {
          data: [
            { id: 1, name: 'API', url: 'https://example.com' },
            { id: 2, name: 'Site', url: 'https://site.example.com' },
          ],
        },
      })
      .mockResolvedValueOnce({
        data: { data: [{ id: 9, monitor_id: 1, status: 'active' }] },
      })

    await incidentsStore.fetchIncidentsData()

    expect(incidentsStore.monitorsMap).toEqual({
      1: { id: 1, name: 'API', url: 'https://example.com' },
      2: { id: 2, name: 'Site', url: 'https://site.example.com' },
    })
    expect(incidentsStore.incidents).toEqual([{ id: 9, monitor_id: 1, status: 'active' }])
    expect(initializeSSE).toHaveBeenCalled()
    expect(incidentsStore.loading).toBe(false)
  })

  it('monitor store fetches monitors and initializes the SSE connection once', async () => {
    const instances: MockEventSource[] = []

    class MockEventSource {
      url: string
      onmessage: ((event: MessageEvent) => void) | null = null
      onerror: ((event: Event) => void) | null = null

      constructor(url: string) {
        this.url = url
        instances.push(this)
      }
    }

    vi.stubGlobal('EventSource', MockEventSource as any)
    Object.defineProperty(window, 'location', {
      value: { origin: 'http://localhost:5173' },
      writable: true,
    })

    const store = useMonitorsStore()
    fetchMock.get.mockResolvedValueOnce({ data: { data: [{ id: 1, name: 'API' }] } })

    await store.fetchMonitors()
    store.initializeSSE()

    expect(store.monitors).toEqual([{ id: 1, name: 'API' }])
    expect(instances).toHaveLength(1)
    expect(instances[0]!.url).toBe('http://localhost:5173/api/monitors/events')
  })

  it('monitor SSE updates monitors and incident state from incoming events', () => {
    let instance: MockRealtimeEventSource | undefined

    class MockRealtimeEventSource {
      onmessage: ((event: MessageEvent) => void) | null = null
      onerror: ((event: Event) => void) | null = null

      constructor(_url: string) {
        instance = this
      }
    }

    vi.stubGlobal('EventSource', MockRealtimeEventSource as any)
    Object.defineProperty(window, 'location', {
      value: { origin: 'http://localhost:5173' },
      writable: true,
    })

    const monitorsStore = useMonitorsStore()
    const incidentsStore = useIncidentsStore()

    monitorsStore.monitors = [{ id: 1, name: 'Old API', uptime_status: 'up' }] as any
    incidentsStore.incidents = [{ id: 7, status: 'active' }] as any

    monitorsStore.initializeSSE()

    instance?.onmessage?.({
      data: JSON.stringify({
        type: 'monitor.updated',
        payload: { id: 1, name: 'New API', uptime_status: 'down' },
      }),
    } as MessageEvent)
    instance?.onmessage?.({
      data: JSON.stringify({
        type: 'incident.created',
        payload: { id: 8, status: 'active' },
      }),
    } as MessageEvent)
    instance?.onmessage?.({
      data: JSON.stringify({
        type: 'incident.resolved',
        payload: { id: 7, status: 'resolved' },
      }),
    } as MessageEvent)
    instance?.onmessage?.({
      data: JSON.stringify({
        type: 'monitor.deleted',
        payload: { id: 1 },
      }),
    } as MessageEvent)

    expect(monitorsStore.monitors).toEqual([])
    expect(incidentsStore.incidents[0]).toEqual({ id: 8, status: 'active' })
    expect(incidentsStore.incidents[1]).toEqual({ id: 7, status: 'resolved' })
  })

  it('monitor store returns created and updated monitors and removes deleted ones', async () => {
    const store = useMonitorsStore()
    store.monitors = [{ id: 1 }, { id: 2 }] as any

    fetchMock.post.mockResolvedValueOnce({ data: { data: { id: 3, name: 'New' } } })
    fetchMock.put.mockResolvedValueOnce({ data: { data: { id: 2, name: 'Updated' } } })

    const created = await store.createMonitor({ name: 'New' } as any)
    const updated = await store.updateMonitor(2, { name: 'Updated' })
    await store.deleteMonitor(1)

    expect(created).toEqual({ id: 3, name: 'New' })
    expect(updated).toEqual({ id: 2, name: 'Updated' })
    expect(fetchMock.delete).toHaveBeenCalledWith('/monitors/1')
    expect(store.monitors).toEqual([{ id: 2 }])
  })
})
