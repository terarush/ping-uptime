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

export interface Series {
  monitorId: number
  name: string
  color: string
  points: DataPoint[]
}

const props = withDefaults(defineProps<{
  seriesList: Series[]
  height?: number
}>(), {
  height: 300,
})

const canvasRef = ref<HTMLCanvasElement | null>(null)
const containerRef = ref<HTMLDivElement | null>(null)
let ctx: CanvasRenderingContext2D | null = null
let dpr = 1

// Tooltip state
const tooltipVisible = ref(false)
const tooltipX = ref(0)
const tooltipY = ref(0)
const tooltipTitle = ref('')
const tooltipItems = ref<{ name: string; color: string; latency: number; status: string }[]>([])

// Create a union of all dates, sorted chronologically
const allDates = computed(() => {
  const dates = new Set<string>()
  props.seriesList.forEach(s => {
    s.points.forEach(p => {
      if (p.date) dates.add(p.date)
    })
  })
  return Array.from(dates).sort()
})

// Calculate the maximum latency across all active series to scale Y-axis
const maxLatency = computed(() => {
  let max = 50
  props.seriesList.forEach(s => {
    s.points.forEach(p => {
      if (p.status === 'up' && p.latency > max) {
        max = p.latency
      }
    })
  })
  return max * 1.1 // Add 10% padding
})

function setupCanvas() {
  const el = canvasRef.value
  if (!el) return
  dpr = typeof window !== 'undefined' ? window.devicePixelRatio || 1 : 1
  const { width } = el.getBoundingClientRect()
  el.width = width * dpr
  el.height = props.height * dpr
  ctx = el.getContext('2d')
  if (!ctx) return
  ctx.scale(dpr, dpr)
}

function draw() {
  const c = ctx
  if (!c || !canvasRef.value) return
  const { width } = canvasRef.value.getBoundingClientRect()
  const height = props.height
  c.clearRect(0, 0, width, height)

  if (props.seriesList.length === 0 || allDates.value.length === 0) {
    c.fillStyle = '#94a3b8'
    c.font = '13px ui-sans-serif, system-ui, sans-serif'
    c.textAlign = 'center'
    c.fillText('No data available', width / 2, height / 2)
    return
  }

  const dates = allDates.value
  const count = dates.length
  const padX = 55
  const padY = 20
  const plotW = width - padX - 25
  const plotH = height - padY * 2 - 10

  const xForIndex = (i: number) => (count === 1 ? padX + plotW / 2 : padX + (i / (count - 1)) * plotW)
  const yForLatency = (latency: number) => padY + plotH - (latency / maxLatency.value) * plotH

  c.save()

  // Draw Y-axis grid lines and labels
  c.strokeStyle = 'rgba(148,163,184,0.08)'
  c.lineWidth = 1
  c.fillStyle = '#94a3b8'
  c.font = '10px ui-sans-serif, system-ui, sans-serif'
  c.textAlign = 'right'
  c.textBaseline = 'middle'

  const gridLines = 4
  for (let i = 0; i <= gridLines; i++) {
    const y = padY + (i / gridLines) * plotH
    const latencyVal = Math.round(maxLatency.value * (1 - i / gridLines))
    c.beginPath()
    c.moveTo(padX, y)
    c.lineTo(padX + plotW, y)
    c.stroke()
    c.fillText(`${latencyVal}ms`, padX - 10, y)
  }

  // Draw X-axis date labels
  c.textAlign = 'center'
  c.textBaseline = 'top'
  c.fillStyle = '#64748b'
  
  if (count > 1) {
    const labelStep = Math.max(Math.floor(count / 6), 1)
    for (let i = 0; i < count; i += labelStep) {
      const x = xForIndex(i)
      const dateStr = dates[i]
      if (dateStr) {
        const parts = dateStr.split(' ')
        // For hourly view show hour, for daily view show date
        const disp = parts.length > 1 ? (parts[1]?.substring(0, 5) ?? dateStr) : dateStr
        c.fillText(disp, x, padY + plotH + 12)
      }
    }
  } else if (count === 1) {
    c.fillText(dates[0] || '', padX + plotW / 2, padY + plotH + 12)
  }

  // Draw each series line
  props.seriesList.forEach(series => {
    // Map points to absolute date indices
    const pointMap = new Map<string, DataPoint>()
    series.points.forEach(p => {
      if (p.date) pointMap.set(p.date, p)
    })

    c.beginPath()
    let first = true
    
    // Draw continuous lines
    for (let i = 0; i < count; i++) {
      const date = dates[i]
      if (date) {
        const p = pointMap.get(date)
        if (p) {
          const x = xForIndex(i)
          const latency = p.status === 'down' || p.latency <= 0 ? 0 : p.latency
          const y = yForLatency(latency)
          
          if (first) {
            c.moveTo(x, y)
            first = false
          } else {
            c.lineTo(x, y)
          }
        }
      }
    }

    c.strokeStyle = series.color
    c.lineWidth = 2
    c.lineJoin = 'round'
    c.lineCap = 'round'
    c.stroke()

    // Draw little circles for data points if count is small (e.g. less than 30 points)
    if (count < 30) {
      for (let i = 0; i < count; i++) {
        const date = dates[i]
        if (date) {
          const p = pointMap.get(date)
          if (p) {
            const x = xForIndex(i)
            const latency = p.status === 'down' || p.latency <= 0 ? 0 : p.latency
            const y = yForLatency(latency)
            
            c.beginPath()
            c.arc(x, y, 3, 0, Math.PI * 2)
            c.fillStyle = series.color
            c.fill()
            c.strokeStyle = '#ffffff'
            c.lineWidth = 1
            c.stroke()
          }
        }
      }
    }
  })

  c.restore()
}

// Mouse interaction for tooltip
function handleMouseMove(e: MouseEvent) {
  if (!canvasRef.value || props.seriesList.length === 0 || allDates.value.length === 0) return

  const rect = canvasRef.value.getBoundingClientRect()
  const mouseX = e.clientX - rect.left
  const dates = allDates.value
  const count = dates.length
  
  const padX = 55
  const padY = 20
  const plotW = rect.width - padX - 25

  if (mouseX < padX || mouseX > padX + plotW) {
    tooltipVisible.value = false
    return
  }

  // Find the closest date index
  const percent = (mouseX - padX) / plotW
  const idx = Math.min(Math.max(Math.round(percent * (count - 1)), 0), count - 1)
  const targetDate = dates[idx]

  if (!targetDate) {
    tooltipVisible.value = false
    return
  }

  // Gather values from all series at this index
  const items: typeof tooltipItems.value = []
  
  props.seriesList.forEach(series => {
    const pt = series.points.find(p => p.date === targetDate)
    if (pt) {
      items.push({
        name: series.name,
        color: series.color,
        latency: Math.round(pt.latency),
        status: pt.status
      })
    }
  })

  // Sort tooltip items by latency descending
  items.sort((a, b) => b.latency - a.latency)

  // Calculate coordinates for tooltip
  const x = padX + (idx / (count - 1)) * plotW
  
  tooltipVisible.value = true
  tooltipX.value = x
  tooltipY.value = padY
  tooltipTitle.value = targetDate
  tooltipItems.value = items
}

function handleMouseLeave() {
  tooltipVisible.value = false
}

let raf = 0
function schedule() {
  cancelAnimationFrame(raf)
  raf = requestAnimationFrame(() => {
    setupCanvas()
    draw()
  })
}

watch(() => props.seriesList, schedule, { deep: true })
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
  <div ref="containerRef" class="w-full relative" :style="{ height: `${height}px` }">
    <canvas
      ref="canvasRef"
      class="w-full h-full block cursor-crosshair"
      @mousemove="handleMouseMove"
      @mouseleave="handleMouseLeave"
    />

    <!-- Vertical Marker Line -->
    <div
      v-if="tooltipVisible"
      class="absolute top-[20px] bottom-[40px] w-px bg-slate-300 dark:bg-slate-700 pointer-events-none transition-all duration-75"
      :style="{ left: `${tooltipX}px` }"
    />

    <!-- Hover Tooltip Box -->
    <div
      v-if="tooltipVisible"
      class="absolute z-50 bg-popover/95 border border-border text-popover-foreground rounded-lg shadow-xl p-3 text-xs pointer-events-none min-w-[200px] backdrop-blur-md transition-all duration-75"
      :style="{
        left: `${tooltipX + 15 + 200 > (containerRef?.clientWidth ?? 0) ? tooltipX - 215 : tooltipX + 15}px`,
        top: `${tooltipY}px`
      }"
    >
      <div class="font-bold border-b border-border/60 pb-1.5 mb-1.5 text-foreground flex items-center justify-between">
        <span>Timeline Details</span>
        <span class="text-[10px] text-muted-foreground font-medium">{{ tooltipTitle }}</span>
      </div>
      <div class="space-y-1.5">
        <div
          v-for="item in tooltipItems"
          :key="item.name"
          class="flex items-center justify-between gap-4"
        >
          <div class="flex items-center gap-1.5 truncate">
            <span class="w-2 h-2 rounded-full shrink-0" :style="{ backgroundColor: item.color }" />
            <span class="font-medium truncate text-[11px]">{{ item.name }}</span>
          </div>
          <div class="flex items-center gap-1 shrink-0">
            <span class="font-bold text-[11px]">{{ item.latency }} ms</span>
            <span
              class="w-1.5 h-1.5 rounded-full"
              :class="[item.status === 'up' ? 'bg-emerald-500' : 'bg-rose-500']"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
