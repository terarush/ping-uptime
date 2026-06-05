<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'

export interface DataPoint {
  date: string
  status: string
  uptime_pct: number
  failed: number
  total: number
  latency: number
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

// Only compute max latency from successful pings (or pings with latency > 0)
const validLatencies = computed(() => props.points.filter(p => p.status === 'up' && p.latency > 0).map(p => p.latency))
const maxLatency = computed(() => {
  const vals = validLatencies.value
  return vals.length > 0 ? Math.max(...vals, 50) : 50
})

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
  const c = ctx
  if (!c || !canvasRef.value) return
  const { width } = canvasRef.value.getBoundingClientRect()
  const height = props.height ?? 180
  c.clearRect(0, 0, width, height)

  if (!props.points.length) {
    c.fillStyle = '#94a3b8'
    c.font = '12px ui-sans-serif, system-ui, -apple-system, sans-serif'
    c.textAlign = 'center'
    c.fillText('No data', width / 2, height / 2)
    return
  }

  const points = props.points
  const count = points.length
  const padX = 48
  const padY = 20
  const plotW = width - padX - 24
  const plotH = height - padY * 2 - 20 // Reserve 20px at the bottom for status dots

  // Pre-process latencies: interpolate down/failed pings (latency <= 0) to avoid dropping the line chart to 0ms
  const cleanLatencies = points.map((p, idx) => {
    if (p.status === 'down' || p.latency <= 0) {
      // Find nearest successful ping backward
      for (let j = idx - 1; j >= 0; j--) {
        if (points[j]!.status === 'up' && points[j]!.latency > 0) {
          return points[j]!.latency
        }
      }
      // Find nearest successful ping forward
      for (let j = idx + 1; j < count; j++) {
        if (points[j]!.status === 'up' && points[j]!.latency > 0) {
          return points[j]!.latency
        }
      }
      return 0
    }
    return p.latency
  })

  const xFor = (i: number) => (count === 1 ? padX + plotW / 2 : padX + (i / (count - 1)) * plotW)
  const yForLatency = (i: number) => padY + plotH - ((cleanLatencies[i] ?? 0) / maxLatency.value) * plotH

  c.save()

  // Draw Y-axis grid lines and latency labels
  c.strokeStyle = 'rgba(148,163,184,0.08)'
  c.lineWidth = 1
  c.fillStyle = '#94a3b8'
  c.font = '9px ui-sans-serif, system-ui, sans-serif'
  c.textAlign = 'right'
  c.textBaseline = 'middle'

  const gridLines = 3
  for (let i = 0; i <= gridLines; i++) {
    const y = padY + (i / gridLines) * plotH
    const latencyVal = Math.round(maxLatency.value * (1 - i / gridLines))
    c.beginPath()
    c.moveTo(padX, y)
    c.lineTo(padX + plotW, y)
    c.stroke()
    c.fillText(`${latencyVal}ms`, padX - 8, y)
  }

  // Draw X-axis date labels
  c.textAlign = 'center'
  c.textBaseline = 'top'
  c.fillStyle = '#64748b'
  if (count > 1) {
    const labelStep = Math.max(Math.floor(count / 5), 1)
    for (let i = 0; i < count; i += labelStep) {
      const x = xFor(i)
      const dateStr = points[i]?.date || ''
      const parts = dateStr.split(' ')
      const disp = parts.length > 1 ? parts[1]!.substring(0, 5) : dateStr
      c.fillText(disp, x, padY + plotH + 20)
    }
  } else if (count === 1) {
    c.fillText(points[0]?.date || '', padX + plotW / 2, padY + plotH + 20)
  }

  // Draw Latency Area Fill
  if (count > 0) {
    c.beginPath()
    c.moveTo(xFor(0), padY + plotH)
    for (let i = 0; i < count; i++) {
      c.lineTo(xFor(i), yForLatency(i))
    }
    c.lineTo(xFor(count - 1), padY + plotH)
    c.closePath()

    const gradient = c.createLinearGradient(0, padY, 0, padY + plotH)
    gradient.addColorStop(0, 'rgba(99, 102, 241, 0.22)')
    gradient.addColorStop(1, 'rgba(99, 102, 241, 0.00)')
    c.fillStyle = gradient
    c.fill()
  }

  // Draw Latency Stroke Line
  c.beginPath()
  for (let i = 0; i < count; i++) {
    const x = xFor(i)
    const y = yForLatency(i)
    if (i === 0) c.moveTo(x, y)
    else c.lineTo(x, y)
  }
  c.strokeStyle = 'rgba(99, 102, 241, 0.9)'
  c.lineWidth = 2
  c.stroke()

  // Draw Status indicators in their dedicated row at the bottom
  points.forEach((p, i) => {
    const x = xFor(i)
    const cfg = STATUS_COLORS[p.status] || STATUS_COLORS.unknown || { bg: '', border: '' }
    c.fillStyle = cfg.bg
    c.strokeStyle = cfg.border
    c.lineWidth = 1.2
    c.beginPath()
    c.arc(x, padY + plotH + 8, 4.5, 0, Math.PI * 2)
    c.fill()
    c.stroke()
  })

  c.restore()
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
