/**
 * Zod Authentication & Registration Schemas
 * Used on the client side with vee-validate to enforce
 * input formats before submitting payloads to the backend.
 */
import { z } from 'zod';

// Login inputs validation schema
export const loginSchema = z.object({
  email: z
    .string({ required_error: 'Email is required' })
    .min(1, 'Email is required')
    .email('Please enter a valid email address'),
  password: z
    .string({ required_error: 'Password is required' })
    .min(6, 'Password must be at least 6 characters'),
});

// Setup / Signup registration inputs validation schema
export const setupSchema = z.object({
  name: z
    .string({ required_error: 'Full name is required' })
    .min(1, 'Full name is required')
    .max(100, 'Name must be less than 100 characters'),
  email: z
    .string({ required_error: 'Email is required' })
    .min(1, 'Email is required')
    .email('Please enter a valid email address'),
  password: z
    .string({ required_error: 'Password is required' })
    .min(6, 'Password must be at least 6 characters'),
});

export type LoginFields = z.infer<typeof loginSchema>;
export type SetupFields = z.infer<typeof setupSchema>;
