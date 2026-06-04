<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { Mail, Lock, User, Eye, EyeOff, Loader2, AlertCircle, CheckCircle2, ShieldAlert } from '@lucide/vue';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card';
import ThemeButton from '@/components/theme-button.vue';
import ExtendedFetch from '@/lib/fetch';
import Cookies from 'js-cookie';
import gsap from 'gsap';
import { useAuth } from '@/composables/useAuth';
import { loginSchema, setupSchema } from '@/validations/auth';
import { Form, FormField, FormItem, FormLabel, FormControl, FormMessage } from '@/components/ui/form';
import { toTypedSchema } from '@vee-validate/zod';
import { siteConfig } from '@/content/config';

const router = useRouter();
const { isAuthenticated, setSession } = useAuth();

// UI & Logic States
const isSetupMode = ref(false); // If true, first-time setup is active
const isRegisterMode = ref(false);
const allowRegistration = ref(true);
const systemName = ref('Ping Uptime');
const checkingAuth = ref(true);
const loading = ref(false);
const error = ref('');
const success = ref('');
const showPassword = ref(false);

// Form Schema
const formSchema = computed(() => {
  return toTypedSchema(isSetupMode.value || isRegisterMode.value ? setupSchema : loginSchema);
});

// Trigger GSAP entry animations on load
const runEntryAnimations = () => {
  gsap.fromTo('.ambient-orb-1',
    { opacity: 0, scale: 0.6 },
    { opacity: 0.6, scale: 1, duration: 2, ease: 'power3.out' }
  );

  gsap.fromTo('.ambient-orb-2',
    { opacity: 0, scale: 0.6 },
    { opacity: 0.5, scale: 1, duration: 2.2, ease: 'power3.out', delay: 0.2 }
  );

  gsap.fromTo('.login-logo',
    { opacity: 0, y: -30 },
    { opacity: 1, y: 0, duration: 0.8, ease: 'back.out(1.5)', delay: 0.3 }
  );

  gsap.fromTo('.login-card',
    { opacity: 0, y: 40, scale: 0.95 },
    { opacity: 1, y: 0, scale: 1, duration: 0.9, ease: 'power4.out', delay: 0.5 }
  );

  gsap.fromTo('.theme-toggle-container',
    { opacity: 0, y: -10 },
    { opacity: 1, y: 0, duration: 0.6, ease: 'power2.out', delay: 0.7 }
  );

  if (isSetupMode.value) {
    gsap.fromTo('.setup-banner',
      { opacity: 0, y: -20, scale: 0.95 },
      { opacity: 1, y: 0, scale: 1, duration: 0.7, ease: 'back.out(1.5)', delay: 0.4 }
    );

    gsap.fromTo('.setup-banner-icon',
      { scale: 0.9 },
      { scale: 1.15, repeat: -1, yoyo: true, duration: 1.2, ease: 'sine.inOut', delay: 1.0 }
    );
  }
};

onMounted(async () => {
  let isSetup = false;

  // 1. Query Setup Status from the backend first
  try {
    const response = await ExtendedFetch.get('/auth/setup-status');
    isSetup = response.data?.data?.is_setup;
    isSetupMode.value = !isSetup; // Show setup if not yet setup
    systemName.value = response.data?.data?.system_name || 'Ping Uptime';
    allowRegistration.value = response.data?.data?.allow_registration ?? true;
  } catch (err) {
    console.error('Failed to query setup status:', err);
    // Default to login mode if API fails or is unreachable
    isSetupMode.value = false;
    isSetup = true; // assume setup is done so we check login
  }

  // 2. Clear stale cookies if setup is still required; otherwise, redirect if already logged in
  if (!isSetup) {
    Cookies.remove('accessToken');
    Cookies.remove('refreshToken');
  } else {
    // If the router guard already authenticated us, redirect to app directly
    if (isAuthenticated.value) {
      router.push(siteConfig.appPath);
      return;
    }
  }

  checkingAuth.value = false;
  // Delay slightly to allow the DOM to render before animating
  setTimeout(runEntryAnimations, 50);
});

// Submit Authentication / Setup Request
const onSubmit = async (values: any) => {
  error.value = '';
  success.value = '';
  loading.value = true;

  try {
    if (isSetupMode.value) {
      await ExtendedFetch.post('/auth/setup', {
        name: values.name,
        email: values.email,
        password: values.password,
      });

      success.value = 'Administrator registered successfully! Logging in...';

      // Auto Login
      const loginResponse = await ExtendedFetch.post('/auth/login', {
        email: values.email,
        password: values.password,
      });

      const accessToken = loginResponse.data?.data?.accessToken;
      const refreshToken = loginResponse.data?.data?.refreshToken;
      const user = loginResponse.data?.data?.user;
      if (accessToken && refreshToken && user) {
        setSession(accessToken, refreshToken, user);
        setTimeout(() => {
          router.push(siteConfig.appPath);
        }, 800);
      } else {
        error.value = 'Setup succeeded, but auto-login failed. Please refresh and log in.';
        isSetupMode.value = false;
      }
    } else if (isRegisterMode.value) {
      await ExtendedFetch.post('/auth/register', {
        name: values.name,
        email: values.email,
        password: values.password,
      });

      success.value = 'Account registered successfully! Logging in...';

      // Auto Login
      const loginResponse = await ExtendedFetch.post('/auth/login', {
        email: values.email,
        password: values.password,
      });

      const accessToken = loginResponse.data?.data?.accessToken;
      const refreshToken = loginResponse.data?.data?.refreshToken;
      const user = loginResponse.data?.data?.user;
      if (accessToken && refreshToken && user) {
        setSession(accessToken, refreshToken, user);
        setTimeout(() => {
          router.push(siteConfig.appPath);
        }, 800);
      } else {
        error.value = 'Registration succeeded, but auto-login failed. Please sign in.';
        isRegisterMode.value = false;
      }
    } else {
      const response = await ExtendedFetch.post('/auth/login', {
        email: values.email,
        password: values.password,
      });

      const accessToken = response.data?.data?.accessToken;
      const refreshToken = response.data?.data?.refreshToken;
      const user = response.data?.data?.user;
      if (accessToken && refreshToken && user) {
        success.value = 'Login successful! Entering dashboard...';
        setSession(accessToken, refreshToken, user);
        setTimeout(() => {
          router.push(siteConfig.appPath);
        }, 800);
      } else {
        error.value = 'Failed to retrieve login session.';
        triggerShake();
      }
    }
  } catch (err: any) {
    console.error('Authentication Error:', err);
    error.value = err.response?.data?.error || err.message || 'An error occurred. Please try again.';
    triggerShake();
  } finally {
    loading.value = false;
  }
};

const onInvalidSubmit = () => {
  triggerShake();
};

// Error shake feedback animation
const triggerShake = () => {
  const tl = gsap.timeline();
  tl.to('.login-card', { x: -6, duration: 0.05 })
    .to('.login-card', { x: 6, duration: 0.05 })
    .to('.login-card', { x: -6, duration: 0.05 })
    .to('.login-card', { x: 6, duration: 0.05 })
    .to('.login-card', { x: -3, duration: 0.05 })
    .to('.login-card', { x: 3, duration: 0.05 })
    .to('.login-card', { x: 0, duration: 0.05 });
};
</script>

<template>
  <div class="relative min-h-screen flex flex-col justify-center items-center px-4 overflow-hidden bg-background">
    <!-- Beautiful Ambient Background Orbs -->
    <div class="ambient-orb-1 absolute top-[-10%] left-[-10%] w-[60%] h-[60%] rounded-full bg-emerald-500/10 dark:bg-emerald-500/5 blur-[120px] pointer-events-none"></div>
    <div class="ambient-orb-2 absolute bottom-[-10%] right-[-10%] w-[60%] h-[60%] rounded-full bg-emerald-600/10 dark:bg-emerald-600/5 blur-[120px] pointer-events-none"></div>

    <!-- Loading blocker while checking auth cookie & setup status -->
    <div v-if="checkingAuth" class="flex flex-col items-center justify-center space-y-4">
      <Loader2 class="h-8 w-8 text-primary animate-spin" />
      <span class="text-xs text-muted-foreground font-semibold">Checking system status...</span>
    </div>

    <div v-else class="w-full max-w-105 flex flex-col items-center space-y-6 z-10">
      <!-- Theme Toggle in top corner -->
      <div class="theme-toggle-container absolute top-6 right-6">
        <ThemeButton variant="rounded" />
      </div>

      <!-- App Logo Section -->
      <div class="login-logo flex flex-col items-center text-center">
        <div class="flex items-center gap-2 mb-2">
          <div class="relative flex h-3.5 w-3.5">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-3.5 w-3.5 bg-emerald-500"></span>
          </div>
          <h1 class="text-2xl font-black tracking-tight text-foreground select-none">{{ systemName }}</h1>
        </div>
        <p class="text-xs text-muted-foreground max-w-70">
          Monitor your services with custom real-time alerts and uptime statistics.
        </p>
      </div>

      <!-- Setup Notice Banner -->
      <div v-if="isSetupMode" class="setup-banner w-full p-4 bg-amber-500/10 border border-amber-500/20 text-amber-600 dark:text-amber-400 rounded-2xl flex gap-3 text-xs leading-relaxed opacity-0">
        <ShieldAlert class="setup-banner-icon h-5 w-5 shrink-0 mt-0.5 text-amber-500" />
        <div>
          <span class="font-bold block">First-Time Setup Required</span>
          There are no accounts registered on this instance. Create the initial administrator account below.
        </div>
      </div>

      <!-- Main Auth Card -->
      <Card class="login-card w-full bg-card/60 dark:bg-card/40 backdrop-blur-xl border border-border/60 dark:border-border/20 shadow-2xl rounded-2xl overflow-hidden p-6 gap-0">
        <!-- Form Header -->
        <CardHeader class="p-0 pb-6 text-center">
          <CardTitle class="text-lg font-bold tracking-tight text-foreground">
            {{ isSetupMode ? 'Setup Admin Account' : (isRegisterMode ? 'Create Account' : 'Welcome back') }}
          </CardTitle>
          <CardDescription class="text-xs text-muted-foreground mt-1">
            {{ isSetupMode ? 'Register the main account for this application' : (isRegisterMode ? 'Sign up to register your dashboard account' : 'Sign in to access your self-hosted panel') }}
          </CardDescription>
        </CardHeader>

        <CardContent class="p-0 space-y-4">
          <!-- Alert Messages -->
          <div v-if="error" class="p-3 text-xs bg-destructive/10 border border-destructive/20 text-destructive rounded-lg flex items-start gap-2 animate-[fadeIn_0.2s_ease-out]">
            <AlertCircle class="h-4 w-4 shrink-0 mt-0.5" />
            <span>{{ error }}</span>
          </div>

          <div v-if="success" class="p-3 text-xs bg-emerald-500/10 border border-emerald-500/20 text-emerald-600 dark:text-emerald-400 rounded-lg flex items-start gap-2 animate-[fadeIn_0.2s_ease-out]">
            <CheckCircle2 class="h-4 w-4 shrink-0 mt-0.5" />
            <span>{{ success }}</span>
          </div>

          <!-- Auth Form -->
          <Form :validation-schema="formSchema" @submit="onSubmit" @invalid-submit="onInvalidSubmit" class="space-y-4">
            <!-- Full Name (Only visible in setup/register mode) -->
            <FormField v-slot="{ componentField }" v-if="isSetupMode || isRegisterMode" name="name">
              <FormItem class="space-y-1.5 animate-[fadeIn_0.25s_ease-out]">
                <FormLabel class="text-xs font-bold text-foreground/80">Full Name</FormLabel>
                <FormControl>
                  <div class="relative">
                    <User class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                    <Input
                      v-bind="componentField"
                      type="text"
                      placeholder="John Doe"
                      class="pl-9 h-10 rounded-lg"
                    />
                  </div>
                </FormControl>
                <FormMessage class="text-[10px] text-destructive" />
              </FormItem>
            </FormField>

            <!-- Email Address -->
            <FormField name="email" v-slot="{ componentField }">
              <FormItem class="space-y-1.5">
                <FormLabel class="text-xs font-bold text-foreground/80">Email Address</FormLabel>
                <FormControl>
                  <div class="relative">
                    <Mail class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                    <Input
                      v-bind="componentField"
                      type="email"
                      placeholder="admin@example.com"
                      class="pl-9 h-10 rounded-lg"
                    />
                  </div>
                </FormControl>
                <FormMessage class="text-[10px] text-destructive" />
              </FormItem>
            </FormField>

            <!-- Password -->
            <FormField name="password" v-slot="{ componentField }">
              <FormItem class="space-y-1.5">
                <FormLabel class="text-xs font-bold text-foreground/80">Password</FormLabel>
                <FormControl>
                  <div class="relative">
                    <Lock class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                    <Input
                      v-bind="componentField"
                      :type="showPassword ? 'text' : 'password'"
                      placeholder="••••••••"
                      class="pl-9 pr-10 h-10 rounded-lg"
                    />
                    <button
                      type="button"
                      @click="showPassword = !showPassword"
                      class="absolute right-3 top-2.5 text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
                    >
                      <Eye v-if="!showPassword" class="h-4 w-4" />
                      <EyeOff v-else class="h-4 w-4" />
                    </button>
                  </div>
                </FormControl>
                <FormMessage class="text-[10px] text-destructive" />
              </FormItem>
            </FormField>

            <!-- Submit Button -->
            <Button
              type="submit"
              :disabled="loading"
              class="w-full h-10 rounded-lg font-bold shadow-lg hover:shadow-xl transition-all duration-300 hover:scale-[1.01] bg-primary hover:bg-primary/90 text-primary-foreground flex items-center justify-center gap-2 cursor-pointer mt-2"
            >
              <Loader2 v-if="loading" class="h-4 w-4 animate-spin" />
              <span v-else>{{ isSetupMode ? 'Create Admin & Get Started' : (isRegisterMode ? 'Sign Up' : 'Sign In to Dashboard') }}</span>
            </Button>
          </Form>

          <!-- Toggle Register / Login links -->
          <div v-if="!isSetupMode && allowRegistration" class="text-center text-xs mt-4 pt-2 border-t border-border/40">
            <a
              href="#"
              @click.prevent="isRegisterMode = !isRegisterMode; error = ''; success = '';"
              class="text-primary hover:underline font-semibold"
            >
              {{ isRegisterMode ? 'Already have an account? Sign in' : "Don't have an account? Register" }}
            </a>
          </div>
        </CardContent>

        <CardFooter class="p-0 pt-6 text-center justify-center">
          <p class="text-[10px] text-muted-foreground">
            Self-hosted service status dashboard. Built with Vue 3 & Go.
          </p>
        </CardFooter>
      </Card>
    </div>
  </div>
</template>
