<script setup lang="ts">
import Navbar from '@/components/Navbar.vue'
import HeroSection from '@/components/HeroSection.vue'
import Sponsors from '@/components/Sponsors.vue'
import Benefits from '@/components/Benefits.vue'
import Features from '@/components/Features.vue'
import Services from '@/components/Services.vue'
import Testimonials from '@/components/Testimonials.vue'
import TeamSection from '@/components/TeamSection.vue'
import Community from '@/components/Community.vue'
import PricingSection from '@/components/PricingSection.vue'
import FAQSection from '@/components/FAQSection.vue'
import ContactSection from '@/components/ContactSection.vue'
import Footer from '@/components/Footer.vue'
</script>

<template>
  <div class="min-h-screen bg-background">
    <Navbar />
    <main>
      <HeroSection />
      <Sponsors />
      <Benefits />
      <Features />
      <Services />
      <Testimonials />
      <TeamSection />
      <Community />
      <PricingSection />
      <FAQSection />
      <ContactSection />
    </main>
    <Footer />
  </div>
</template>
