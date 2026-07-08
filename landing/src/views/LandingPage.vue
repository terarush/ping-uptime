<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import AppNavbar from '@/components/AppNavbar.vue'
import HeroSection from '@/components/HeroSection.vue'
import BenefitsSection from '@/components/BenefitsSection.vue'
import FeaturesSection from '@/components/FeaturesSection.vue'
import ServicesSection from '@/components/ServicesSection.vue'
import TeamSection from '@/components/TeamSection.vue'
import CtaSection from '@/components/CtaSection.vue'
import ContactSection from '@/components/ContactSection.vue'
import AppFooter from '@/components/AppFooter.vue'

const route = useRoute()

onMounted(() => {
  if (route.hash) {
    const id = route.hash.slice(1)
    setTimeout(() => {
      const el = document.getElementById(id)
      if (el) el.scrollIntoView({ behavior: 'smooth', block: 'start' })
    }, 100)
  }
})
</script>

<template>
  <div class="min-h-screen bg-background">
    <AppNavbar />
    <main>
      <HeroSection />
      <BenefitsSection />
      <FeaturesSection />
      <ServicesSection />
      <TeamSection />
      <CtaSection />
      <ContactSection />
    </main>
    <AppFooter />
  </div>
</template>
