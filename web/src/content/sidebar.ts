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
        href: "/dashboard",
        icon: LayoutDashboard,
      },
      {
        title: "Monitors",
        href: "/dashboard/monitors",
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
        href: "/dashboard/incidents",
        icon: FileText,
      },
      {
        title: "Status Pages",
        href: "/dashboard/status-pages",
        icon: ShieldCheck,
      },
    ],
  },
  {
    groupName: "Administration",
    items: [
      {
        title: "Users",
        href: "/dashboard/users",
        icon: Users,
      },
      {
        title: "Notifications",
        href: "/dashboard/settings/notifications",
        icon: Bell,
      },
      {
        title: "Settings",
        href: "/dashboard/settings",
        icon: Settings,
      },
    ],
  },
];
