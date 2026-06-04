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
        href: "/",
        icon: LayoutDashboard,
      },
      {
        title: "Monitors",
        href: "/monitors",
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
        href: "/incidents",
        icon: FileText,
      },
      {
        title: "Status Pages",
        href: "/status-pages",
        icon: ShieldCheck,
      },
    ],
  },
  {
    groupName: "Administration",
    items: [
      {
        title: "Users",
        href: "/users",
        icon: Users,
      },
      {
        title: "Notifications",
        href: "/settings/notifications",
        icon: Bell,
      },
      {
        title: "Settings",
        href: "/settings",
        icon: Settings,
      },
    ],
  },
];
