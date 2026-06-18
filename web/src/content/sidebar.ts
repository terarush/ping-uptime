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
    groupName: 'Overview',
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
        title: 'Status Pages',
        href: `${siteConfig.appPath}/status-pages`,
        icon: Globe,
      },
      {
        title: 'Analytics',
        href: `${siteConfig.appPath}/analytics`,
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
        title: 'Maintenance',
        href: `${siteConfig.appPath}/maintenances`,
        icon: Wrench,
      },
    ],
  },
  {
    groupName: 'Preferences',
    items: [
      {
        title: 'Notifications',
        href: `${siteConfig.appPath}/settings/notifications`,
        icon: Bell,
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
        title: 'Audit Logs',
        href: `${siteConfig.appPath}/audit-logs`,
        icon: ClipboardList,
      },
      {
        title: 'Settings',
        href: `${siteConfig.appPath}/settings`,
        icon: Settings,
      },
    ],
  },
]
