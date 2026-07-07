<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { Mail, ArrowLeft, Loader2, AlertCircle, CheckCircle2 } from '@lucide/vue';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from '@/components/ui/card';
import ThemeButton from '@/components/theme-button.vue';
import ExtendedFetch from '@/lib/fetch';
import gsap from 'gsap';
import { Form, FormField, FormItem, FormLabel, FormControl, FormMessage } from '@/components/ui/form';
import { toTypedSchema } from '@vee-validate/zod';
import { z } from 'zod';

const forgotSchema = z.object({
  email: z.string({ required_error: 'Email is required' }).min(1, 'Email is required').email('Enter a valid email'),
});
const formSchema = toTypedSchema(forgotSchema);

const loading = ref(false);
const sent = ref(false);
const error = ref('');

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
});

const onSubmit = async (values: { email: string }) => {
  error.value = '';
  loading.value = true;
  try {
    await ExtendedFetch.post('/auth/forgot-password', { email: values.email });
    sent.value = true;
  } catch (err: any) {
    error.value = err.response?.data?.error || err.message || 'Something went wrong';
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
        <h1 class="text-2xl font-black tracking-tight text-foreground select-none">Forgot Password</h1>
        <p class="text-xs text-muted-foreground max-w-70 mt-1">Enter your email to receive a password reset link.</p>
      </div>

      <Card class="login-card w-full bg-card/60 dark:bg-card/40 backdrop-blur-xl border border-border/60 dark:border-border/20 shadow-2xl rounded-2xl overflow-hidden p-6 gap-0">
        <CardHeader class="p-0 pb-6 text-center">
          <CardTitle class="text-lg font-bold tracking-tight text-foreground">Reset your password</CardTitle>
          <CardDescription class="text-xs text-muted-foreground mt-1">
            We'll send you a reset link if the account exists.
          </CardDescription>
        </CardHeader>

        <CardContent class="p-0 space-y-4">
          <div v-if="error" class="p-3 text-xs bg-destructive/10 border border-destructive/20 text-destructive rounded-lg flex items-start gap-2">
            <AlertCircle class="h-4 w-4 shrink-0 mt-0.5" />
            <span>{{ error }}</span>
          </div>

          <div v-if="sent" class="p-3 text-xs bg-emerald-500/10 border border-emerald-500/20 text-emerald-600 dark:text-emerald-400 rounded-lg flex items-start gap-2">
            <CheckCircle2 class="h-4 w-4 shrink-0 mt-0.5" />
            <span>If that email is registered, you'll receive a password reset link shortly.</span>
          </div>

          <Form v-if="!sent" :validation-schema="formSchema" @submit="onSubmit" class="space-y-4">
            <FormField name="email" v-slot="{ componentField }">
              <FormItem class="space-y-1.5">
                <FormLabel class="text-xs font-bold text-foreground/80">Email Address</FormLabel>
                <FormControl>
                  <div class="relative">
                    <Mail class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
                    <Input v-bind="componentField" type="email" placeholder="admin@example.com" class="pl-9 h-10 rounded-lg" />
                  </div>
                </FormControl>
                <FormMessage class="text-[10px] text-destructive" />
              </FormItem>
            </FormField>

            <Button type="submit" :disabled="loading"
              class="w-full h-10 rounded-lg font-bold shadow-lg hover:shadow-xl transition-all duration-300 hover:scale-[1.01] bg-primary hover:bg-primary/90 text-primary-foreground flex items-center justify-center gap-2 cursor-pointer mt-2">
              <Loader2 v-if="loading" class="h-4 w-4 animate-spin" />
              <span v-else>Send Reset Link</span>
            </Button>
          </Form>

          <div class="text-center text-xs mt-4 pt-2 border-t border-border/40">
            <RouterLink to="/" class="inline-flex items-center gap-1 text-primary hover:underline font-semibold">
              <ArrowLeft class="h-3 w-3" /> Back to Sign In
            </RouterLink>
          </div>
        </CardContent>

        <CardFooter class="p-0 pt-6 text-center justify-center">
          <p class="text-[10px] text-muted-foreground">Self-hosted service status dashboard. Built with Vue 3 & Go.</p>
        </CardFooter>
      </Card>
    </div>
  </div>
</template>
