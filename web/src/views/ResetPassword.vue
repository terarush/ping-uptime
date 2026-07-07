<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { Lock, Eye, EyeOff, Loader2, AlertCircle, CheckCircle2 } from '@lucide/vue';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card';
import ThemeButton from '@/components/theme-button.vue';
import ExtendedFetch from '@/lib/fetch';
import gsap from 'gsap';
import { Form, FormField, FormItem, FormLabel, FormControl, FormMessage } from '@/components/ui/form';
import { toTypedSchema } from '@vee-validate/zod';
import { z } from 'zod';

const route = useRoute();
const router = useRouter();

const token = (route.query.token as string) || '';

const resetSchema = z.object({
  password: z.string({ required_error: 'Password is required' }).min(6, 'Password must be at least 6 characters'),
  confirmPassword: z.string({ required_error: 'Please confirm your password' }).min(1, 'Please confirm your password'),
}).refine(data => data.password === data.confirmPassword, {
  message: 'Passwords do not match',
  path: ['confirmPassword'],
});
const formSchema = toTypedSchema(resetSchema);

const loading = ref(false);
const success = ref(false);
const error = ref('');
const showPassword = ref(false);
const showConfirm = ref(false);

onMounted(() => {
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
  if (!token) {
    error.value = 'Invalid or missing reset token.';
  }
});

const onSubmit = async (values: { password: string }) => {
  if (!token) { error.value = 'Invalid reset token.'; return; }
  error.value = '';
  loading.value = true;
  try {
    await ExtendedFetch.post('/auth/reset-password', { token, password: values.password });
    success.value = true;
    setTimeout(() => router.push('/'), 2500);
  } catch (err: any) {
    error.value = err.response?.data?.error || err.message || 'Reset failed. Token may be invalid or expired.';
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <div class="relative min-h-screen flex flex-col justify-center items-center px-4 overflow-hidden bg-background">
    <div class="ambient-orb-1 absolute top-[-10%] left-[-10%] w-[60%] h-[60%] rounded-full bg-emerald-500/10 dark:bg-emerald-500/5 blur-[120px] pointer-events-none"></div>
    <div class="ambient-orb-2 absolute bottom-[-10%] right-[-10%] w-[60%] h-[60%] rounded-full bg-emerald-600/10 dark:bg-emerald-600/5 blur-[120px] pointer-events-none"></div>

    <div class="w-full max-w-105 flex flex-col items-center space-y-6 z-10">
      <div class="theme-toggle-container absolute top-6 right-6">
        <ThemeButton variant="rounded" />
      </div>

      <div class="login-logo flex flex-col items-center text-center">
        <h1 class="text-2xl font-black tracking-tight text-foreground select-none">Reset Password</h1>
        <p class="text-xs text-muted-foreground max-w-70 mt-1">Choose a new password for your account.</p>
      </div>

      <Card class="login-card w-full bg-card/60 dark:bg-card/40 backdrop-blur-xl border border-border/60 dark:border-border/20 shadow-2xl rounded-2xl overflow-hidden p-6 gap-0">
        <CardHeader class="p-0 pb-6 text-center">
          <CardTitle class="text-lg font-bold tracking-tight text-foreground">Enter new password</CardTitle>
          <CardDescription class="text-xs text-muted-foreground mt-1">
            Must be at least 6 characters.
          </CardDescription>
        </CardHeader>

        <CardContent class="p-0 space-y-4">
          <div v-if="error" class="p-3 text-xs bg-destructive/10 border border-destructive/20 text-destructive rounded-lg flex items-start gap-2">
            <AlertCircle class="h-4 w-4 shrink-0 mt-0.5" />
            <span>{{ error }}</span>
          </div>

          <div v-if="success" class="p-3 text-xs bg-emerald-500/10 border border-emerald-500/20 text-emerald-600 dark:text-emerald-400 rounded-lg flex items-start gap-2">
            <CheckCircle2 class="h-4 w-4 shrink-0 mt-0.5" />
            <span>Password reset successfully! Redirecting to sign in...</span>
          </div>

          <Form v-if="!success" :validation-schema="formSchema" @submit="onSubmit" class="space-y-4">
            <FormField name="password" v-slot="{ componentField }">
              <FormItem class="space-y-1.5">
                <FormLabel class="text-xs font-bold text-foreground/80">New Password</FormLabel>
                <FormControl>
                  <div class="relative">
                    <Lock class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                    <Input v-bind="componentField" :type="showPassword ? 'text' : 'password'" placeholder="••••••••" class="pl-9 pr-10 h-10 rounded-lg" />
                    <button type="button" @click="showPassword = !showPassword" class="absolute right-3 top-2.5 text-muted-foreground hover:text-foreground transition-colors cursor-pointer">
                      <Eye v-if="!showPassword" class="h-4 w-4" />
                      <EyeOff v-else class="h-4 w-4" />
                    </button>
                  </div>
                </FormControl>
                <FormMessage class="text-[10px] text-destructive" />
              </FormItem>
            </FormField>

            <FormField name="confirmPassword" v-slot="{ componentField }">
              <FormItem class="space-y-1.5">
                <FormLabel class="text-xs font-bold text-foreground/80">Confirm Password</FormLabel>
                <FormControl>
                  <div class="relative">
                    <Lock class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                    <Input v-bind="componentField" :type="showConfirm ? 'text' : 'password'" placeholder="••••••••" class="pl-9 pr-10 h-10 rounded-lg" />
                    <button type="button" @click="showConfirm = !showConfirm" class="absolute right-3 top-2.5 text-muted-foreground hover:text-foreground transition-colors cursor-pointer">
                      <Eye v-if="!showConfirm" class="h-4 w-4" />
                      <EyeOff v-else class="h-4 w-4" />
                    </button>
                  </div>
                </FormControl>
                <FormMessage class="text-[10px] text-destructive" />
              </FormItem>
            </FormField>

            <Button type="submit" :disabled="loading"
              class="w-full h-10 rounded-lg font-bold shadow-lg hover:shadow-xl transition-all duration-300 hover:scale-[1.01] bg-primary hover:bg-primary/90 text-primary-foreground flex items-center justify-center gap-2 cursor-pointer mt-2">
              <Loader2 v-if="loading" class="h-4 w-4 animate-spin" />
              <span v-else>Reset Password</span>
            </Button>
          </Form>
        </CardContent>

        <CardFooter class="p-0 pt-6 text-center justify-center">
          <p class="text-[10px] text-muted-foreground">Self-hosted service status dashboard. Built with Vue 3 & Go.</p>
        </CardFooter>
      </Card>
    </div>
  </div>
</template>
