<script setup lang="ts">
import AppNavbar from '@/components/AppNavbar.vue'
import HeroSection from '@/components/HeroSection.vue'
import BenefitsSection from '@/components/BenefitsSection.vue'
import FeaturesSection from '@/components/FeaturesSection.vue'
import ServicesSection from '@/components/ServicesSection.vue'
import TeamSection from '@/components/TeamSection.vue'
import ContactSection from '@/components/ContactSection.vue'
import AppFooter from '@/components/AppFooter.vue'
</script>

<template>
  <div class="min-h-screen bg-background">
    <AppNavbar />
    <main>
      <HeroSection />
      <BenefitsSection />
      <FeaturesSection />
      <ServicesSection />
      <!-- <TestimonialsSection /> -->
      <TeamSection />
      <!-- <FAQSection /> -->
      <ContactSection />
    </main>
    <AppFooter />
  </div>
</template>
