<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useAuth } from '@/composables/useAuth';
import { useUsers, type User } from '@/composables/useUsers';
import { userSchema } from '@/validations/user';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import UserTable from '@/components/user-table.vue';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog';
import {
  Users as UsersIcon,
  UserPlus,
  Lock,
  Loader2,
  Search,
  ShieldAlert,
  UserCheck,
  UserX
} from '@lucide/vue';
import gsap from 'gsap';
import { toast } from 'vue-sonner';

// Current logged in user (to prevent self-deletion or self-blocking)
const { currentUser } = useAuth();

// Use external Users composable
const {
  users,
  loading,
  error,
  fetchUsersData,
  createUser,
  updateUser,
  deleteUser
} = useUsers();

const searchQuery = ref('');
const filterUserRole = ref('__all__');
const success = ref('');

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

// Fetch all users with animation callback
const fetchAll = async () => {
  try {
    await fetchUsersData();
  } catch (err) {
    console.error('Failed to load users:', err);
  } finally {
    setTimeout(animateTableRows, 50);
  }
};

// Filter users by search query + role
const filteredUsers = computed(() => {
  if (!users.value) return [];
  let result = users.value;
  const q = searchQuery.value.toLowerCase().trim();
  if (q) {
    result = result.filter(u =>
      u.name.toLowerCase().includes(q) ||
      u.email.toLowerCase().includes(q)
    );
  }
  if (filterUserRole.value && filterUserRole.value !== '__all__') {
    result = result.filter(u => u.role === filterUserRole.value);
  }
  return result;
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

  const rawPayload: any = {
    name: formName.value,
    email: formEmail.value,
    role: formRole.value,
    is_blocked: formIsBlocked.value,
    password: formPassword.value || undefined,
  };

  // Perform client-side Zod validation
  const validation = userSchema.safeParse(rawPayload);
  if (!validation.success) {
    const firstError = validation.error.errors[0]?.message || 'Validation failed';
    toast.error(firstError);
    error.value = firstError;
    formLoading.value = false;
    return;
  }

  try {
    if (isEditMode.value && actionUser.value) {
      await updateUser(actionUser.value.id, rawPayload);
      toast.success(`User "${formName.value}" updated successfully!`);
      success.value = `User "${formName.value}" updated successfully!`;
    } else {
      if (!formPassword.value) {
        error.value = 'Password is required for new users.';
        toast.error('Password is required for new users.');
        formLoading.value = false;
        return;
      }
      await createUser(rawPayload);
      toast.success(`User "${formName.value}" created successfully!`);
      success.value = `User "${formName.value}" created successfully!`;
    }

    isFormDialogOpen.value = false;
    await fetchAll();
  } catch (err: any) {
    console.error('Failed to save user:', err);
    const msg = err.response?.data?.error || 'Failed to save user data.';
    toast.error(msg);
    error.value = msg;
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
    await deleteUser(actionUser.value.id);
    toast.success(`User "${actionUser.value.name}" deleted successfully!`);
    success.value = `User "${actionUser.value.name}" deleted successfully!`;
    isDeleteDialogOpen.value = false;
    await fetchAll();
  } catch (err: any) {
    console.error('Failed to delete user:', err);
    const msg = err.response?.data?.error || 'Failed to delete user.';
    toast.error(msg);
    error.value = msg;
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
  fetchAll();
  gsap.fromTo('.ambient-orb',
    { opacity: 0, scale: 0.8 },
    { opacity: 0.6, scale: 1, duration: 2.5, ease: 'power3.out' }
  );
});
</script>

<template>
  <div class="relative p-8 space-y-8 max-w-7xl mx-auto min-h-[calc(100vh-4rem)]">
    <!-- Ambient Background Orb -->
    <div class="ambient-orb absolute top-[-5%] right-[-5%] w-[45%] h-[45%] rounded-full bg-emerald-500/10 dark:bg-emerald-500/5 blur-[100px] pointer-events-none"></div>

    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 z-10 relative">
      <div>
        <h2 class="text-2xl font-black tracking-tight text-foreground flex items-center gap-2">
          <UsersIcon class="w-6 h-6 text-primary" />
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
      <CardHeader class="pb-3 border-b border-border/40 space-y-3">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
          <div>
            <CardTitle class="text-sm font-bold text-foreground">Registered Users</CardTitle>
            <CardDescription class="text-xs">A list of all users registered in this self-hosted system.</CardDescription>
          </div>

          <!-- Search Bar -->
          <div class="relative w-full sm:w-56">
            <Search class="absolute left-2.5 top-2 h-4 w-4 text-muted-foreground" />
            <Input v-model="searchQuery" placeholder="Search by name or email..." class="pl-8 h-8 rounded-lg" />
          </div>
        </div>
        <div class="flex flex-wrap items-center gap-2">
          <Select v-model="filterUserRole">
            <SelectTrigger class="w-28 h-8">
              <SelectValue placeholder="Role" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="__all__">All roles</SelectItem>
              <SelectItem value="admin">Admin</SelectItem>
              <SelectItem value="user">User</SelectItem>
            </SelectContent>
          </Select>
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
        <div v-if="loading && users.length === 0" class="p-8 space-y-4">
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

        <!-- User Table Component -->
        <UserTable
          v-else
          :users="filteredUsers"
          :current-user="currentUser"
          @edit="openEditDialog"
          @delete="openDeleteDialog"
        />
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
