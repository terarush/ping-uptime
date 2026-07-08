<script setup lang="ts">
/**
 * Navigation Sidebar Component
 * Connects with Vue Router to highight active items,
 * and calls the Auth composable logout actions on logout trigger.
 */
import { computed, onMounted } from 'vue'
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
  useSidebar,
} from '@/components/ui/sidebar'
import { sidebarContent } from '@/content/sidebar'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import ThemeButton from '@/components/theme-button.vue'
import { LogOut, Download } from '@lucide/vue'
import { useAuth } from '@/composables/useAuth'
import { useAppTitle } from '@/composables/useAppTitle'
import { useUpdateChecker } from '@/composables/useUpdateChecker'
import { siteConfig } from '@/content/config'

// Initialize router and auth dependencies
const { appTitle, fetchAppTitle } = useAppTitle()
const route = useRoute()
const router = useRouter()
const { currentUser, logout } = useAuth()
const { isMobile, setOpenMobile } = useSidebar()
const { isUpdateAvailable, latestRelease } = useUpdateChecker()

onMounted(fetchAppTitle)

const isAdmin = computed(() => currentUser.value?.role === 'admin')

const handleItemClick = () => {
  if (isMobile.value) {
    setOpenMobile(false)
  }
}

// Filter sidebar content based on permissions
const filteredSidebarContent = computed(() => {
  return sidebarContent.filter(group => {
    if (group.admin && !isAdmin.value) return false
    return true
  })
})

/**
 * Handles clearing sessions and redirecting standard users to the login index
 */
const handleLogout = () => {
  logout()
  router.push('/')
}
</script>


<template>
  <Sidebar collapsible="icon" variant="sidebar">
    <SidebarHeader class="h-16 shrink-0 border-b border-border/50 p-0 px-6 flex flex-row items-center justify-between group-data-[collapsible=icon]:justify-center group-data-[collapsible=icon]:px-0">
      <div class="flex items-center gap-2 overflow-hidden">
        <!-- Pulse Green Dot Logo -->
        <div class="relative flex h-3 w-3 shrink-0">
          <span
            class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"
          ></span>
          <span class="relative inline-flex rounded-full h-3 w-3 bg-emerald-500"></span>
        </div>
        <span class="font-bold text-sm tracking-tight text-foreground whitespace-nowrap group-data-[collapsible=icon]:hidden"
          >{{ appTitle }}</span
        >
      </div>
    </SidebarHeader>

    <!-- Content: Nav Groups & Items -->
    <SidebarContent class="py-2 overflow-y-auto [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none]">
      <SidebarGroup v-for="group in filteredSidebarContent" :key="group.groupName">
        <SidebarGroupLabel v-if="group.groupName">{{ group.groupName }}</SidebarGroupLabel>
        <SidebarGroupContent>
          <SidebarMenu>
            <SidebarMenuItem v-for="item in group.items" :key="item.title">
              <SidebarMenuButton
                as-child
                :is-active="route.path === item.href"
                :tooltip="item.title"
              >
                <RouterLink :to="item.href" @click="handleItemClick" class="flex items-center gap-3">
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
      <!-- Update available banner -->
      <SidebarMenu v-if="isUpdateAvailable">
        <SidebarMenuItem>
          <SidebarMenuButton as-child variant="default" tooltip="Update available">
            <a
              :href="latestRelease?.html_url || `https://github.com/${siteConfig.repoOwner}/${siteConfig.repoName}/releases/latest`"
              target="_blank"
              rel="noopener noreferrer"
              class="flex items-center gap-3 text-amber-600 dark:text-amber-400 hover:bg-amber-500/10 w-full cursor-pointer"
            >
              <Download class="w-4 h-4 shrink-0" />
              <span class="group-data-[collapsible=icon]:hidden text-sm font-semibold">
                {{ latestRelease?.tag_name }} available
              </span>
            </a>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
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
