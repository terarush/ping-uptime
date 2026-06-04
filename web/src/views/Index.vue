<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useRouter } from 'vue-router';
import { Mail, Lock, User, Eye, EyeOff, Loader2, AlertCircle, CheckCircle2 } from '@lucide/vue';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card';
import ThemeButton from '@/components/theme-button.vue';
import ExtendedFetch from '@/lib/fetch';
import Cookies from 'js-cookie';
import gsap from 'gsap';

const router = useRouter();

// Form States
const isSignUp = ref(false);
const name = ref('');
const email = ref('');
const password = ref('');
const showPassword = ref(false);

// UI States
const loading = ref(false);
const error = ref('');
const success = ref('');
const checkingAuth = ref(true);

// Clear messages on tab toggle
watch(isSignUp, () => {
  error.value = '';
  success.value = '';
});

// Trigger GSAP entry animations on load
onMounted(() => {
  // Check if user is already authenticated via cookie
  const token = Cookies.get('accessToken');
  if (token) {
    router.push('/app');
    return;
  }
  checkingAuth.value = false;

  // Run GSAP entry animations
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
});

// Submit Authentication Request
const handleSubmit = async () => {
  error.value = '';
  success.value = '';
  loading.value = true;

  try {
    if (isSignUp.value) {
      // Sign Up flow
      if (!name.value.trim()) {
        error.value = 'Please enter your full name.';
        loading.value = false;
        triggerShake();
        return;
      }
      if (password.value.length < 6) {
        error.value = 'Password must be at least 6 characters.';
        loading.value = false;
        triggerShake();
        return;
      }

      await ExtendedFetch.post('/auth/register', {
        name: name.value,
        email: email.value,
        password: password.value,
      });

      // Automatically sign in the user on registration success
      success.value = 'Account created successfully! Logging in...';

      const loginResponse = await ExtendedFetch.post('/auth/login', {
        email: email.value,
        password: password.value,
      });

      const token = loginResponse.data?.data?.token;
      if (token) {
        Cookies.set('accessToken', token, { expires: 7 });
        setTimeout(() => {
          router.push('/app');
        }, 800);
      } else {
        error.value = 'Registration succeeded, but auto-login failed. Please sign in manually.';
        isSignUp.value = false;
      }
    } else {
      // Sign In flow
      const response = await ExtendedFetch.post('/auth/login', {
        email: email.value,
        password: password.value,
      });

      const token = response.data?.data?.token;
      if (token) {
        success.value = 'Login successful! Entering dashboard...';
        Cookies.set('accessToken', token, { expires: 7 });
        setTimeout(() => {
          router.push('/app');
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

    <!-- Loading blocker while checking auth cookie -->
    <div v-if="checkingAuth" class="flex flex-col items-center justify-center space-y-4">
      <Loader2 class="h-8 w-8 text-primary animate-spin" />
      <span class="text-xs text-muted-foreground font-semibold">Resuming session...</span>
    </div>

    <div v-else class="w-full max-w-[420px] flex flex-col items-center space-y-6 z-10">
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
          <h1 class="text-2xl font-black tracking-tight text-foreground select-none">Ping Uptime</h1>
        </div>
        <p class="text-xs text-muted-foreground max-w-[280px]">
          Monitor your services with custom real-time alerts and uptime statistics.
        </p>
      </div>

      <!-- Main Login Card -->
      <Card class="login-card w-full bg-card/60 dark:bg-card/40 backdrop-blur-xl border border-border/60 dark:border-border/20 shadow-2xl rounded-2xl overflow-hidden p-6 gap-0">
        <!-- Form Header -->
        <CardHeader class="p-0 pb-6 text-center">
          <CardTitle class="text-lg font-bold tracking-tight text-foreground">
            {{ isSignUp ? 'Create your account' : 'Welcome back' }}
          </CardTitle>
          <CardDescription class="text-xs text-muted-foreground mt-1">
            {{ isSignUp ? 'Enter your details below to get started' : 'Sign in to access your self-hosted panel' }}
          </CardDescription>
        </CardHeader>

        <CardContent class="p-0 space-y-4">
          <!-- Animated Tab Toggle (Sign In / Sign Up) -->
          <div class="grid w-full grid-cols-2 p-1 bg-muted/65 dark:bg-muted/40 rounded-xl mb-2">
            <button
              type="button"
              @click="isSignUp = false"
              :class="[
                'py-2 text-xs font-semibold rounded-lg transition-all duration-300 select-none cursor-pointer',
                !isSignUp ? 'bg-card text-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground'
              ]"
            >
              Sign In
            </button>
            <button
              type="button"
              @click="isSignUp = true"
              :class="[
                'py-2 text-xs font-semibold rounded-lg transition-all duration-300 select-none cursor-pointer',
                isSignUp ? 'bg-card text-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground'
              ]"
            >
              Sign Up
            </button>
          </div>

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
          <form @submit.prevent="handleSubmit" class="space-y-4">
            <!-- Full Name (Only visible when signing up) -->
            <Transition name="expand">
              <div v-if="isSignUp" class="space-y-1.5 overflow-hidden">
                <label class="text-xs font-bold text-foreground/80">Full Name</label>
                <div class="relative">
                  <User class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                  <Input
                    v-model="name"
                    type="text"
                    placeholder="John Doe"
                    class="pl-9 h-10 rounded-lg"
                    required
                  />
                </div>
              </div>
            </Transition>

            <!-- Email Address -->
            <div class="space-y-1.5">
              <label class="text-xs font-bold text-foreground/80">Email Address</label>
              <div class="relative">
                <Mail class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  v-model="email"
                  type="email"
                  placeholder="admin@ping-uptime.com"
                  class="pl-9 h-10 rounded-lg"
                  required
                />
              </div>
            </div>

            <!-- Password -->
            <div class="space-y-1.5">
              <div class="flex items-center justify-between">
                <label class="text-xs font-bold text-foreground/80">Password</label>
              </div>
              <div class="relative">
                <Lock class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  v-model="password"
                  :type="showPassword ? 'text' : 'password'"
                  placeholder="••••••••"
                  class="pl-9 pr-10 h-10 rounded-lg"
                  required
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
            </div>

            <!-- Submit Button -->
            <Button
              type="submit"
              :disabled="loading"
              class="w-full h-10 rounded-lg font-bold shadow-lg hover:shadow-xl transition-all duration-300 hover:scale-[1.01] bg-primary hover:bg-primary/90 text-primary-foreground flex items-center justify-center gap-2 cursor-pointer mt-2"
            >
              <Loader2 v-if="loading" class="h-4 w-4 animate-spin" />
              <span v-else>{{ isSignUp ? 'Register & Sign In' : 'Sign In to Dashboard' }}</span>
            </Button>
          </form>
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

<style scoped>
.expand-enter-active,
.expand-leave-active {
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  max-height: 80px;
  opacity: 1;
}
.expand-enter-from,
.expand-leave-to {
  max-height: 0px;
  opacity: 0;
  margin-top: 0 !important;
  margin-bottom: 0 !important;
  padding-top: 0 !important;
  padding-bottom: 0 !important;
}
</style>
