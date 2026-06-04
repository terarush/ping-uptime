<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useAuth } from '@/composables/useAuth';
import ExtendedFetch from '@/lib/fetch';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog';
import {
  Users,
  UserPlus,
  Trash2,
  Edit,
  Mail,
  Search,
  Lock,
  Loader2,
  Shield,
  UserCheck,
  ShieldAlert,
  UserX,
  Plus
} from '@lucide/vue';
import gsap from 'gsap';

// Types
interface User {
  id: number;
  name: string;
  email: string;
  role: string;
  is_blocked: boolean;
  avatar?: string;
  created_at: string;
}

// Refs & States
const users = ref<User[]>();
const loading = ref(true);
const searchQuery = ref('');
const error = ref('');
const success = ref('');

// Current logged in user (to prevent self-deletion or self-blocking)
const { currentUser } = useAuth();

// Dialog states
const isFormDialogOpen = ref(false);
const isDeleteDialogOpen = ref(false);
const formLoading = ref(false);
const actionUser = ref<User | null>(null); // User to edit or delete

// Form field refs
const formName = ref('');
const formEmail = ref('');
const formPassword = ref('');
const formRole = ref('user');
const formIsBlocked = ref(false);

const isEditMode = computed(() => !!actionUser.value);

// Fetch all users
const fetchUsers = async () => {
  loading.value = true;
  error.value = '';
  try {
    const response = await ExtendedFetch.get('/users');
    // Ensure we parse correctly depending on standardized response wrapper
    users.value = response.data?.data || response.data || [];
  } catch (err: any) {
    console.error('Failed to fetch users:', err);
    error.value = err.response?.data?.error || 'Failed to load user accounts.';
  } finally {
    loading.value = false;
    setTimeout(animateTableRows, 50);
  }
};

// Filter users by search query
const filteredUsers = computed(() => {
  if (!users.value) return [];
  const query = searchQuery.value.toLowerCase().trim();
  if (!query) return users.value;
  return users.value.filter(u =>
    u.name.toLowerCase().includes(query) ||
    u.email.toLowerCase().includes(query) ||
    u.role.toLowerCase().includes(query)
  );
});

// Reset form
const resetForm = () => {
  formName.value = '';
  formEmail.value = '';
  formPassword.value = '';
  formRole.value = 'user';
  formIsBlocked.value = false;
  actionUser.value = null;
  error.value = '';
  success.value = '';
};

// Open Dialog for Create
const openCreateDialog = () => {
  resetForm();
  isFormDialogOpen.value = true;
};

// Open Dialog for Edit
const openEditDialog = (user: User) => {
  resetForm();
  actionUser.value = user;
  formName.value = user.name;
  formEmail.value = user.email;
  formRole.value = user.role;
  formIsBlocked.value = user.is_blocked;
  isFormDialogOpen.value = true;
};

// Open Dialog for Delete Confirmation
const openDeleteDialog = (user: User) => {
  actionUser.value = user;
  isDeleteDialogOpen.value = true;
};

// Submit form (Create / Edit)
const handleFormSubmit = async () => {
  error.value = '';
  success.value = '';
  formLoading.value = true;

  try {
    const payload: any = {
      name: formName.value,
      email: formEmail.value,
      role: formRole.value,
      is_blocked: formIsBlocked.value,
    };

    if (formPassword.value) {
      payload.password = formPassword.value;
    }

    if (isEditMode.value && actionUser.value) {
      await ExtendedFetch.put(`/users/${actionUser.value.id}`, payload);
      success.value = `User "${formName.value}" updated successfully!`;
    } else {
      if (!formPassword.value) {
        error.value = 'Password is required for new users.';
        formLoading.value = false;
        return;
      }
      payload.password = formPassword.value;
      await ExtendedFetch.post('/users', payload);
      success.value = `User "${formName.value}" created successfully!`;
    }

    isFormDialogOpen.value = false;
    await fetchUsers();
  } catch (err: any) {
    console.error('Failed to save user:', err);
    error.value = err.response?.data?.error || 'Failed to save user data.';
  } finally {
    formLoading.value = false;
  }
};

// Delete user
const handleDeleteConfirm = async () => {
  if (!actionUser.value) return;
  error.value = '';
  formLoading.value = true;

  try {
    await ExtendedFetch.delete(`/users/${actionUser.value.id}`);
    success.value = `User "${actionUser.value.name}" deleted successfully!`;
    isDeleteDialogOpen.value = false;
    await fetchUsers();
  } catch (err: any) {
    console.error('Failed to delete user:', err);
    error.value = err.response?.data?.error || 'Failed to delete user.';
  } finally {
    formLoading.value = false;
  }
};

// GSAP Animations
const animateTableRows = () => {
  gsap.fromTo('.user-row',
    { opacity: 0, y: 15 },
    { opacity: 1, y: 0, duration: 0.4, stagger: 0.05, ease: 'power2.out' }
  );
};

onMounted(() => {
  fetchUsers();
  gsap.fromTo('.ambient-orb',
    { opacity: 0, scale: 0.8 },
    { opacity: 0.6, scale: 1, duration: 2.5, ease: 'power3.out' }
  );
});
</script>

<template>
  <div class="relative p-8 space-y-8 max-w-400 mx-auto min-h-[calc(100vh-4rem)]">
    <!-- Ambient Background Orb -->
    <div class="ambient-orb absolute top-[-5%] right-[-5%] w-[45%] h-[45%] rounded-full bg-emerald-500/10 dark:bg-emerald-500/5 blur-[100px] pointer-events-none"></div>

    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 z-10 relative">
      <div>
        <h2 class="text-2xl font-black tracking-tight text-foreground flex items-center gap-2">
          <Users class="w-6 h-6 text-primary" />
          <span>User Accounts</span>
        </h2>
        <p class="text-xs text-muted-foreground">Manage administrator permissions and dashboard access accounts.</p>
      </div>

      <!-- Add User Trigger -->
      <Button @click="openCreateDialog" class="inline-flex items-center gap-1.5 shadow-sm shadow-primary/10">
        <UserPlus class="w-4 h-4" />
        <span>Create User</span>
      </Button>
    </div>

    <!-- Main Container Card -->
    <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md z-10 relative">
      <CardHeader class="pb-3 border-b border-border/40">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <CardTitle class="text-sm font-bold text-foreground">Registered Users</CardTitle>
            <CardDescription class="text-xs">A list of all users registered in this self-hosted system.</CardDescription>
          </div>

          <!-- Search Bar -->
          <div class="relative w-full sm:w-72">
            <Search class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              v-model="searchQuery"
              placeholder="Search by name, email or role..."
              class="pl-9 h-9 rounded-lg"
            />
          </div>
        </div>
      </CardHeader>

      <CardContent class="p-0">
        <!-- Error alert -->
        <div v-if="error" class="m-6 p-3 text-xs bg-destructive/10 border border-destructive/20 text-destructive rounded-lg flex items-start gap-2">
          <ShieldAlert class="h-4 w-4 shrink-0 mt-0.5" />
          <span>{{ error }}</span>
        </div>

        <!-- Success alert -->
        <div v-if="success" class="m-6 p-3 text-xs bg-emerald-500/10 border border-emerald-500/20 text-emerald-600 dark:text-emerald-400 rounded-lg flex items-start gap-2">
          <UserCheck class="h-4 w-4 shrink-0 mt-0.5" />
          <span>{{ success }}</span>
        </div>

        <!-- Loading Skeleton -->
        <div v-if="loading" class="p-8 space-y-4">
          <div v-for="n in 3" :key="n" class="flex items-center space-x-4">
            <div class="h-10 w-10 rounded-full bg-muted animate-pulse"></div>
            <div class="space-y-2 flex-1">
              <div class="h-4 bg-muted rounded w-1/3 animate-pulse"></div>
              <div class="h-3 bg-muted rounded w-1/4 animate-pulse"></div>
            </div>
            <div class="h-6 w-16 bg-muted rounded animate-pulse"></div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else-if="filteredUsers.length === 0" class="p-16 flex flex-col items-center text-center space-y-3">
          <UserX class="w-12 h-12 text-muted-foreground/60" />
          <h3 class="text-sm font-bold text-foreground">No Users Found</h3>
          <p class="text-xs text-muted-foreground max-w-64">No user records matched your criteria or search query.</p>
        </div>

        <!-- User Table -->
        <Table v-else>
          <TableHeader>
            <TableRow>
              <TableHead class="w-14"></TableHead>
              <TableHead class="text-xs font-bold uppercase text-muted-foreground">User</TableHead>
              <TableHead class="text-xs font-bold uppercase text-muted-foreground">Role</TableHead>
              <TableHead class="text-xs font-bold uppercase text-muted-foreground">Status</TableHead>
              <TableHead class="text-xs font-bold uppercase text-muted-foreground text-right pr-6">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="user in filteredUsers" :key="user.id" class="user-row hover:bg-muted/40 transition-colors">
              <!-- Avatar -->
              <TableCell class="py-4 pl-6">
                <Avatar class="h-9 w-9 ring-1 ring-primary/10">
                  <AvatarImage :src="user.avatar || `https://api.dicebear.com/7.x/initials/svg?seed=${user.name}`" alt="Avatar" />
                  <AvatarFallback class="font-bold text-xs uppercase">{{ user.name.slice(0, 2) }}</AvatarFallback>
                </Avatar>
              </TableCell>

              <!-- Profile Details -->
              <TableCell class="py-4">
                <div class="flex flex-col gap-0.5">
                  <span class="text-xs font-bold text-foreground flex items-center gap-1.5">
                    {{ user.name }}
                    <Badge v-if="currentUser && currentUser.id === user.id" variant="secondary" class="text-[9px] font-bold px-1 py-0 h-4 bg-muted text-muted-foreground">You</Badge>
                  </span>
                  <span class="text-[10px] text-muted-foreground flex items-center gap-1">
                    <Mail class="w-2.5 h-2.5" />
                    {{ user.email }}
                  </span>
                </div>
              </TableCell>

              <!-- Role Badge -->
              <TableCell class="py-4">
                <Badge variant="outline" :class="[
                  'text-[10px] font-bold py-0.5 px-2 border/50',
                  user.role === 'admin'
                    ? 'bg-primary/5 text-primary border-primary/20'
                    : 'bg-muted text-muted-foreground border-border'
                ]">
                  <Shield v-if="user.role === 'admin'" class="w-2.5 h-2.5 inline mr-1" />
                  {{ user.role === 'admin' ? 'Administrator' : 'Standard User' }}
                </Badge>
              </TableCell>

              <!-- Status Badge -->
              <TableCell class="py-4">
                <Badge variant="outline" :class="[
                  'text-[10px] font-bold py-0.5 px-2 border/50',
                  user.is_blocked
                    ? 'bg-destructive/5 text-destructive border-destructive/20'
                    : 'bg-emerald-500/5 text-emerald-600 dark:text-emerald-400 border-emerald-500/20'
                ]">
                  {{ user.is_blocked ? 'Blocked' : 'Active' }}
                </Badge>
              </TableCell>

              <!-- Action Menu -->
              <TableCell class="py-4 text-right pr-6">
                <div class="flex items-center justify-end gap-1.5">
                  <!-- Edit Action -->
                  <Button @click="openEditDialog(user)" size="icon" variant="ghost" class="h-8 w-8 text-muted-foreground hover:text-foreground">
                    <Edit class="w-4 h-4" />
                  </Button>

                  <!-- Delete Action (Disabled for logged-in user itself) -->
                  <Button
                    @click="openDeleteDialog(user)"
                    size="icon"
                    variant="ghost"
                    class="h-8 w-8 text-muted-foreground hover:text-destructive"
                    :disabled="currentUser && currentUser.id === user.id"
                  >
                    <Trash2 class="w-4 h-4" />
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>

    <!-- Dialog Modal: Create or Edit User -->
    <Dialog v-model:open="isFormDialogOpen">
      <DialogContent class="sm:max-w-md bg-card border-border/50">
        <DialogHeader>
          <DialogTitle class="text-foreground font-bold tracking-tight">
            {{ isEditMode ? 'Edit User details' : 'Create New Account' }}
          </DialogTitle>
          <DialogDescription class="text-muted-foreground text-xs">
            {{ isEditMode ? 'Modify dashboard access permissions and profile information.' : 'Add a new administrator or standard operator to this panel.' }}
          </DialogDescription>
        </DialogHeader>

        <!-- Form fields -->
        <div class="space-y-4 py-4">
          <!-- Full Name -->
          <div class="space-y-1.5">
            <Label for="name" class="text-xs font-bold text-foreground/80">Full Name</Label>
            <Input id="name" v-model="formName" placeholder="e.g. John Doe" class="h-10 rounded-lg" />
          </div>

          <!-- Email -->
          <div class="space-y-1.5">
            <Label for="email" class="text-xs font-bold text-foreground/80">Email Address</Label>
            <Input id="email" type="email" v-model="formEmail" placeholder="e.g. name@example.com" class="h-10 rounded-lg" />
          </div>

          <!-- Password -->
          <div class="space-y-1.5">
            <div class="flex items-center justify-between">
              <Label for="password" class="text-xs font-bold text-foreground/80">Password</Label>
              <span v-if="isEditMode" class="text-[9px] text-muted-foreground">(Leave empty to keep current password)</span>
            </div>
            <div class="relative">
              <Lock class="absolute left-3 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                id="password"
                type="password"
                v-model="formPassword"
                placeholder="••••••••"
                class="pl-9 h-10 rounded-lg"
              />
            </div>
          </div>

          <!-- Role selection -->
          <div class="space-y-1.5">
            <Label class="text-xs font-bold text-foreground/80">Access Privileges</Label>
            <select
              v-model="formRole"
              class="flex h-10 w-full rounded-lg border border-input bg-card px-3 py-2 text-sm text-foreground ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
            >
              <option value="user">Standard User (Operators)</option>
              <option value="admin">Administrator (Full Access)</option>
            </select>
          </div>

          <!-- Block Toggle (Only visible in edit mode and for users other than current self) -->
          <div v-if="isEditMode && currentUser && currentUser.id !== actionUser?.id" class="flex items-center justify-between p-3 border border-border/50 rounded-xl bg-muted/20">
            <div class="flex flex-col gap-0.5">
              <span class="text-xs font-bold text-foreground">Block Account</span>
              <span class="text-[10px] text-muted-foreground">Deny user from logging in and accessing APIs.</span>
            </div>
            <input
              type="checkbox"
              v-model="formIsBlocked"
              class="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
            />
          </div>
        </div>

        <DialogFooter class="flex sm:justify-end gap-2">
          <Button variant="outline" @click="isFormDialogOpen = false" :disabled="formLoading">
            Cancel
          </Button>
          <Button @click="handleFormSubmit" :disabled="formLoading || !formName || !formEmail" class="gap-1.5">
            <Loader2 v-if="formLoading" class="w-4 h-4 animate-spin" />
            <span>{{ isEditMode ? 'Update Account' : 'Register User' }}</span>
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Dialog Modal: Delete Confirmation -->
    <Dialog v-model:open="isDeleteDialogOpen">
      <DialogContent class="sm:max-w-md bg-card border-border/50">
        <DialogHeader>
          <DialogTitle class="text-foreground font-bold tracking-tight flex items-center gap-2 text-destructive">
            <ShieldAlert class="w-5 h-5" />
            <span>Delete User Account</span>
          </DialogTitle>
          <DialogDescription class="text-muted-foreground text-xs leading-relaxed">
            Are you absolutely sure you want to delete user <span class="font-bold text-foreground">"{{ actionUser?.name }}"</span>? This action is permanent and cannot be undone. All logs and metadata associated with this user will be removed.
          </DialogDescription>
        </DialogHeader>

        <DialogFooter class="flex sm:justify-end gap-2">
          <Button variant="outline" @click="isDeleteDialogOpen = false" :disabled="formLoading">
            Cancel
          </Button>
          <Button variant="destructive" @click="handleDeleteConfirm" :disabled="formLoading" class="gap-1.5">
            <Loader2 v-if="formLoading" class="w-4 h-4 animate-spin" />
            <span>Delete Account</span>
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
