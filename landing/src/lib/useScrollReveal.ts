import { ref, onMounted, onUnmounted } from 'vue'

export function useScrollReveal(threshold = 0.1) {
  const els = ref<(HTMLElement | null)[]>([])
  let observer: IntersectionObserver | null = null
  let reduceMotion = false

  function observe(el: HTMLElement | null) {
    if (el && !els.value.includes(el)) {
      els.value.push(el)
      if (el) el.classList.toggle('visible', reduceMotion)
      if (observer && !reduceMotion) observer.observe(el)
    }
  }

  onMounted(() => {
    reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    if (reduceMotion) {
      els.value.forEach((el) => el?.classList.add('visible'))
      return
    }
    observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            entry.target.classList.add('visible')
            observer?.unobserve(entry.target)
          }
        })
      },
      { threshold }
    )
    els.value.forEach((el) => {
      if (el) observer?.observe(el)
    })
  })

  onUnmounted(() => {
    observer?.disconnect()
  })

  return { observe }
}
