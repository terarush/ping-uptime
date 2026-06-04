/**
 * Composable for User Management
 * delegates to users Pinia store.
 */
import { storeToRefs } from 'pinia';
import { useUsersStore, type User } from '@/stores/users';

export type { User };

export function useUsers() {
  const store = useUsersStore();
  const { users, loading, error } = storeToRefs(store);

  return {
    users,
    loading,
    error,
    fetchUsersData: store.fetchUsersData,
    createUser: store.createUser,
    updateUser: store.updateUser,
    deleteUser: store.deleteUser,
  };
}
