<script setup lang="ts">
import Navbar from '@/components/Navbar.vue'
import HeroSection from '@/components/HeroSection.vue'
import Benefits from '@/components/Benefits.vue'
import Features from '@/components/Features.vue'
import Services from '@/components/Services.vue'
import TeamSection from '@/components/TeamSection.vue'
import ContactSection from '@/components/ContactSection.vue'
import Footer from '@/components/Footer.vue'
</script>

<template>
  <div class="min-h-screen bg-background">
    <Navbar />
    <main>
      <HeroSection />
      <Benefits />
      <Features />
      <Services />
      <!-- <Testimonials /> -->
      <TeamSection />
      <!-- <FAQSection /> -->
      <ContactSection />
    </main>
    <Footer />
  </div>
</template>
