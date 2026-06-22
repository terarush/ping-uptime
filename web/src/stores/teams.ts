import { defineStore } from 'pinia';
import { ref } from 'vue';
import ExtendedFetch from '@/lib/fetch';

export interface Team {
  id: number;
  name: string;
  member_count?: number;
  created_at: string;
}

export interface TeamMember {
  id: number;
  team_id: number;
  user_id: number;
  role: string;
  status: string;
  created_at: string;
  user?: { id: number; email: string; name?: string };
}

export const useTeamsStore = defineStore('teams', () => {
  const teams = ref<Team[]>([]);
  const members = ref<TeamMember[]>([]);
  const loading = ref(false);
  const error = ref('');

  const fetchAll = async () => {
    loading.value = true;
    error.value = '';
    try {
      const res = await ExtendedFetch.get('/teams');
      teams.value = res.data?.data || res.data || [];
    } catch (err: any) {
      console.error('Failed to fetch teams:', err);
      error.value = err.response?.data?.error || 'Failed to load teams.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const create = async (name: string) => {
    error.value = '';
    try {
      const res = await ExtendedFetch.post('/teams', { name });
      const team = res.data?.data || res.data;
      teams.value.push(team);
      return team;
    } catch (err: any) {
      console.error('Failed to create team:', err);
      error.value = err.response?.data?.error || 'Failed to create team.';
      throw err;
    }
  };

  const update = async (id: number, name: string) => {
    error.value = '';
    try {
      const res = await ExtendedFetch.put(`/teams/${id}`, { name });
      const updated = res.data?.data || res.data;
      const idx = teams.value.findIndex(t => t.id === id);
      if (idx !== -1) teams.value[idx] = updated;
      return updated;
    } catch (err: any) {
      console.error('Failed to update team:', err);
      error.value = err.response?.data?.error || 'Failed to update team.';
      throw err;
    }
  };

  const deleteTeam = async (id: number) => {
    error.value = '';
    try {
      await ExtendedFetch.delete(`/teams/${id}`);
      teams.value = teams.value.filter(t => t.id !== id);
    } catch (err: any) {
      console.error('Failed to delete team:', err);
      error.value = err.response?.data?.error || 'Failed to delete team.';
      throw err;
    }
  };

  const fetchMembers = async (teamID: number) => {
    loading.value = true;
    error.value = '';
    try {
      const res = await ExtendedFetch.get(`/teams/${teamID}/members`);
      members.value = res.data?.data || res.data || [];
    } catch (err: any) {
      console.error('Failed to fetch members:', err);
      error.value = err.response?.data?.error || 'Failed to load members.';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const inviteMember = async (teamID: number, userID: number, role: string = 'member') => {
    error.value = '';
    try {
      const res = await ExtendedFetch.post(`/teams/${teamID}/members`, { user_id: userID, role });
      const member = res.data?.data || res.data;
      members.value.push(member);
      return member;
    } catch (err: any) {
      console.error('Failed to invite member:', err);
      error.value = err.response?.data?.error || 'Failed to invite member.';
      throw err;
    }
  };

  const updateMember = async (teamID: number, memberID: number, role: string) => {
    error.value = '';
    try {
      const res = await ExtendedFetch.put(`/teams/${teamID}/members/${memberID}`, { role });
      const updated = res.data?.data || res.data;
      const idx = members.value.findIndex(m => m.id === memberID);
      if (idx !== -1) members.value[idx] = updated;
      return updated;
    } catch (err: any) {
      console.error('Failed to update member:', err);
      error.value = err.response?.data?.error || 'Failed to update member.';
      throw err;
    }
  };

  const removeMember = async (teamID: number, memberID: number) => {
    error.value = '';
    try {
      await ExtendedFetch.delete(`/teams/${teamID}/members/${memberID}`);
      members.value = members.value.filter(m => m.id !== memberID);
    } catch (err: any) {
      console.error('Failed to remove member:', err);
      error.value = err.response?.data?.error || 'Failed to remove member.';
      throw err;
    }
  };

  const acceptInvite = async (teamID: number) => {
    error.value = '';
    try {
      await ExtendedFetch.post(`/teams/${teamID}/members/accept`);
    } catch (err: any) {
      console.error('Failed to accept invitation:', err);
      error.value = err.response?.data?.error || 'Failed to accept invitation.';
      throw err;
    }
  };

  const rejectInvite = async (teamID: number) => {
    error.value = '';
    try {
      await ExtendedFetch.post(`/teams/${teamID}/members/reject`);
    } catch (err: any) {
      console.error('Failed to reject invitation:', err);
      error.value = err.response?.data?.error || 'Failed to reject invitation.';
      throw err;
    }
  };

  return {
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
  };
});
