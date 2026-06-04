import {
  LayoutDashboard,
  Activity,
  Settings,
  Users,
  Bell,
  FileText,
  ShieldCheck
} from '@lucide/vue';

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
        href: "/app",
        icon: LayoutDashboard,
      },
      {
        title: "Monitors",
        href: "/app/monitors",
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
        href: "/app/incidents",
        icon: FileText,
      },
      {
        title: "Status Pages",
        href: "/app/status-pages",
        icon: ShieldCheck,
      },
    ],
  },
  {
    groupName: "Administration",
    items: [
      {
        title: "Users",
        href: "/app/users",
        icon: Users,
      },
      {
        title: "Notifications",
        href: "/app/settings/notifications",
        icon: Bell,
      },
      {
        title: "Settings",
        href: "/app/settings",
        icon: Settings,
      },
    ],
  },
];
