/**
 * Status Page Validation Schema
 * Client-side Zod validation rules for status pages
 */
import { z } from 'zod';

export const statusPageSchema = z.object({
  name: z
    .string({ required_error: 'Status page name is required' })
    .min(1, 'Status page name is required')
    .max(100, 'Name must be less than 100 characters'),
  slug: z
    .string({ required_error: 'Slug is required' })
    .min(1, 'Slug is required')
    .regex(/^[a-z0-9-_]+$/, 'Slug can only contain lowercase letters, numbers, dashes, and underscores')
    .max(50, 'Slug must be less than 50 characters'),
  description: z
    .string()
    .max(500, 'Description must be less than 500 characters')
    .optional()
    .or(z.literal('')),
  monitor_ids: z
    .array(z.number())
    .default([]),
});

export type StatusPageFields = z.infer<typeof statusPageSchema>;
