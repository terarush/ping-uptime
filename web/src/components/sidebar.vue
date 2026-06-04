<script setup lang="ts">
import {
  Sidebar,
  SidebarHeader,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarGroupContent,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  SidebarMenuBadge,
} from '@/components/ui/sidebar'
import { sidebarContent } from '@/content/sidebar'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import ThemeButton from '@/components/theme-button.vue'
import { LogOut } from '@lucide/vue'
import { siteConfig } from '@/content/config'
import { useAuth } from '@/composables/useAuth'

const route = useRoute()
const router = useRouter()
const { logout } = useAuth()

const handleLogout = () => {
  logout()
  router.push('/')
}
</script>


<template>
  <Sidebar collapsible="icon" variant="sidebar">
    <!-- Header: App Logo -->
    <SidebarHeader class="border-b border-border/50 py-4 px-4 flex items-center justify-between">
      <div
        class="flex items-center gap-2 overflow-hidden transition-[width] group-data-[collapsible=icon]:w-0 group-data-[collapsible=icon]:p-0"
      >
        <!-- Pulse Green Dot Logo -->
        <div class="relative flex h-3 w-3">
          <span
            class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"
          ></span>
          <span class="relative inline-flex rounded-full h-3 w-3 bg-emerald-500"></span>
        </div>
        <span class="font-bold text-sm tracking-tight text-foreground whitespace-nowrap"
          >{{ siteConfig.name }}</span
        >
      </div>
    </SidebarHeader>

    <!-- Content: Nav Groups & Items -->
    <SidebarContent class="py-2">
      <SidebarGroup v-for="group in sidebarContent" :key="group.groupName">
        <SidebarGroupLabel v-if="group.groupName">{{ group.groupName }}</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem v-for="item in group.items" :key="item.title">
              <SidebarMenuButton
                as-child
                :is-active="route.path === item.href"
                :tooltip="item.title"
              >
                <RouterLink :to="item.href" class="flex items-center gap-3">
                  <component :is="item.icon" class="w-4 h-4 shrink-0" />
                  <span class="group-data-[collapsible=icon]:hidden text-sm">{{ item.title }}</span>
                </RouterLink>
              </SidebarMenuButton>
              <SidebarMenuBadge
                v-if="item.badge"
                class="group-data-[collapsible=icon]:hidden bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 text-[10px] font-bold px-1.5 py-0.5 rounded"
              >
                {{ item.badge }}
              </SidebarMenuBadge>
            </SidebarMenuItem>
          </SidebarMenu>
        </SidebarGroupContent>
      </SidebarGroup>
    </SidebarContent>

    <!-- Footer: Theme Toggle & Logout/User profile -->
    <SidebarFooter class="border-t border-border/50 p-4 flex flex-col gap-2">
      <!-- Theme switcher item -->
      <SidebarMenu>
        <SidebarMenuItem>
          <ThemeButton is-sidebar-item />
        </SidebarMenuItem>
      </SidebarMenu>

      <!-- User logout profile item -->
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton as-child variant="default" tooltip="Logout">
            <button
              @click="handleLogout"
              class="flex items-center gap-3 text-destructive hover:bg-destructive/10 hover:text-destructive w-full cursor-pointer"
            >
              <LogOut class="w-4 h-4 shrink-0" />
              <span class="group-data-[collapsible=icon]:hidden text-sm">Sign Out</span>
            </button>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarFooter>
  </Sidebar>
</template>
