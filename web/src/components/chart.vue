<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'

export interface DataPoint {
  date: string
  status: string
  uptime_pct: number
  failed: number
  total: number
}

interface StatusColors {
  [key: string]: { bg: string; border: string }
}

const STATUS_COLORS: StatusColors = {
  up: { bg: 'rgba(16,185,129,0.18)', border: 'rgba(16,185,129,0.9)' },
  down: { bg: 'rgba(239,68,68,0.22)', border: 'rgba(239,68,68,0.95)' },
  unknown: { bg: 'rgba(148,163,184,0.15)', border: 'rgba(148,163,184,0.8)' },
}

const props = withDefaults(defineProps<{ points: DataPoint[]; height?: number }>(), {
  height: 180,
})

const canvasRef = ref<HTMLCanvasElement | null>(null)
let ctx: CanvasRenderingContext2D | null = null
let dpr = 1

const maxTotal = computed(() => Math.max(...props.points.map(p => p.total), 1))
const maxFailed = computed(() => Math.max(...props.points.map(p => p.failed), 1))

function setupCanvas() {
  const el = canvasRef.value
  if (!el) return
  dpr = typeof window !== 'undefined' ? window.devicePixelRatio || 1 : 1
  const { width, height } = el.getBoundingClientRect()
  el.width = width * dpr
  el.height = (props.height ?? 180) * dpr
  ctx = el.getContext('2d')
  if (!ctx) return
  ctx.scale(dpr, dpr)
}

function draw() {
  if (!ctx || !canvasRef.value) return
  const { width } = canvasRef.value.getBoundingClientRect()
  const height = props.height ?? 180
  ctx.clearRect(0, 0, width, height)

  if (!props.points.length) {
    ctx.fillStyle = '#94a3b8'
    ctx.font = '12px ui-sans-serif, system-ui, -apple-system, sans-serif'
    ctx.textAlign = 'center'
    ctx.fillText('No data', width / 2, height / 2)
    return
  }

  const points = props.points
  const count = points.length
  const padX = 36
  const padY = 24
  const plotW = width - padX * 2
  const plotH = height - padY * 2

  const xFor = (i: number) => (count === 1 ? width / 2 : padX + (i / (count - 1)) * plotW)
  const yForTotal = (i: number) => padY + plotH - (points[i].total / maxTotal.value) * plotH
  const yForFailed = (i: number) => padY + plotH - (points[i].failed / maxTotal.value) * plotH

  ctx.save()

  ['total', 'failed'].forEach((key, order) => {
    const yFn = key === 'total' ? yForTotal : yForFailed
    const color = key === 'total' ? 'rgba(99,102,241,0.35)' : 'rgba(239,68,68,0.55)'
    const lineColor = key === 'total' ? 'rgba(99,102,241,0.75)' : 'rgba(239,68,68,0.85)'
    const width2 = key === 'total' ? 1.5 : 1

    ctx.beginPath()
    for (let i = 0; i < count; i++) {
      const x = xFor(i)
      const y = yFn(i)
      if (i === 0) ctx.moveTo(x, y)
      else ctx.lineTo(x, y)
    }
    ctx.strokeStyle = lineColor
    ctx.lineWidth = width2
    ctx.stroke()

    ctx.beginPath()
    for (let i = 0; i < count; i++) {
      const x = xFor(i)
      const y = yFn(i)
      ctx.rect(x - 2, y - 2, order === 0 ? 3 : 2, order === 0 ? 3 : 2)
    }
    ctx.fillStyle = color
    ctx.fill()
  })

  points.forEach((p, i) => {
    const x = xFor(i)
    const cfg = STATUS_COLORS[p.status] || STATUS_COLORS.unknown
    ctx.fillStyle = cfg.bg
    ctx.strokeStyle = cfg.border
    ctx.lineWidth = 1
    ctx.beginPath()
    ctx.arc(x, padY + plotH - 3, 4.5, 0, Math.PI * 2)
    ctx.fill()
    ctx.stroke()
  })

  ctx.restore()
}

let raf = 0
function schedule() {
  cancelAnimationFrame(raf)
  raf = requestAnimationFrame(() => {
    setupCanvas()
    draw()
  })
}

watch(() => props.points, schedule, { deep: true })
watch(() => props.height, schedule)

onMounted(() => {
  schedule()
  window.addEventListener('resize', schedule)
})
onUnmounted(() => {
  cancelAnimationFrame(raf)
  window.removeEventListener('resize', schedule)
})
</script>

<template>
  <div class="w-full" :style="{ height: `${height ?? 180}px` }">
    <canvas ref="canvasRef" class="w-full h-full block" />
  </div>
</template>
