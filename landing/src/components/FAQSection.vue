<script setup lang="ts">
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion'

const faqs = [
  { question: 'Is this template free?', answer: 'Yes, this template is completely free to use for both personal and commercial projects.' },
  { question: 'Duis aute irure dolor in reprehenderit in voluptate velit?', answer: 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.' },
  { question: 'Lorem ipsum dolor sit amet Consectetur natus dolor minus quibusdam?', answer: 'Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.' },
  { question: 'Excepteur sint occaecat cupidata non proident sunt?', answer: 'Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.' },
  { question: 'Enim ad minim veniam, quis nostrud exercitation ullamco laboris?', answer: 'Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation.' },
]
</script>

<template>
  <section id="faq" class="py-20 md:py-28">
    <div class="mx-auto max-w-3xl px-4 sm:px-6 lg:px-8">
      <div class="text-center mb-12">
        <h2 class="text-3xl font-bold tracking-tight sm:text-4xl">FAQS</h2>
        <p class="mt-3 text-lg text-muted-foreground">Common Questions</p>
      </div>

      <Accordion type="single" collapsible class="w-full">
        <AccordionItem v-for="(faq, i) in faqs" :key="i" :value="`item-${i}`">
          <AccordionTrigger class="text-left font-medium text-sm">
            {{ faq.question }}
          </AccordionTrigger>
          <AccordionContent class="text-sm text-muted-foreground">
            {{ faq.answer }}
          </AccordionContent>
        </AccordionItem>
      </Accordion>
    </div>
  </section>
</template>
