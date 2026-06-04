/**
 * User Validation Schema
 * Client-side Zod validation rules for user profile/account editing & registration
 */
import { z } from 'zod';

export const userSchema = z.object({
  name: z
    .string({ required_error: 'Full name is required' })
    .min(1, 'Full name is required')
    .max(100, 'Name must be less than 100 characters'),
  email: z
    .string({ required_error: 'Email is required' })
    .min(1, 'Email is required')
    .email('Please enter a valid email address'),
  password: z
    .string()
    .min(6, 'Password must be at least 6 characters')
    .optional()
    .or(z.literal('')),
  role: z
    .string({ required_error: 'Access privilege is required' })
    .min(1, 'Access privilege is required'),
  is_blocked: z
    .boolean()
    .default(false),
});

export type UserFields = z.infer<typeof userSchema>;
