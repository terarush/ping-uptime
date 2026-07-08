import { ref, onMounted, onUnmounted } from 'vue'
import { siteConfig } from '@/content/config'

interface GitHubRelease {
  tag_name: string
  html_url: string
  published_at: string
  prerelease: boolean
}

const latestRelease = ref<GitHubRelease | null>(null)
const isUpdateAvailable = ref(false)
const loading = ref(false)
const error = ref(false)
let intervalId: ReturnType<typeof setInterval> | null = null

function compareVersions(a: string, b: string): number {
  const aClean = a.startsWith('v') ? a.slice(1) : a
  const bClean = b.startsWith('v') ? b.slice(1) : b
  const aParts = aClean.split('.').map(Number)
  const bParts = bClean.split('.').map(Number)
  for (let i = 0; i < Math.max(aParts.length, bParts.length); i++) {
    const an = aParts[i] || 0
    const bn = bParts[i] || 0
    if (an > bn) return 1
    if (an < bn) return -1
  }
  return 0
}

async function checkForUpdate() {
  if (loading.value) return
  loading.value = true
  error.value = false
  try {
    const res = await fetch(
      `https://api.github.com/repos/${siteConfig.repoOwner}/${siteConfig.repoName}/releases/latest`,
      { headers: { Accept: 'application/vnd.github.v3+json' } }
    )
    if (!res.ok) throw new Error('Failed to fetch')
    const data: GitHubRelease = await res.json()
    latestRelease.value = data
    // Skip if current version is unknown (dev build without git)
    const current = __APP_VERSION__
    if (current === '0.0.0') return
    isUpdateAvailable.value =
      compareVersions(data.tag_name, current) > 0
  } catch {
    error.value = true
  } finally {
    loading.value = false
  }
}

export function useUpdateChecker() {
  onMounted(() => {
    checkForUpdate()
    // Check every hour
    intervalId = setInterval(checkForUpdate, 3600_000)
  })

  onUnmounted(() => {
    if (intervalId) {
      clearInterval(intervalId)
      intervalId = null
    }
  })

  return {
    latestRelease,
    isUpdateAvailable,
    loading,
    error,
    checkForUpdate,
  }
}
