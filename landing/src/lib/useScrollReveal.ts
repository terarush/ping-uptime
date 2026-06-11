import { ref, onMounted, onUnmounted } from 'vue'

export function useScrollReveal(threshold = 0.1) {
  const els = ref<(HTMLElement | null)[]>([])
  let observer: IntersectionObserver | null = null

  function observe(el: HTMLElement | null) {
    if (el && !els.value.includes(el)) {
      els.value.push(el)
      if (observer) observer.observe(el)
    }
  }

  onMounted(() => {
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
