import { spawnSync } from 'node:child_process'
import { mkdirSync } from 'node:fs'
import { join } from 'node:path'

const tmpDir = join('node_modules', '.tmp')
const storageFile = join(tmpDir, 'vitest-localstorage')
const vitestEntry = join('node_modules', 'vitest', 'vitest.mjs')

const env = { ...process.env }

// Vitest spawns its own worker processes, which don't inherit CLI flags
// passed to this wrapper. NODE_OPTIONS is read by every descendant Node
// process, so it's the only reliable way to silence the warning everywhere.
if (process.allowedNodeEnvironmentFlags.has('--localstorage-file')) {
  mkdirSync(tmpDir, { recursive: true })
  const flag = `--localstorage-file=${storageFile}`
  env.NODE_OPTIONS = env.NODE_OPTIONS ? `${env.NODE_OPTIONS} ${flag}` : flag
}

const result = spawnSync(process.execPath, [vitestEntry, ...process.argv.slice(2)], {
  stdio: 'inherit',
  env,
})

if (result.error) {
  throw result.error
}

process.exit(result.status ?? 1)
