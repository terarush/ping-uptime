import { beforeEach, describe, expect, it } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'

import { useAnalytics } from '@/composables/useAnalytics'
import { useApiTokens } from '@/composables/useApiTokens'
import { useAuditLogs } from '@/composables/useAuditLogs'
import { useAuth } from '@/composables/useAuth'
import { useBackup } from '@/composables/useBackup'
import { useIncidents } from '@/composables/useIncidents'
import { useIntegrations } from '@/composables/useIntegrations'
import { useMaintenances } from '@/composables/useMaintenances'
import { useMonitors } from '@/composables/useMonitors'
import { useNotificationChannels } from '@/composables/useNotificationChannels'
import { useNotificationLogs } from '@/composables/useNotificationLogs'
import { useSettings } from '@/composables/useSettings'
import { useSslMonitors } from '@/composables/useSslMonitors'
import { useStatusPages } from '@/composables/useStatusPages'
import { useSubscribers } from '@/composables/useSubscribers'
import { useTags } from '@/composables/useTags'
import { useTeams } from '@/composables/useTeams'
import { useUsers } from '@/composables/useUsers'
import { useAnalyticsStore } from '@/stores/analytics'
import { useApiTokensStore } from '@/stores/apiTokens'
import { useAuditLogsStore } from '@/stores/auditLogs'
import { useAuthStore } from '@/stores/auth'
import { useBackupStore } from '@/stores/backup'
import { useIncidentsStore } from '@/stores/incidents'
import { useIntegrationsStore } from '@/stores/integrations'
import { useMaintenancesStore } from '@/stores/maintenances'
import { useMonitorsStore } from '@/stores/monitors'
import { useNotificationsStore } from '@/stores/notifications'
import { useNotificationLogsStore } from '@/stores/notificationLogs'
import { useSettingsStore } from '@/stores/settings'
import { useSslMonitorsStore } from '@/stores/sslMonitors'
import { useStatusPagesStore } from '@/stores/statusPages'
import { useSubscribersStore } from '@/stores/subscribers'
import { useTagsStore } from '@/stores/tags'
import { useTeamsStore } from '@/stores/teams'
import { useUsersStore } from '@/stores/users'

interface BindingSpec {
  name: string
  useComposable: () => Record<string, any>
  useStore: () => Record<string, any>
  refs: string[]
  actions: Array<string | [string, string]>
}

const specs: BindingSpec[] = [
  {
    name: 'useAnalytics',
    useComposable: useAnalytics,
    useStore: useAnalyticsStore,
    refs: ['stats', 'chartPoints', 'chartMonitorID', 'chartWindow', 'loading', 'error'],
    actions: ['fetchDashboardStats', 'fetchChart'],
  },
  {
    name: 'useApiTokens',
    useComposable: useApiTokens,
    useStore: useApiTokensStore,
    refs: ['tokens', 'loading', 'error'],
    actions: ['fetchAll', 'create', 'revoke'],
  },
  {
    name: 'useAuditLogs',
    useComposable: useAuditLogs,
    useStore: useAuditLogsStore,
    refs: ['items', 'loading'],
    actions: ['fetchAll'],
  },
  {
    name: 'useAuth',
    useComposable: useAuth,
    useStore: useAuthStore,
    refs: ['currentUser', 'isAuthenticated', 'isVerifying'],
    actions: ['verifyToken', 'logout', 'setSession'],
  },
  {
    name: 'useBackup',
    useComposable: useBackup,
    useStore: useBackupStore,
    refs: ['records', 'loading'],
    actions: ['exportBackup', 'downloadBackup', 'importBackup', 'getHistory', 'deleteRecord'],
  },
  {
    name: 'useIncidents',
    useComposable: useIncidents,
    useStore: useIncidentsStore,
    refs: ['incidents', 'monitorsMap', 'loading', 'error'],
    actions: ['fetchIncidentsData'],
  },
  {
    name: 'useIntegrations',
    useComposable: useIntegrations,
    useStore: useIntegrationsStore,
    refs: ['integrations', 'loading', 'error'],
    actions: ['fetchAll', 'create', 'update', 'deleteIntegration', 'test'],
  },
  {
    name: 'useMaintenances',
    useComposable: useMaintenances,
    useStore: useMaintenancesStore,
    refs: ['items', 'loading', 'error'],
    actions: ['fetchAll', 'create', 'update', 'remove'],
  },
  {
    name: 'useMonitors',
    useComposable: useMonitors,
    useStore: useMonitorsStore,
    refs: ['monitors', 'loading', 'error'],
    actions: ['fetchMonitors', 'createMonitor', 'updateMonitor', 'deleteMonitor'],
  },
  {
    name: 'useNotificationChannels',
    useComposable: useNotificationChannels,
    useStore: useNotificationsStore,
    refs: ['channels', 'loading', 'error'],
    actions: [
      ['fetchChannelsData', 'fetchChannelsData'],
      ['createChannel', 'createChannel'],
      ['updateChannel', 'updateChannel'],
      ['deleteChannel', 'deleteChannel'],
    ],
  },
  {
    name: 'useNotificationLogs',
    useComposable: useNotificationLogs,
    useStore: useNotificationLogsStore,
    refs: ['logs', 'loading', 'error'],
    actions: ['fetchLogs'],
  },
  {
    name: 'useSettings',
    useComposable: useSettings,
    useStore: useSettingsStore,
    refs: ['settings', 'loading', 'saveLoading', 'error'],
    actions: ['fetchSettingsData', 'saveSettingsData'],
  },
  {
    name: 'useSslMonitors',
    useComposable: useSslMonitors,
    useStore: useSslMonitorsStore,
    refs: ['certs', 'loading', 'error'],
    actions: ['fetchAll', 'fetchExpiring', 'checkMonitor', 'checkAll', 'deleteCert'],
  },
  {
    name: 'useStatusPages',
    useComposable: useStatusPages,
    useStore: useStatusPagesStore,
    refs: ['statusPages', 'availableMonitors', 'loading', 'error'],
    actions: ['fetchStatusPagesData', 'createStatusPage', 'updateStatusPage', 'deleteStatusPage'],
  },
  {
    name: 'useSubscribers',
    useComposable: useSubscribers,
    useStore: useSubscribersStore,
    refs: ['subscribers', 'statusPages', 'loading', 'error'],
    actions: ['fetchSubscribers', 'fetchStatusPages', 'unsubscribe'],
  },
  {
    name: 'useTags',
    useComposable: useTags,
    useStore: useTagsStore,
    refs: ['tags', 'loading', 'error'],
    actions: [
      'fetchAll',
      'create',
      'update',
      'deleteTag',
      'attachToMonitor',
      'getMonitorTags',
      'detachFromMonitor',
    ],
  },
  {
    name: 'useTeams',
    useComposable: useTeams,
    useStore: useTeamsStore,
    refs: ['teams', 'members', 'loading', 'error'],
    actions: [
      'fetchAll',
      'create',
      'update',
      'deleteTeam',
      'fetchMembers',
      'inviteMember',
      'updateMember',
      'removeMember',
      'acceptInvite',
      'rejectInvite',
    ],
  },
  {
    name: 'useUsers',
    useComposable: useUsers,
    useStore: useUsersStore,
    refs: ['users', 'loading', 'error'],
    actions: ['fetchUsersData', 'createUser', 'updateUser', 'deleteUser'],
  },
]

describe('store-backed composables', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it.each(specs)('$name exposes the expected store refs and actions', (spec) => {
    const exposed = spec.useComposable()
    const store = spec.useStore()

    for (const refName of spec.refs) {
      expect(exposed[refName].value).toEqual(store[refName])
    }

    for (const action of spec.actions) {
      const [exposedName, storeName] = Array.isArray(action) ? action : [action, action]

      expect(exposed[exposedName]).toBe(store[storeName])
    }
  })
})
