import { spawnSync } from 'node:child_process'
import { mkdirSync } from 'node:fs'
import { join } from 'node:path'

const tmpDir = join('node_modules', '.tmp')
const storageFile = join(tmpDir, 'vitest-localstorage')
const vitestEntry = join('node_modules', 'vitest', 'vitest.mjs')

const nodeArgs = []

if (process.allowedNodeEnvironmentFlags.has('--localstorage-file')) {
  mkdirSync(tmpDir, { recursive: true })
  nodeArgs.push(`--localstorage-file=${storageFile}`)
}

const result = spawnSync(process.execPath, [...nodeArgs, vitestEntry, ...process.argv.slice(2)], {
  stdio: 'inherit',
})

if (result.error) {
  throw result.error
}

process.exit(result.status ?? 1)
