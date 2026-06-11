/**
 * Singleton composable — all consumers share same reactive appTitle ref.
 * Any component can call fetchAppTitle() and all watchers update.
 */
import { ref } from 'vue';
import { siteConfig } from '@/content/config';

const appTitle = ref(siteConfig.name);
const fetchAppTitle = async () => {
  try {
    const res = await fetch('/api/settings/public/system-name');
    if (res.ok) {
      const json = await res.json();
      const val = json?.data?.value;
      if (val) appTitle.value = val;
    }
  } catch {
    // silent — keep default
  }
};

export function useAppTitle() {
  return { appTitle, fetchAppTitle };
}
