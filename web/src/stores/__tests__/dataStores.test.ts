import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

const fetchMock = vi.hoisted(() => ({
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  delete: vi.fn(),
  uploadForm: vi.fn(),
}))

vi.mock('@/lib/fetch', () => ({
  default: fetchMock,
}))

import { useAnalyticsStore } from '@/stores/analytics'
import { useApiTokensStore } from '@/stores/apiTokens'
import { useAuditLogsStore } from '@/stores/auditLogs'
import { useIntegrationsStore } from '@/stores/integrations'
import { useMaintenancesStore } from '@/stores/maintenances'
import { useNotificationLogsStore } from '@/stores/notificationLogs'
import { useNotificationsStore } from '@/stores/notifications'
import { useSettingsStore } from '@/stores/settings'
import { useSslMonitorsStore } from '@/stores/sslMonitors'
import { useStatusPagesStore } from '@/stores/statusPages'
import { useSubscribersStore } from '@/stores/subscribers'
import { useTagsStore } from '@/stores/tags'
import { useTeamsStore } from '@/stores/teams'
import { useUsersStore } from '@/stores/users'

describe('data stores', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('analytics store fetches dashboard stats and chart data', async () => {
    const store = useAnalyticsStore()
    const stats = [{ monitor_id: 1, points: [] }]
    const points = [{ date: '2026-01-01', status: 'up' }]

    fetchMock.get
      .mockResolvedValueOnce({ data: { data: { data: stats } } })
      .mockResolvedValueOnce({ data: { data: { points } } })

    await store.fetchDashboardStats('7d')
    await store.fetchChart(12, '24h')

    expect(fetchMock.get).toHaveBeenNthCalledWith(1, '/analytics/dashboard?window=7d')
    expect(fetchMock.get).toHaveBeenNthCalledWith(2, '/analytics/monitors/12/chart?window=24h')
    expect(store.stats).toEqual(stats)
    expect(store.chartPoints).toEqual(points)
    expect(store.chartMonitorID).toBe(12)
    expect(store.chartWindow).toBe('24h')
    expect(store.loading).toBe(false)
  })

  it('analytics store clears chart points and records the error when chart fetch fails', async () => {
    const store = useAnalyticsStore()
    const errorSpy = vi.spyOn(console, 'error').mockImplementation(() => {})
    store.chartPoints = [{ date: 'stale' }] as any
    const error = { response: { data: { error: 'chart failed' } } }

    fetchMock.get.mockRejectedValueOnce(error)

    await expect(store.fetchChart(4)).rejects.toBe(error)

    expect(errorSpy).toHaveBeenCalledWith('Failed to fetch chart:', error)
    expect(store.chartPoints).toEqual([])
    expect(store.error).toBe('chart failed')
    expect(store.loading).toBe(false)

    errorSpy.mockRestore()
  })

  it('api tokens store fetches, creates, and revokes tokens', async () => {
    const store = useApiTokensStore()
    const tokens = [{ id: 1, name: 'ops' }]
    const created = { id: 2, raw_token: 'secret' }

    fetchMock.get.mockResolvedValueOnce({ data: { data: tokens } })
    fetchMock.post.mockResolvedValueOnce({ data: { data: created } })
    store.tokens = [{ id: 1, name: 'ops' }] as any

    await store.fetchAll()
    const result = await store.create('deploy', '2026-12-31')
    await store.revoke(1)

    expect(store.tokens).toEqual([])
    expect(result).toEqual(created)
    expect(fetchMock.post).toHaveBeenCalledWith('/api-tokens', {
      name: 'deploy',
      expires_at: '2026-12-31',
    })
    expect(fetchMock.delete).toHaveBeenCalledWith('/api-tokens/1')
  })

  it('audit logs store builds the filter query string and updates current filter', async () => {
    const store = useAuditLogsStore()
    const items = [{ id: 1 }]

    fetchMock.get.mockResolvedValueOnce({ data: { data: items } })

    await store.fetchAll({
      user_id: 7,
      entity_type: 'monitor',
      action: 'delete',
      from: '2026-01-01',
      to: '2026-01-31',
      limit: 25,
    })

    expect(fetchMock.get).toHaveBeenCalledWith(
      '/audit-logs?user_id=7&entity_type=monitor&action=delete&from=2026-01-01&to=2026-01-31&limit=25',
    )
    expect(store.items).toEqual(items)
    expect(store.currentFilter).toEqual({
      user_id: 7,
      entity_type: 'monitor',
      action: 'delete',
      from: '2026-01-01',
      to: '2026-01-31',
      limit: 25,
    })
    expect(store.loading).toBe(false)
  })

  it('integrations store fetches, deletes, and tests integrations', async () => {
    const store = useIntegrationsStore()
    store.integrations = [{ id: 1, name: 'Slack' }] as any

    fetchMock.get.mockResolvedValueOnce({ data: { data: [{ id: 1, name: 'Slack' }] } })
    fetchMock.delete.mockResolvedValueOnce({ data: { success: true } })
    fetchMock.post.mockResolvedValueOnce({ data: { success: true } })

    await store.fetchAll()
    const testResult = await store.test(3)
    await store.deleteIntegration(1)

    expect(testResult).toEqual({ success: true })
    expect(store.integrations).toEqual([])
    expect(fetchMock.post).toHaveBeenCalledWith('/integrations/3/test')
  })

  it('maintenances store fetches items and removes a maintenance from local state', async () => {
    const store = useMaintenancesStore()
    store.items = [{ id: 1 }, { id: 2 }] as any
    fetchMock.get.mockResolvedValueOnce({ data: { data: [{ id: 1 }] } })

    await store.fetchAll()
    await store.remove(2)

    expect(store.items).toEqual([{ id: 1 }])
    expect(fetchMock.delete).toHaveBeenCalledWith('/maintenances/2')
  })

  it('notification logs store only includes truthy filters in the request URL', async () => {
    const store = useNotificationLogsStore()
    const logs = [{ id: 1 }]
    fetchMock.get.mockResolvedValueOnce({ data: { data: logs } })

    await store.fetchLogs({ status: 'failed', recipient: '', type: 'email' })

    expect(fetchMock.get).toHaveBeenCalledWith('/notification-logs?status=failed&type=email')
    expect(store.logs).toEqual(logs)
  })

  it('notifications store fetches channels and filters deleted channels from state', async () => {
    const store = useNotificationsStore()
    store.channels = [{ id: 1 }, { id: 2 }] as any

    fetchMock.get.mockResolvedValueOnce({ data: { data: [{ id: 1 }, { id: 2 }] } })
    fetchMock.delete.mockResolvedValueOnce({ data: { success: true } })

    await store.fetchChannelsData()
    const response = await store.deleteChannel(1)

    expect(response).toEqual({ success: true })
    expect(store.channels).toEqual([{ id: 2 }])
  })

  it('settings store fetches current settings and posts each payload when saving', async () => {
    const store = useSettingsStore()
    const settings = [{ key: 'system_name', value: 'Ping' }]
    fetchMock.get.mockResolvedValueOnce({ data: { data: settings } })
    fetchMock.post.mockResolvedValueOnce({}).mockResolvedValueOnce({})

    await store.fetchSettingsData()
    await store.saveSettingsData([
      { key: 'system_name', value: 'Ping', description: '' },
      { key: 'admin_email', value: 'ops@example.com', description: '' },
    ])

    expect(store.settings).toEqual(settings)
    expect(fetchMock.post).toHaveBeenNthCalledWith(1, '/settings', {
      key: 'system_name',
      value: 'Ping',
      description: '',
    })
    expect(fetchMock.post).toHaveBeenNthCalledWith(2, '/settings', {
      key: 'admin_email',
      value: 'ops@example.com',
      description: '',
    })
    expect(store.saveLoading).toBe(false)
  })

  it('ssl monitors store fetches expiring certs, updates a checked cert, and removes deleted certs', async () => {
    const store = useSslMonitorsStore()
    store.certs = [{ id: 1, monitor_id: 9, domain: 'old.example.com' }] as any

    fetchMock.get.mockResolvedValueOnce({ data: { data: [{ id: 2, monitor_id: 2 }] } })
    fetchMock.post
      .mockResolvedValueOnce({
        data: { data: { id: 1, monitor_id: 9, domain: 'new.example.com' } },
      })
      .mockResolvedValueOnce({ data: { data: [{ id: 3 }] } })
    fetchMock.get.mockResolvedValueOnce({ data: { data: [{ id: 4, monitor_id: 4 }] } })

    await store.fetchExpiring(14)
    const checked = await store.checkMonitor(9)
    const checkAllResult = await store.checkAll()
    await store.deleteCert(4)

    expect(fetchMock.get).toHaveBeenNthCalledWith(1, '/ssl-monitors/expiring?days=14')
    expect(checked).toEqual({ id: 1, monitor_id: 9, domain: 'new.example.com' })
    expect(checkAllResult).toEqual([{ id: 3 }])
    expect(store.certs).toEqual([])
  })

  it('status pages store fetches both pages and monitors, then removes deleted pages from state', async () => {
    const store = useStatusPagesStore()
    store.statusPages = [{ id: 1 }, { id: 2 }] as any

    fetchMock.get
      .mockResolvedValueOnce({ data: { data: [{ id: 1, name: 'Public' }] } })
      .mockResolvedValueOnce({ data: { data: [{ id: 9, name: 'API' }] } })

    await store.fetchStatusPagesData()
    await store.deleteStatusPage(2)

    expect(store.statusPages).toEqual([{ id: 1, name: 'Public' }])
    expect(store.availableMonitors).toEqual([{ id: 9, name: 'API' }])
    expect(fetchMock.delete).toHaveBeenCalledWith('/status-pages/2')
  })

  it('subscribers store fetches subscribers, maps status pages, and removes unsubscribed emails', async () => {
    const store = useSubscribersStore()
    store.subscribers = [{ email: 'ops@example.com' }, { email: 'dev@example.com' }] as any

    fetchMock.get
      .mockResolvedValueOnce({ data: { data: [{ id: 1, email: 'ops@example.com' }] } })
      .mockResolvedValueOnce({ data: { data: [{ id: 4, name: 'Public' }] } })

    await store.fetchSubscribers(4)
    await store.fetchStatusPages()
    await store.unsubscribe(4, 'ops@example.com')

    expect(store.subscribers).toEqual([])
    expect(store.statusPages).toEqual([{ id: 4, name: 'Public' }])
    expect(fetchMock.post).toHaveBeenCalledWith('/status-pages/4/subscribers/unsubscribe', {
      email: 'ops@example.com',
    })
  })

  it('tags store mutates local tag state and calls monitor tag endpoints', async () => {
    const store = useTagsStore()
    store.tags = [{ id: 1, name: 'ops' }] as any
    fetchMock.post.mockResolvedValueOnce({ data: { data: { id: 2, name: 'prod' } } })
    fetchMock.put.mockResolvedValueOnce({ data: { data: { id: 1, name: 'Ops Team' } } })
    fetchMock.get.mockResolvedValueOnce({ data: { data: [{ id: 1, name: 'Ops Team' }] } })

    const created = await store.create({ name: 'prod' })
    const updated = await store.update(1, { name: 'Ops Team' })
    await store.attachToMonitor(8, [1, 2])
    const monitorTags = await store.getMonitorTags(8)
    await store.detachFromMonitor(8, 1)
    await store.deleteTag(2)

    expect(created).toEqual({ id: 2, name: 'prod' })
    expect(updated).toEqual({ id: 1, name: 'Ops Team' })
    expect(monitorTags).toEqual([{ id: 1, name: 'Ops Team' }])
    expect(store.tags).toEqual([{ id: 1, name: 'Ops Team' }])
    expect(fetchMock.post).toHaveBeenCalledWith('/monitors/8/tags', { tag_ids: [1, 2] })
    expect(fetchMock.delete).toHaveBeenCalledWith('/monitors/8/tags/1')
  })

  it('teams store manages teams and team members through local state updates', async () => {
    const store = useTeamsStore()
    store.teams = [{ id: 1, name: 'Ops' }] as any
    store.members = [{ id: 9, role: 'member' }] as any

    fetchMock.post
      .mockResolvedValueOnce({ data: { data: { id: 2, name: 'Platform' } } })
      .mockResolvedValueOnce({ data: { data: { id: 11, role: 'member' } } })
      .mockResolvedValueOnce({})
      .mockResolvedValueOnce({})
    fetchMock.put
      .mockResolvedValueOnce({ data: { data: { id: 1, name: 'Operations' } } })
      .mockResolvedValueOnce({ data: { data: { id: 9, role: 'admin' } } })

    await store.create('Platform')
    await store.update(1, 'Operations')
    await store.inviteMember(1, 42)
    await store.updateMember(1, 9, 'admin')
    await store.removeMember(1, 11)
    await store.deleteTeam(2)
    await store.acceptInvite(1)
    await store.rejectInvite(1)

    expect(store.teams).toEqual([{ id: 1, name: 'Operations' }])
    expect(store.members).toEqual([{ id: 9, role: 'admin' }])
    expect(fetchMock.post).toHaveBeenCalledWith('/teams/1/members/accept')
    expect(fetchMock.post).toHaveBeenCalledWith('/teams/1/members/reject')
  })

  it('users store fetches users and removes a deleted user from state', async () => {
    const store = useUsersStore()
    store.users = [{ id: 1 }, { id: 2 }] as any

    fetchMock.get.mockResolvedValueOnce({ data: { data: [{ id: 1 }, { id: 2 }] } })
    fetchMock.delete.mockResolvedValueOnce({ data: { success: true } })

    await store.fetchUsersData()
    const result = await store.deleteUser(1)

    expect(result).toEqual({ success: true })
    expect(store.users).toEqual([{ id: 2 }])
  })
})
