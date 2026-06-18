/**
 * Composable for Team Management
 * delegates to teams Pinia store.
 */
import { storeToRefs } from 'pinia';
import { useTeamsStore, type Team, type TeamMember } from '@/stores/teams';

export type { Team, TeamMember };

export function useTeams() {
  const store = useTeamsStore();
  const { teams, members, loading, error } = storeToRefs(store);

  return {
    teams,
    members,
    loading,
    error,
    fetchAll: store.fetchAll,
    create: store.create,
    update: store.update,
    deleteTeam: store.deleteTeam,
    fetchMembers: store.fetchMembers,
    inviteMember: store.inviteMember,
    updateMember: store.updateMember,
    removeMember: store.removeMember,
    acceptInvite: store.acceptInvite,
    rejectInvite: store.rejectInvite,
  };
}
