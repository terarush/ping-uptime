import {
  LayoutDashboard,
  Activity,
  Settings,
  Users,
  Bell,
  FileText,
  ShieldCheck
} from '@lucide/vue';
import { siteConfig } from './config';

export interface SidebarItem {
  title: string;
  href: string;
  icon: any;
  badge?: string;
}

export interface SidebarGroup {
  groupName?: string;
  items: SidebarItem[];
}

export const sidebarContent: SidebarGroup[] = [
  {
    groupName: "Overview",
    items: [
      {
        title: "Dashboard",
        href: `${siteConfig.appPath}`,
        icon: LayoutDashboard,
      },
      {
        title: "Monitors",
        href: `${siteConfig.appPath}/monitors`,
        icon: Activity,
        badge: "Live",
      },
    ],
  },
  {
    groupName: "Alerts & Incidents",
    items: [
      {
        title: "Incident Logs",
        href: `${siteConfig.appPath}/incidents`,
        icon: FileText,
      },
      {
        title: "Status Pages",
        href: `${siteConfig.appPath}/status-pages`,
        icon: ShieldCheck,
      },
    ],
  },
  {
    groupName: "Administration",
    items: [
      {
        title: "Users",
        href: `${siteConfig.appPath}/users`,
        icon: Users,
      },
      {
        title: "Notifications",
        href: `${siteConfig.appPath}/settings/notifications`,
        icon: Bell,
      },
      {
        title: "Settings",
        href: `${siteConfig.appPath}/settings`,
        icon: Settings,
      },
    ],
  },
];
