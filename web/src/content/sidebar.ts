import {
  LayoutDashboard,
  Activity,
  Settings,
  Users,
  Bell,
  FileText,
  ShieldCheck,
  BarChart3,
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
        title: 'Settings',
        href: `${siteConfig.appPath}/settings`,
        icon: Settings,
      },
    ],
  },
]
