<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useTeams, type Team, type TeamMember } from '@/composables/useTeams';
import { useUsers, type User } from '@/composables/useUsers';
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
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
  UserMinus,
  UserCheck,
  UserX,
  Loader2,
  ShieldAlert,
  Plus,
  Send,
  LogIn,
  X,
  UsersRound,
  Mail,
  BadgeCheck,
  Clock,
  Ban,
  Search,
} from '@lucide/vue';
import gsap from 'gsap';
import { toast } from 'vue-sonner';

const {
  teams,
  members,
  loading,
  error,
  fetchAll,
  create,
  update,
  deleteTeam,
  fetchMembers,
  inviteMember,
  updateMember,
  removeMember,
  acceptInvite,
  rejectInvite,
} = useTeams();

const { users, fetchUsersData } = useUsers();

// Dialog states
const isCreateDialogOpen = ref(false);
const isEditDialogOpen = ref(false);
const isDeleteDialogOpen = ref(false);
const isInviteDialogOpen = ref(false);
const formLoading = ref(false);
const actionTeam = ref<Team | null>(null);
const selectedTeam = ref<Team | null>(null);
const actionMember = ref<TeamMember | null>(null);

// Form fields
const formName = ref('');
const inviteUserID = ref<number>(0);
const inviteRole = ref('member');

const isEditMode = computed(() => !!actionTeam.value);

const searchQuery = ref('');

const filteredTeams = computed(() => {
  if (!searchQuery.value.trim()) return teams.value;
  const q = searchQuery.value.toLowerCase().trim();
  return teams.value.filter(t => t.name.toLowerCase().includes(q));
});

// Filter members for selected team
const teamMembers = computed(() => {
  if (!selectedTeam.value) return [];
  return members.value.filter(m => m.team_id === selectedTeam.value!.id);
});

const fetchAllData = async () => {
  try {
    await fetchAll();
    await fetchUsersData();
  } catch (err) {
    console.error('Failed to load data:', err);
  } finally {
    setTimeout(animateRows, 50);
  }
};

// Reset form
const resetForm = () => {
  formName.value = '';
  actionTeam.value = null;
  error.value = '';
};

const resetInviteForm = () => {
  inviteUserID.value = 0;
  inviteRole.value = 'member';
  error.value = '';
};

// Team CRUD
const openCreateDialog = () => {
  resetForm();
  isCreateDialogOpen.value = true;
};

const openEditDialog = (team: Team) => {
  resetForm();
  actionTeam.value = team;
  formName.value = team.name;
  isEditDialogOpen.value = true;
};

const openDeleteDialog = (team: Team) => {
  actionTeam.value = team;
  isDeleteDialogOpen.value = true;
};

const openInviteDialog = () => {
  resetInviteForm();
  isInviteDialogOpen.value = true;
};

const handleCreate = async () => {
  if (!formName.value) return;
  error.value = '';
  formLoading.value = true;
  try {
    await create(formName.value);
    toast.success(`Team "${formName.value}" created!`);
    isCreateDialogOpen.value = false;
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Failed to create team.');
  } finally {
    formLoading.value = false;
  }
};

const handleUpdate = async () => {
  if (!formName.value || !actionTeam.value) return;
  error.value = '';
  formLoading.value = true;
  try {
    await update(actionTeam.value.id, formName.value);
    toast.success(`Team updated to "${formName.value}"`);
    isEditDialogOpen.value = false;
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Failed to update team.');
  } finally {
    formLoading.value = false;
  }
};

const handleDeleteConfirm = async () => {
  if (!actionTeam.value) return;
  formLoading.value = true;
  try {
    await deleteTeam(actionTeam.value.id);
    toast.success(`Team "${actionTeam.value.name}" deleted`);
    if (selectedTeam.value?.id === actionTeam.value.id) {
      selectedTeam.value = null;
      members.value = [];
    }
    isDeleteDialogOpen.value = false;
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Failed to delete team.');
  } finally {
    formLoading.value = false;
  }
};

// Member actions
const selectTeam = async (team: Team) => {
  selectedTeam.value = team;
  await fetchMembers(team.id);
  setTimeout(animateRows, 50);
};

const handleInvite = async () => {
  if (!selectedTeam.value || !inviteUserID.value) return;
  error.value = '';
  formLoading.value = true;
  try {
    await inviteMember(selectedTeam.value.id, inviteUserID.value, inviteRole.value);
    toast.success('Invitation sent!');
    isInviteDialogOpen.value = false;
    await fetchMembers(selectedTeam.value.id);
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Failed to invite member.');
  } finally {
    formLoading.value = false;
  }
};

const handleRoleChange = async (member: TeamMember, role: string) => {
  if (!selectedTeam.value) return;
  try {
    await updateMember(selectedTeam.value.id, member.id, role);
    toast.success(`Member role updated to ${role}`);
    await fetchMembers(selectedTeam.value.id);
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Failed to update role.');
  }
};

const handleRemoveMember = async (member: TeamMember) => {
  if (!selectedTeam.value) return;
  try {
    await removeMember(selectedTeam.value.id, member.id);
    toast.success('Member removed');
    await fetchMembers(selectedTeam.value.id);
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Failed to remove member.');
  }
};

// Invitation actions
const handleAcceptInvite = async (teamID: number) => {
  try {
    await acceptInvite(teamID);
    toast.success('Invitation accepted!');
    await fetchMembers(teamID);
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Failed to accept invitation.');
  }
};

const handleRejectInvite = async (teamID: number) => {
  try {
    await rejectInvite(teamID);
    toast.success('Invitation rejected');
    await fetchMembers(teamID);
  } catch (err: any) {
    toast.error(err.response?.data?.error || 'Failed to reject invitation.');
  }
};

// Helpers
const getUserName = (userID: number): string => {
  const u = users.value.find(u => u.id === userID);
  return u?.name || u?.email || `User #${userID}`;
};

const getUserEmail = (userID: number): string => {
  const u = users.value.find(u => u.id === userID);
  return u?.email || '';
};

const statusBadge = (status: string) => {
  switch (status) {
    case 'accepted': return 'bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border-emerald-500/25';
    case 'pending': return 'bg-amber-500/15 text-amber-600 dark:text-amber-400 border-amber-500/25';
    case 'rejected': return 'bg-red-500/15 text-red-600 dark:text-red-400 border-red-500/25';
    default: return 'bg-muted text-muted-foreground border-border';
  }
};

const roleBadge = (role: string) => {
  return role === 'admin'
    ? 'bg-purple-500/15 text-purple-600 dark:text-purple-400 border-purple-500/25'
    : 'bg-blue-500/15 text-blue-600 dark:text-blue-400 border-blue-500/25';
};

// GSAP Animations
const animateRows = () => {
  gsap.fromTo('.team-row',
    { opacity: 0, y: 15 },
    { opacity: 1, y: 0, duration: 0.4, stagger: 0.05, ease: 'power2.out' }
  );
  gsap.fromTo('.member-row',
    { opacity: 0, x: -10 },
    { opacity: 1, x: 0, duration: 0.3, stagger: 0.03, ease: 'power2.out' }
  );
};

onMounted(() => {
  fetchAllData();
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
          <span>Teams &amp; Members</span>
        </h2>
        <p class="text-xs text-muted-foreground">Create and manage teams with role-based member access.</p>
      </div>
      <Button @click="openCreateDialog" class="inline-flex items-center gap-1.5 shadow-sm shadow-primary/10">
        <Plus class="w-4 h-4" />
        <span>Create Team</span>
      </Button>
    </div>

    <!-- Main Container -->
    <div class="grid grid-cols-1 lg:grid-cols-5 gap-6 z-10 relative">
      <!-- Team List -->
      <div class="lg:col-span-2">
        <Card class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
          <CardHeader class="pb-3 border-b border-border/40">
            <CardTitle class="text-sm font-bold text-foreground">All Teams</CardTitle>
            <CardDescription class="text-xs">{{ teams.length }} team{{ teams.length !== 1 ? 's' : '' }}</CardDescription>
            <div class="relative w-full mt-2">
              <Search class="absolute left-3 top-2 h-4 w-4 text-muted-foreground" />
              <Input v-model="searchQuery" placeholder="Search teams..." class="pl-8 h-8" />
            </div>
          </CardHeader>
          <CardContent class="p-0">
            <!-- Loading Skeleton -->
            <div v-if="loading && teams.length === 0" class="p-6 space-y-3">
              <div v-for="n in 3" :key="n" class="h-14 bg-muted animate-pulse rounded-lg"></div>
            </div>

            <!-- Empty State -->
            <div v-else-if="filteredTeams.length === 0" class="p-12 flex flex-col items-center text-center space-y-3">
              <UsersRound class="w-10 h-10 text-muted-foreground/60" />
              <h3 class="text-sm font-bold text-foreground">No Teams Yet</h3>
              <p class="text-xs text-muted-foreground max-w-48">Create your first team to start collaborating.</p>
            </div>

            <!-- Team Cards -->
            <div v-else class="divide-y divide-border/40">
              <div
                v-for="team in filteredTeams" :key="team.id"
                class="team-row flex items-center justify-between p-4 cursor-pointer transition-colors hover:bg-muted/40"
                :class="{ 'bg-muted/30 border-l-2 border-l-primary': selectedTeam?.id === team.id }"
                @click="selectTeam(team)"
              >
                <div class="flex items-center gap-3 min-w-0">
                  <div class="w-9 h-9 rounded-full bg-primary/10 flex items-center justify-center shrink-0">
                    <UsersIcon class="w-4 h-4 text-primary" />
                  </div>
                  <div class="min-w-0">
                    <p class="text-sm font-semibold text-foreground truncate">{{ team.name }}</p>
                    <p class="text-[10px] text-muted-foreground">
                      {{ team.member_count ?? 0 }} member{{ (team.member_count ?? 0) !== 1 ? 's' : '' }}
                    </p>
                  </div>
                </div>
                <div class="flex items-center gap-1 shrink-0">
                  <Button variant="ghost" size="icon" class="h-7 w-7" @click.stop="openEditDialog(team)">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="w-3.5 h-3.5"><path d="M17 3a2.85 2.85 0 1 1 4 4L7.5 20.5 2 22l1.5-5.5Z"/><path d="m15 5 4 4"/></svg>
                  </Button>
                  <Button variant="ghost" size="icon" class="h-7 w-7 text-destructive" @click.stop="openDeleteDialog(team)">
                    <X class="w-3.5 h-3.5" />
                  </Button>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      <!-- Members Panel -->
      <div class="lg:col-span-3">
        <Card v-if="selectedTeam" class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
          <CardHeader class="pb-3 border-b border-border/40">
            <div class="flex items-center justify-between">
              <div>
                <CardTitle class="text-sm font-bold text-foreground">{{ selectedTeam.name }} &mdash; Members</CardTitle>
                <CardDescription class="text-xs">{{ teamMembers.length }} member{{ teamMembers.length !== 1 ? 's' : '' }}</CardDescription>
              </div>
              <Button size="sm" variant="outline" class="gap-1" @click="openInviteDialog">
                <UserPlus class="w-3.5 h-3.5" />
                <span>Invite</span>
              </Button>
            </div>
          </CardHeader>
          <CardContent class="p-0">
            <!-- Loading -->
            <div v-if="loading && members.length === 0" class="p-6 space-y-3">
              <div v-for="n in 3" :key="n" class="h-12 bg-muted animate-pulse rounded-lg"></div>
            </div>

            <!-- Empty State -->
            <div v-else-if="teamMembers.length === 0" class="p-12 flex flex-col items-center text-center space-y-3">
              <UserPlus class="w-10 h-10 text-muted-foreground/60" />
              <h3 class="text-sm font-bold text-foreground">No Members Yet</h3>
              <p class="text-xs text-muted-foreground max-w-48">Invite users to this team to collaborate.</p>
            </div>

            <!-- Member Table -->
            <div v-else class="divide-y divide-border/40">
              <div
                v-for="member in teamMembers"
                :key="member.id"
                class="member-row flex items-center justify-between p-3.5 hover:bg-muted/20 transition-colors"
              >
                <div class="flex items-center gap-3 min-w-0 flex-1">
                  <div class="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center shrink-0 text-xs font-bold text-primary">
                    {{ (member.user?.name || getUserName(member.user_id)).charAt(0).toUpperCase() }}
                  </div>
                  <div class="min-w-0">
                    <p class="text-sm font-semibold text-foreground truncate">
                      {{ member.user?.name || getUserName(member.user_id) }}
                    </p>
                    <p class="text-[10px] text-muted-foreground truncate">{{ member.user?.email || getUserEmail(member.user_id) }}</p>
                  </div>
                </div>

                <div class="flex items-center gap-2 shrink-0">
                  <!-- Role Badge -->
                  <span class="text-[10px] font-semibold px-2 py-0.5 rounded-full border" :class="roleBadge(member.role)">
                    {{ member.role }}
                  </span>

                  <!-- Status Pill -->
                  <span v-if="member.status !== 'accepted'" class="text-[10px] font-semibold px-2 py-0.5 rounded-full border flex items-center gap-1" :class="statusBadge(member.status)">
                    <template v-if="member.status === 'pending'"><Clock class="w-2.5 h-2.5" /> Pending</template>
                    <template v-else-if="member.status === 'rejected'"><Ban class="w-2.5 h-2.5" /> Rejected</template>
                  </span>

                  <!-- Actions -->
                  <div class="flex items-center gap-1">
                    <!-- Role toggle (only for accepted members) -->
                    <select
                      v-if="member.status === 'accepted'"
                      :value="member.role"
                      @change="(e) => handleRoleChange(member, (e.target as HTMLSelectElement).value)"
                      class="h-7 text-[10px] rounded border border-border bg-card text-foreground px-1"
                    >
                      <option value="member">member</option>
                      <option value="admin">admin</option>
                    </select>

                    <!-- Accept/Reject for pending invites (shown to the invitee) -->
                    <Button v-if="member.status === 'pending'" variant="ghost" size="icon" class="h-7 w-7 text-emerald-500" @click="handleAcceptInvite(selectedTeam.id)">
                      <UserCheck class="w-3.5 h-3.5" />
                    </Button>
                    <Button v-if="member.status === 'pending'" variant="ghost" size="icon" class="h-7 w-7 text-destructive" @click="handleRejectInvite(selectedTeam.id)">
                      <UserX class="w-3.5 h-3.5" />
                    </Button>

                    <!-- Remove -->
                    <Button variant="ghost" size="icon" class="h-7 w-7 text-destructive" @click="handleRemoveMember(member)">
                      <UserMinus class="w-3.5 h-3.5" />
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- No team selected -->
        <Card v-else class="border-border/50 bg-card/60 dark:bg-card/40 backdrop-blur-md">
          <CardContent class="p-16 flex flex-col items-center text-center space-y-3">
            <UsersRound class="w-14 h-14 text-muted-foreground/40" />
            <h3 class="text-sm font-bold text-foreground">Select a Team</h3>
            <p class="text-xs text-muted-foreground max-w-56">Choose a team from the list to view and manage its members.</p>
          </CardContent>
        </Card>
      </div>
    </div>

    <!-- Create Team Dialog -->
    <Dialog v-model:open="isCreateDialogOpen">
      <DialogContent class="sm:max-w-md bg-card border-border/50">
        <DialogHeader>
          <DialogTitle class="text-foreground font-bold tracking-tight">Create Team</DialogTitle>
          <DialogDescription class="text-muted-foreground text-xs">Create a new team to organize members and collaborate.</DialogDescription>
        </DialogHeader>
        <div class="space-y-4 py-4">
          <div class="space-y-1.5">
            <Label for="create-name" class="text-xs font-bold text-foreground/80">Team Name</Label>
            <Input id="create-name" v-model="formName" placeholder="e.g. Platform Engineering" class="h-10 rounded-lg" />
          </div>
        </div>
        <DialogFooter class="flex sm:justify-end gap-2">
          <Button variant="outline" @click="isCreateDialogOpen = false" :disabled="formLoading">Cancel</Button>
          <Button @click="handleCreate" :disabled="formLoading || !formName" class="gap-1.5">
            <Loader2 v-if="formLoading" class="w-4 h-4 animate-spin" />
            <Plus class="w-4 h-4" />
            <span>Create</span>
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Edit Team Dialog -->
    <Dialog v-model:open="isEditDialogOpen">
      <DialogContent class="sm:max-w-md bg-card border-border/50">
        <DialogHeader>
          <DialogTitle class="text-foreground font-bold tracking-tight">Edit Team</DialogTitle>
          <DialogDescription class="text-muted-foreground text-xs">Rename this team.</DialogDescription>
        </DialogHeader>
        <div class="space-y-4 py-4">
          <div class="space-y-1.5">
            <Label for="edit-name" class="text-xs font-bold text-foreground/80">Team Name</Label>
            <Input id="edit-name" v-model="formName" placeholder="e.g. Platform Engineering" class="h-10 rounded-lg" />
          </div>
        </div>
        <DialogFooter class="flex sm:justify-end gap-2">
          <Button variant="outline" @click="isEditDialogOpen = false" :disabled="formLoading">Cancel</Button>
          <Button @click="handleUpdate" :disabled="formLoading || !formName" class="gap-1.5">
            <Loader2 v-if="formLoading" class="w-4 h-4 animate-spin" />
            <span>Save</span>
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Delete Team Dialog -->
    <Dialog v-model:open="isDeleteDialogOpen">
      <DialogContent class="sm:max-w-md bg-card border-border/50">
        <DialogHeader>
          <DialogTitle class="text-foreground font-bold tracking-tight flex items-center gap-2 text-destructive">
            <ShieldAlert class="w-5 h-5" />
            <span>Delete Team</span>
          </DialogTitle>
          <DialogDescription class="text-muted-foreground text-xs leading-relaxed">
            Are you sure you want to delete <span class="font-bold text-foreground">"{{ actionTeam?.name }}"</span>? This action is permanent and will remove all memberships.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter class="flex sm:justify-end gap-2">
          <Button variant="outline" @click="isDeleteDialogOpen = false" :disabled="formLoading">Cancel</Button>
          <Button variant="destructive" @click="handleDeleteConfirm" :disabled="formLoading" class="gap-1.5">
            <Loader2 v-if="formLoading" class="w-4 h-4 animate-spin" />
            <span>Delete</span>
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Invite Member Dialog -->
    <Dialog v-model:open="isInviteDialogOpen">
      <DialogContent class="sm:max-w-md bg-card border-border/50">
        <DialogHeader>
          <DialogTitle class="text-foreground font-bold tracking-tight">Invite Member</DialogTitle>
          <DialogDescription class="text-muted-foreground text-xs">Invite a user to join <span class="font-bold text-foreground">{{ selectedTeam?.name }}</span>.</DialogDescription>
        </DialogHeader>
        <div class="space-y-4 py-4">
          <div class="space-y-1.5">
            <Label for="invite-user" class="text-xs font-bold text-foreground/80">Select User</Label>
            <select
              id="invite-user"
              v-model="inviteUserID"
              class="flex h-10 w-full rounded-lg border border-input bg-card px-3 py-2 text-sm text-foreground ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            >
              <option :value="0" disabled>Choose a user...</option>
              <option v-for="u in users" :key="u.id" :value="u.id">{{ u.name }} ({{ u.email }})</option>
            </select>
          </div>
          <div class="space-y-1.5">
            <Label for="invite-role" class="text-xs font-bold text-foreground/80">Role</Label>
            <select
              id="invite-role"
              v-model="inviteRole"
              class="flex h-10 w-full rounded-lg border border-input bg-card px-3 py-2 text-sm text-foreground ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            >
              <option value="member">Member</option>
              <option value="admin">Admin</option>
            </select>
          </div>
        </div>
        <DialogFooter class="flex sm:justify-end gap-2">
          <Button variant="outline" @click="isInviteDialogOpen = false" :disabled="formLoading">Cancel</Button>
          <Button @click="handleInvite" :disabled="formLoading || !inviteUserID" class="gap-1.5">
            <Loader2 v-if="formLoading" class="w-4 h-4 animate-spin" />
            <Send class="w-4 h-4" />
            <span>Send Invitation</span>
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
