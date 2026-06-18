import {
  LayoutDashboard,
  Activity,
  Settings,
  Users,
  Bell,
  FileText,
  ShieldCheck,
  BarChart3,
  Globe,
  Wrench,
  ClipboardList,
  Tag,
  KeyRound,
  Mail,
  UsersRound,
  Link,
  FolderArchive,
  UserPlus,
} from '@lucide/vue'
import { siteConfig } from './config'

export interface SidebarItem {
  title: string
  href: string
  icon: any
  badge?: string
}

export interface SidebarGroup {
  groupName?: string
  admin?: boolean
  items: SidebarItem[]
}
export const sidebarContent: SidebarGroup[] = [
  {
    groupName: 'Monitoring',
    items: [
      {
        title: 'Dashboard',
        href: `${siteConfig.appPath}`,
        icon: LayoutDashboard,
      },
      {
        title: 'Monitors',
        href: `${siteConfig.appPath}/monitors`,
        icon: Activity,
        badge: 'Live',
      },
      {
        title: 'SSL Certificates',
        href: `${siteConfig.appPath}/ssl-monitors`,
        icon: ShieldCheck,
        badge: 'Expiry',
      },
    ],
  },
  {
    groupName: 'Status & Reports',
    items: [
      {
        title: 'Status Pages',
        href: `${siteConfig.appPath}/status-pages`,
        icon: Globe,
      },
      {
        title: 'Subscribers',
        href: `${siteConfig.appPath}/subscribers`,
        icon: Mail,
      },
      {
        title: 'Reports',
        href: `${siteConfig.appPath}/reports`,
        icon: BarChart3,
      },
    ],
  },
  {
    groupName: 'Alerts & Incidents',
    items: [
      {
        title: 'Incident Logs',
        href: `${siteConfig.appPath}/incidents`,
        icon: FileText,
      },
      {
        title: 'Maintenances',
        href: `${siteConfig.appPath}/maintenances`,
        icon: Wrench,
      },
      {
        title: 'Notification History',
        href: `${siteConfig.appPath}/notification-logs`,
        icon: ClipboardList,
      },
    ],
  },
  {
    groupName: 'Notifications',
    items: [
      {
        title: 'Channels',
        href: `${siteConfig.appPath}/settings/notifications`,
        icon: Bell,
      },
      {
        title: 'Integrations',
        href: `${siteConfig.appPath}/integrations`,
        icon: Link,
      },
    ],
  },
  {
    groupName: 'Administration',
    admin: true,
    items: [
      {
        title: 'Users',
        href: `${siteConfig.appPath}/users`,
        icon: Users,
      },
      {
        title: 'Teams',
        href: `${siteConfig.appPath}/teams`,
        icon: UsersRound,
      },
      {
        title: 'Audit Logs',
        href: `${siteConfig.appPath}/audit-logs`,
        icon: FileText,
      },
      {
        title: 'API Tokens',
        href: `${siteConfig.appPath}/api-tokens`,
        icon: KeyRound,
      },
      {
        title: 'Tags & Labels',
        href: `${siteConfig.appPath}/tags`,
        icon: Tag,
      },
      {
        title: 'Backup & Export',
        href: `${siteConfig.appPath}/backup`,
        icon: FolderArchive,
      },
      {
        title: 'Settings',
        href: `${siteConfig.appPath}/settings`,
        icon: Settings,
      },
    ],
  },
]
