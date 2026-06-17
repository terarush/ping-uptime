<script setup lang="ts">
import AppNavbar from '@/components/AppNavbar.vue'
import AppFooter from '@/components/AppFooter.vue'
import { useScrollReveal } from '@/lib/useScrollReveal'

const { observe } = useScrollReveal()

const steps = [
  {
    num: '01',
    title: 'Fork the Repository',
    desc: 'Head over to the ping-uptime repo on GitHub and click the Fork button. This creates a copy under your account.',
    code: "gh repo fork terarush/ping-uptime --clone",
  },
  {
    num: '02',
    title: 'Create a Feature Branch',
    desc: 'Always branch from dev. Use a descriptive name like feat/your-feature or fix/your-bugfix.',
    code: "git checkout dev\ngit checkout -b feat/my-contribution",
  },
  {
    num: '03',
    title: 'Make Your Changes',
    desc: 'Write clean, modular code. Follow the existing patterns — Go backend uses DDD modules, Vue frontend uses composables & Pinia stores.',
    code: "# Backend: go run main.go\n# Frontend: npm run dev --prefix web\n# Landing:  npm run dev --prefix landing",
  },
  {
    num: '04',
    title: 'Commit & Push',
    desc: 'Write clear, conventional commit messages. Then push to your fork.',
    code: "git add .\ngit commit -m \"feat: add description of change\"\ngit push origin feat/my-contribution",
  },
  {
    num: '05',
    title: 'Open a Pull Request',
    desc: 'Go to the original repo and open a PR from your branch to dev. Include a clear description of what you changed and why.',
    code: "gh pr create --base dev --head your-username:feat/my-contribution",
  },
]
</script>

<template>
  <div>
    <AppNavbar />
    <main>
      <!-- Hero -->
      <section class="relative pt-32 pb-20 md:pt-40 md:pb-28 overflow-hidden">
        <div class="absolute inset-0 grid-bg opacity-[0.03] dark:opacity-[0.06]" />
        <div class="relative mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 text-center">
          <div class="mx-auto max-w-3xl">
            <span
              class="inline-flex items-center gap-2 rounded-full border border-border bg-background/80 px-3.5 py-1 text-xs font-mono font-medium text-muted-foreground backdrop-blur-sm mb-5 tracking-wider"
            >
              <span class="h-1.5 w-1.5 rounded-full bg-primary/60" />
              Contribute
            </span>
            <h1 class="font-display text-4xl font-bold tracking-tight sm:text-5xl lg:text-6xl">
              Contribute to
              <span class="text-primary">ping-uptime</span>
            </h1>
            <p class="mt-4 text-lg text-muted-foreground max-w-xl mx-auto">
              Help us build the best self-hosted uptime monitor. All contributions — code, docs, issues — are welcome.
            </p>
            <div class="mt-8 flex items-center justify-center gap-4">
              <a
                href="https://github.com/terarush/ping-uptime"
                target="_blank"
                rel="noopener noreferrer"
                class="inline-flex items-center gap-2 rounded-xl bg-primary px-6 py-3 text-sm font-semibold text-primary-foreground shadow-lg shadow-primary/25 hover:shadow-xl hover:shadow-primary/30 transition-all"
              >
                <img src="/github.svg" class="h-4 w-4 invert" aria-hidden="true" />
                View Repository
              </a>
              <a
                href="https://github.com/terarush/ping-uptime/issues"
                target="_blank"
                rel="noopener noreferrer"
                class="inline-flex items-center gap-2 rounded-xl border border-border bg-background px-6 py-3 text-sm font-semibold text-foreground hover:bg-accent transition-colors"
              >
                Open Issues
              </a>
            </div>
          </div>
        </div>
      </section>

      <!-- How to contribute -->
      <section class="py-20 md:py-28 bg-muted/30">
        <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div class="mx-auto max-w-2xl text-center mb-16" :ref="(el: any) => observe(el)">
            <h2 class="font-display text-3xl font-bold tracking-tight sm:text-4xl">
              How to Contribute
            </h2>
            <p class="mt-3 text-lg text-muted-foreground">
              Five simple steps to get your first PR merged
            </p>
          </div>
          <div class="mx-auto max-w-3xl space-y-8">
            <div
              v-for="(step, i) in steps"
              :key="step.num"
              :ref="(el: any) => observe(el)"
              class="relative flex gap-6 rounded-2xl border border-border/50 bg-card/50 backdrop-blur-md p-6 sm:p-8 hover:shadow-md transition-shadow"
            >
              <!-- Number badge -->
              <div
                class="hidden sm:flex h-14 w-14 shrink-0 items-center justify-center rounded-xl bg-primary/10 text-primary font-mono text-sm font-bold"
              >
                {{ step.num }}
              </div>
              <div class="min-w-0 flex-1">
                <h3 class="font-display text-lg font-semibold mb-2">{{ step.title }}</h3>
                <p class="text-sm text-muted-foreground mb-4">{{ step.desc }}</p>
                <pre
                  class="overflow-x-auto rounded-lg bg-muted/80 p-4 text-xs font-mono text-foreground/80"
                ><code>{{ step.code }}</code></pre>
              </div>
              <!-- Connector line -->
              <div
                v-if="i < steps.length - 1"
                class="absolute left-[1.625rem] top-[5.5rem] bottom-[-2rem] w-px bg-border hidden sm:block"
              />
            </div>
          </div>
        </div>
      </section>

      <!-- Development setup -->
      <section class="py-20 md:py-28">
        <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div class="mx-auto max-w-2xl text-center mb-16" :ref="(el: any) => observe(el)">
            <span
              class="inline-flex items-center gap-2 rounded-full border border-border bg-background/80 px-3.5 py-1 text-xs font-mono font-medium text-muted-foreground backdrop-blur-sm mb-5 tracking-wider"
            >
              <span class="h-1.5 w-1.5 rounded-full bg-primary/60" />
              Setup
            </span>
            <h2 class="font-display text-3xl font-bold tracking-tight sm:text-4xl">
              Development Setup
            </h2>
          </div>
          <div class="mx-auto max-w-3xl space-y-6" :ref="(el: any) => observe(el)">
            <div class="rounded-2xl border border-border/50 bg-card/50 backdrop-blur-md p-6 sm:p-8">
              <h3 class="font-display text-base font-semibold mb-4">Prerequisites</h3>
              <ul class="space-y-2 text-sm text-muted-foreground">
                <li class="flex items-start gap-2">
                  <span class="mt-0.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary/60" />
                  Go 1.22+ — download from golang.org
                </li>
                <li class="flex items-start gap-2">
                  <span class="mt-0.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary/60" />
                  Node.js 20+ — download from nodejs.org
                </li>
                <li class="flex items-start gap-2">
                  <span class="mt-0.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary/60" />
                  Git — for version control
                </li>
              </ul>
            </div>
            <div class="rounded-2xl border border-border/50 bg-card/50 backdrop-blur-md p-6 sm:p-8">
              <h3 class="font-display text-base font-semibold mb-4">Run Locally</h3>
              <pre
                class="overflow-x-auto rounded-lg bg-muted/80 p-4 text-xs font-mono text-foreground/80"
              ><code># Backend
go run main.go

# Web frontend (separate terminal)
npm run dev --prefix web

# Landing page (separate terminal)
npm run dev --prefix landing</code></pre>
            </div>
            <div class="rounded-2xl border border-border/50 bg-card/50 backdrop-blur-md p-6 sm:p-8">
              <h3 class="font-display text-base font-semibold mb-4">Build &amp; Test</h3>
              <pre
                class="overflow-x-auto rounded-lg bg-muted/80 p-4 text-xs font-mono text-foreground/80"
              ><code># Full build
npm run build --prefix web && go build -o tmp/main main.go

# Backend tests
go test ./...

# Frontend type-check
npm run type-check --prefix web

# Frontend unit tests
npm run test:unit --prefix web</code></pre>
            </div>
          </div>
        </div>
      </section>

      <!-- Guidelines -->
      <section class="py-20 md:py-28 bg-muted/30">
        <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div class="mx-auto max-w-2xl text-center mb-16" :ref="(el: any) => observe(el)">
            <span
              class="inline-flex items-center gap-2 rounded-full border border-border bg-background/80 px-3.5 py-1 text-xs font-mono font-medium text-muted-foreground backdrop-blur-sm mb-5 tracking-wider"
            >
              <span class="h-1.5 w-1.5 rounded-full bg-primary/60" />
              Guidelines
            </span>
            <h2 class="font-display text-3xl font-bold tracking-tight sm:text-4xl">
              Code Guidelines
            </h2>
          </div>
          <div class="mx-auto max-w-3xl space-y-6" :ref="(el: any) => observe(el)">
            <div class="rounded-2xl border border-border/50 bg-card/50 backdrop-blur-md p-6 sm:p-8">
              <h3 class="font-display text-base font-semibold mb-3">Go Backend</h3>
              <ul class="space-y-2 text-sm text-muted-foreground">
                <li class="flex items-start gap-2">
                  <span class="mt-0.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary/60" />
                  Modular layout under <code class="text-xs font-mono bg-muted px-1 py-0.5 rounded">/modules/&lt;name&gt;</code> with interfaces for handlers, services, and repositories
                </li>
                <li class="flex items-start gap-2">
                  <span class="mt-0.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary/60" />
                  Group routes and apply auth middleware at the group level
                </li>
                <li class="flex items-start gap-2">
                  <span class="mt-0.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary/60" />
                  Use standard JSON response helpers (<code class="text-xs font-mono bg-muted px-1 py-0.5 rounded">SuccessResponse</code>, <code class="text-xs font-mono bg-muted px-1 py-0.5 rounded">ErrorResponse</code>)
                </li>
                <li class="flex items-start gap-2">
                  <span class="mt-0.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary/60" />
                  Keep it clean — <code class="text-xs font-mono bg-muted px-1 py-0.5 rounded">go fmt</code> before committing
                </li>
              </ul>
            </div>
            <div class="rounded-2xl border border-border/50 bg-card/50 backdrop-blur-md p-6 sm:p-8">
              <h3 class="font-display text-base font-semibold mb-3">Vue Frontend</h3>
              <ul class="space-y-2 text-sm text-muted-foreground">
                <li class="flex items-start gap-2">
                  <span class="mt-0.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary/60" />
                  Use Tailwind CSS v4 + shadcn-vue components — avoid direct styling when predefined tokens exist
                </li>
                <li class="flex items-start gap-2">
                  <span class="mt-0.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary/60" />
                  Views go flattened under <code class="text-xs font-mono bg-muted px-1 py-0.5 rounded">@/views</code> — no subfolders
                </li>
                <li class="flex items-start gap-2">
                  <span class="mt-0.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary/60" />
                  State management via Pinia stores wrapped in reactive composables
                </li>
                <li class="flex items-start gap-2">
                  <span class="mt-0.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary/60" />
                  Form validation with vee-validate + Zod schemas under <code class="text-xs font-mono bg-muted px-1 py-0.5 rounded">@/validations</code>
                </li>
              </ul>
            </div>
            <div class="rounded-2xl border border-border/50 bg-card/50 backdrop-blur-md p-6 sm:p-8">
              <h3 class="font-display text-base font-semibold mb-3">Commit Messages</h3>
              <p class="text-sm text-muted-foreground mb-3">Use conventional commits for clarity:</p>
              <pre
                class="overflow-x-auto rounded-lg bg-muted/80 p-4 text-xs font-mono text-foreground/80"
              ><code>feat:     new feature
fix:      bug fix
docs:     documentation only
style:    formatting, missing semicolons
refactor: code change that neither fixes nor adds
perf:     performance improvement
chore:    build process, tooling</code></pre>
            </div>
          </div>
        </div>
      </section>

      <!-- Reporting issues -->
      <section class="py-20 md:py-28">
        <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div class="mx-auto max-w-2xl text-center mb-16" :ref="(el: any) => observe(el)">
            <span
              class="inline-flex items-center gap-2 rounded-full border border-border bg-background/80 px-3.5 py-1 text-xs font-mono font-medium text-muted-foreground backdrop-blur-sm mb-5 tracking-wider"
            >
              <span class="h-1.5 w-1.5 rounded-full bg-primary/60" />
              Issues
            </span>
            <h2 class="font-display text-3xl font-bold tracking-tight sm:text-4xl">
              Reporting Issues
            </h2>
            <p class="mt-3 text-lg text-muted-foreground">
              Found a bug? Have a feature request? Open an issue on GitHub.
            </p>
          </div>
          <div class="mx-auto max-w-3xl grid gap-4 sm:grid-cols-2" :ref="(el: any) => observe(el)">
            <a
              href="https://github.com/terarush/ping-uptime/issues/new?template=bug_report.md"
              target="_blank"
              rel="noopener noreferrer"
              class="rounded-2xl border border-border/50 bg-card/50 backdrop-blur-md p-6 hover:shadow-md hover:border-destructive/30 transition-all group"
            >
              <span class="inline-flex h-10 w-10 items-center justify-center rounded-lg bg-destructive/10 text-destructive mb-4">!</span>
              <h3 class="font-display text-base font-semibold mb-1">Bug Report</h3>
              <p class="text-sm text-muted-foreground">Something not working right? Let us know.</p>
            </a>
            <a
              href="https://github.com/terarush/ping-uptime/issues/new?template=feature_request.md"
              target="_blank"
              rel="noopener noreferrer"
              class="rounded-2xl border border-border/50 bg-card/50 backdrop-blur-md p-6 hover:shadow-md hover:border-primary/30 transition-all group"
            >
              <span class="inline-flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary mb-4">+</span>
              <h3 class="font-display text-base font-semibold mb-1">Feature Request</h3>
              <p class="text-sm text-muted-foreground">Have an idea? We'd love to hear it.</p>
            </a>
          </div>
        </div>
      </section>

      <!-- Get in touch -->
      <section class="py-20 md:py-28 bg-muted/30">
        <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div class="mx-auto max-w-2xl text-center" :ref="(el: any) => observe(el)">
            <h2 class="font-display text-3xl font-bold tracking-tight sm:text-4xl">
              Questions?
            </h2>
            <p class="mt-3 text-lg text-muted-foreground">
              Join our community or reach out directly.
            </p>
            <div class="mt-8 flex items-center justify-center gap-4">
              <a
                href="https://discord.gg/G32T9q55xh"
                target="_blank"
                rel="noopener noreferrer"
                class="inline-flex items-center gap-2 rounded-xl bg-indigo-600 px-6 py-3 text-sm font-semibold text-white shadow-lg hover:bg-indigo-500 transition-colors"
              >
                Join Discord
              </a>
              <a
                href="https://github.com/terarush/ping-uptime/discussions"
                target="_blank"
                rel="noopener noreferrer"
                class="inline-flex items-center gap-2 rounded-xl border border-border bg-background px-6 py-3 text-sm font-semibold text-foreground hover:bg-accent transition-colors"
              >
                GitHub Discussions
              </a>
            </div>
          </div>
        </div>
      </section>
    </main>
    <AppFooter />
  </div>
</template>
