/**
 * Setting Validation Schema
 * Client-side Zod validation rules for global system settings including SMTP
 */
import { z } from 'zod';

export const settingSchema = z.object({
  system_name: z
    .string({ required_error: 'System name is required' })
    .min(1, 'System name is required')
    .max(100, 'System name must be less than 100 characters'),
  admin_email: z
    .string()
    .email('Please enter a valid email address')
    .optional()
    .or(z.literal('')),
  allow_registration: z
    .string()
    .default('true'),
  smtp_host: z.string().optional().or(z.literal('')),
  smtp_port: z.string().optional().or(z.literal('')),
  smtp_username: z.string().optional().or(z.literal('')),
  smtp_password: z.string().optional().or(z.literal('')),
  smtp_sender: z.string().optional().or(z.literal('')),
  smtp_encryption: z.string().optional().or(z.literal('')),
});

export type SettingFields = z.infer<typeof settingSchema>;
