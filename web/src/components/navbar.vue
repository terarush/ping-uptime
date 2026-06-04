<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue';
import { useRoute, useRouter, RouterLink } from 'vue-router';
import {
  Menu as MenuIcon,
  Search as SearchIcon,
  X as XIcon,
  User as UserIcon,
  LogOut as LogOutIcon,
  LayoutDashboard,
  Crown,
} from '@lucide/vue';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import ThemeToggle from './theme-button.vue';
import { siteConfig, navConfig } from '@/content/config';
import gsap from 'gsap';

// Mock authentication context for Vue adaptation
const isAuthenticated = ref(true); // default to true for demonstration
const user = ref({
  name: 'Administrator',
  email: 'admin@ping-uptime.com',
  role: 'admin',
  avatar: ''
});

const isScrolled = ref(false);
const isSearchOpen = ref(false);
const isMobileMenuOpen = ref(false);

const route = useRoute();
const router = useRouter();

const headerRef = ref<HTMLElement | null>(null);
const mobileMenuRef = ref<HTMLDivElement | null>(null);
const mobileBackdropRef = ref<HTMLDivElement | null>(null);

const getInitials = (name: string) => {
  return name
    .split(' ')
    .map((n) => n[0])
    .join('')
    .toUpperCase()
    .slice(0, 2);
};

const getAvatarUrl = (avatar: string) => {
  return avatar || 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?q=80&w=100&auto=format&fit=crop';
};

const handleNavClick = (href: string) => {
  if (href.startsWith('#') && route.path === '/') {
    const target = document.querySelector(href);
    if (target) {
      target.scrollIntoView({ behavior: 'smooth' });
    }
  } else if (href.startsWith('#') && route.path !== '/') {
    router.push('/' + href);
  }
};

const logout = async () => {
  isAuthenticated.value = false;
};

const handleLogout = async () => {
  await logout();
  router.push('/');
};

// Handle route navigation to close mobile menu
watch(() => route.path, () => {
  isMobileMenuOpen.value = false;
});

// Watch mobile menu open state for GSAP animation
watch(isMobileMenuOpen, async (isOpen) => {
  if (isOpen) {
    await nextTick();
    if (mobileBackdropRef.value) {
      gsap.fromTo(
        mobileBackdropRef.value,
        { opacity: 0 },
        { opacity: 1, duration: 0.3, ease: 'power2.out' }
      );
    }

    if (mobileMenuRef.value) {
      gsap.fromTo(
        mobileMenuRef.value,
        { y: 100, opacity: 0 },
        { y: 0, opacity: 1, duration: 0.5, ease: 'power3.out' }
      );

      const menuItems = mobileMenuRef.value.querySelectorAll('.menu-item');
      gsap.fromTo(
        menuItems,
        { y: 20, opacity: 0 },
        {
          y: 0,
          opacity: 1,
          duration: 0.4,
          stagger: 0.05,
          ease: 'power2.out',
          delay: 0.2,
        }
      );
    }
  } else {
    if (mobileMenuRef.value) {
      gsap.to(mobileMenuRef.value, {
        y: 100,
        opacity: 0,
        duration: 0.3,
        ease: 'power2.in',
      });
    }
    if (mobileBackdropRef.value) {
      gsap.to(mobileBackdropRef.value, {
        opacity: 0,
        duration: 0.2,
        ease: 'power2.in',
      });
    }
  }
});

const handleScroll = () => {
  isScrolled.value = window.scrollY > 20;
};

onMounted(() => {
  window.addEventListener('scroll', handleScroll);

  // GSAP animation for header on mount
  if (headerRef.value) {
    gsap.fromTo(
      headerRef.value,
      { y: -100, opacity: 0 },
      { y: 0, opacity: 1, duration: 0.8, ease: 'power3.out' }
    );
  }
});

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll);
});
</script>

<template>
  <header
    ref="headerRef"
    class="fixed top-4 left-0 right-0 z-50 px-4 opacity-0"
  >
    <nav
      :class="[
        'relative backdrop-blur-xl border border-border shadow-md shadow-gray-900/10 transition-all duration-500 rounded-xl mx-auto py-3 max-w-7xl bg-background/80'
      ]"
    >
      <div class="relative px-6 flex items-center justify-between">
        <div class="flex items-center space-x-8">
          <RouterLink
            to="/"
            class="flex items-center space-x-3 group"
          >
            <div class="relative">
              <img
                src="/favicon.ico"
                :alt="siteConfig.name"
                class="h-8 w-8 object-contain transition-transform duration-300 group-hover:scale-110 rounded-md"
              />
            </div>
            <span class="text-xl font-bold text-gray-900 dark:text-white">
              {{ siteConfig.name }}
            </span>
          </RouterLink>
          <div class="hidden lg:flex items-center space-x-1">
            <template v-for="item in navConfig.mainNav" :key="item.href">
              <!-- Use RouterLink component for internal routes -->
              <RouterLink
                v-if="item.href.startsWith('/')"
                :to="item.href"
                class="px-4 py-2 rounded-xl text-gray-700 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200 font-medium"
              >
                {{ item.label }}
              </RouterLink>
              
              <!-- Use anchor tag for hash links -->
              <a
                v-else
                :href="item.href"
                @click.prevent="handleNavClick(item.href)"
                class="px-4 py-2 rounded-xl text-gray-700 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200 font-medium"
              >
                {{ item.label }}
              </a>
            </template>
          </div>
        </div>
        
        <div class="flex items-center space-x-3">
          <Dialog :open="isSearchOpen" @update:open="isSearchOpen = $event">
            <DialogTrigger as-child>
              <Button
                variant="ghost"
                size="sm"
                class="hidden lg:flex items-center space-x-2 px-3 py-2 rounded-xl bg-gray-100/50 dark:bg-gray-800/50 hover:bg-gray-200/50 dark:hover:bg-gray-700/50 transition-all duration-200 min-w-[200px] justify-start"
              >
                <SearchIcon class="h-4 w-4 text-gray-500" />
                <span class="text-gray-500 text-sm">Search...</span>
              </Button>
            </DialogTrigger>
            <DialogContent class="max-w-2xl">
              <DialogHeader>
                <DialogTitle>Search</DialogTitle>
              </DialogHeader>
              <div class="flex items-center space-x-2 border-b pb-4">
                <SearchIcon class="h-5 w-5 text-gray-400" />
                <Input
                  placeholder="Search products, docs, and more..."
                  class="border-0 focus-visible:ring-0 text-lg"
                  auto-focus
                />
              </div>
              <div class="py-4">
                <p class="text-sm text-gray-500">
                  Start typing to search...
                </p>
              </div>
            </DialogContent>
          </Dialog>

          <Button
            variant="ghost"
            size="sm"
            class="lg:hidden p-2 rounded-xl hover:bg-gray-100/50 dark:hover:bg-gray-800/50"
          >
            <SearchIcon class="h-5 w-5" />
          </Button>

          <ThemeToggle variant="rounded" />

          <!-- User Menu - Show when authenticated -->
          <DropdownMenu v-if="isAuthenticated && user">
            <DropdownMenuTrigger as-child>
              <Button
                variant="ghost"
                class="hidden sm:flex items-center space-x-2 px-3 py-2 rounded-xl hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200"
              >
                <Avatar class="h-8 w-8 border-2 border-primary/20">
                  <AvatarImage
                    v-if="user.avatar"
                    :src="getAvatarUrl(user.avatar)"
                    :alt="user.name"
                  />
                  <AvatarFallback class="bg-primary/10 text-primary text-xs font-semibold">
                    {{ getInitials(user.name) }}
                  </AvatarFallback>
                </Avatar>
                <div class="flex flex-col items-start">
                  <div class="flex items-center gap-1.5">
                    <span class="text-sm font-medium">{{ user.name }}</span>
                    <Crown v-if="user.role === 'admin'" class="h-3 w-3 text-yellow-500" />
                  </div>
                </div>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" class="w-56">
              <DropdownMenuLabel>
                <div class="flex flex-col space-y-1">
                  <div class="flex items-center gap-1.5">
                    <p class="text-sm font-medium leading-none">{{ user.name }}</p>
                    <Crown v-if="user.role === 'admin'" class="h-3.5 w-3.5 text-yellow-500" />
                  </div>
                  <p class="text-xs leading-none text-muted-foreground">{{ user.email }}</p>
                </div>
              </DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem @click="router.push('/dashboard')">
                <LayoutDashboard class="mr-2 h-4 w-4" />
                Dashboard
              </DropdownMenuItem>
              <DropdownMenuItem @click="router.push('/dashboard/profile')">
                <UserIcon class="mr-2 h-4 w-4" />
                Profile
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                @click="handleLogout"
                class="text-destructive focus:text-destructive"
              >
                <LogOutIcon class="mr-2 h-4 w-4" />
                Logout
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          <Button
            v-else
            class="hidden sm:flex bg-gray-900 dark:bg-gray-100 text-white dark:text-gray-900 rounded-xl px-6 py-2 font-semibold shadow-lg hover:shadow-xl transition-all duration-300 hover:scale-105"
            as-child
          >
            <RouterLink to="/login">Get Started</RouterLink>
          </Button>

          <Button
            variant="ghost"
            class="lg:hidden p-2 rounded-md hover:bg-gray-100/50 dark:hover:bg-gray-800/50"
            @click="isMobileMenuOpen = !isMobileMenuOpen"
          >
            <XIcon v-if="isMobileMenuOpen" class="h-6 w-6" />
            <MenuIcon v-else class="h-6 w-6" />
          </Button>
        </div>
      </div>
    </nav>
  </header>

  <!-- Mobile Menu Backdrop & Drawer -->
  <div v-if="isMobileMenuOpen" class="fixed inset-0 z-40 lg:hidden">
    <div
      ref="mobileBackdropRef"
      class="fixed inset-0 bg-black/20 backdrop-blur-sm opacity-0"
      @click="isMobileMenuOpen = false"
    />
    <div class="fixed bottom-0 left-0 right-0 p-4 pointer-events-none">
      <div
        ref="mobileMenuRef"
        class="relative bg-background/95 backdrop-blur-xl border border-border shadow-2xl rounded-2xl p-6 pointer-events-auto translate-y-full opacity-0"
      >
        <div class="space-y-4">
          <!-- User Info in Mobile Menu -->
          <div v-if="isAuthenticated && user" class="menu-item flex items-center space-x-3 p-3 rounded-xl bg-gray-100/50 dark:bg-gray-800/50">
            <Avatar class="h-10 w-10 border-2 border-primary/20">
              <AvatarImage
                v-if="user.avatar"
                :src="getAvatarUrl(user.avatar)"
                :alt="user.name"
              />
              <AvatarFallback class="bg-primary/10 text-primary font-semibold">
                {{ getInitials(user.name) }}
              </AvatarFallback>
            </Avatar>
            <div class="flex-1">
              <div class="flex items-center gap-1.5">
                <p class="text-sm font-medium">{{ user.name }}</p>
                <Crown v-if="user.role === 'admin'" class="h-3 w-3 text-yellow-500" />
              </div>
              <p class="text-xs text-muted-foreground">{{ user.email }}</p>
            </div>
          </div>

          <div class="space-y-2">
            <template v-for="item in navConfig.mainNav" :key="item.href">
              <!-- Use RouterLink for internal routes -->
              <RouterLink
                v-if="item.href.startsWith('/')"
                :to="item.href"
                class="menu-item flex items-center space-x-3 p-3 rounded-xl hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200"
                @click="isMobileMenuOpen = false"
              >
                <span class="font-medium text-gray-900 dark:text-white">
                  {{ item.label }}
                </span>
              </RouterLink>
              
              <!-- Use anchor tag for hash links -->
              <a
                v-else
                :href="item.href"
                class="menu-item flex items-center space-x-3 p-3 rounded-xl hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200"
                @click.prevent="isMobileMenuOpen = false; handleNavClick(item.href)"
              >
                <span class="font-medium text-gray-900 dark:text-white">
                  {{ item.label }}
                </span>
              </a>
            </template>

            <!-- Additional menu items when authenticated -->
            <template v-if="isAuthenticated">
              <RouterLink
                to="/dashboard"
                class="menu-item flex items-center space-x-3 p-3 rounded-xl hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200"
                @click="isMobileMenuOpen = false"
              >
                <LayoutDashboard class="h-5 w-5 text-gray-500" />
                <span class="font-medium text-gray-900 dark:text-white">Dashboard</span>
              </RouterLink>
              <RouterLink
                to="/dashboard/profile"
                class="menu-item flex items-center space-x-3 p-3 rounded-xl hover:bg-gray-100/50 dark:hover:bg-gray-800/50 transition-all duration-200"
                @click="isMobileMenuOpen = false"
              >
                <UserIcon class="h-5 w-5 text-gray-500" />
                <span class="font-medium text-gray-900 dark:text-white">Profile</span>
              </RouterLink>
            </template>
          </div>

          <div class="pt-2 menu-item">
            <Button
              v-if="isAuthenticated"
              class="w-full bg-destructive text-destructive-foreground rounded-xl py-3 font-semibold"
              @click="handleLogout(); isMobileMenuOpen = false;"
            >
              <LogOutIcon class="mr-2 h-5 w-5" />
              Logout
            </Button>
            <Button
              v-else
              class="w-full bg-gray-900 dark:bg-gray-100 text-white dark:text-gray-900 rounded-xl py-3 font-semibold"
              as-child
            >
              <RouterLink to="/login">Get Started</RouterLink>
            </Button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
