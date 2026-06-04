/**
 * Monitor Validation Schema
 * Client-side Zod validation rules for monitor configurations
 */
import { z } from 'zod';

export const monitorSchema = z.object({
  name: z
    .string({ required_error: 'Friendly name is required' })
    .min(1, 'Friendly name is required')
    .max(100, 'Friendly name must be less than 100 characters'),
  url: z
    .string({ required_error: 'Target URL is required' })
    .min(1, 'Target URL is required')
    .url('Please enter a valid target URL'),
  type: z
    .string({ required_error: 'Check type is required' })
    .min(1, 'Check type is required'),
  interval: z
    .number({ required_error: 'Interval is required' })
    .min(10, 'Interval must be at least 10 seconds'),
  timeout: z
    .number({ required_error: 'Timeout is required' })
    .min(1, 'Timeout must be at least 1 second'),
  status: z
    .string()
    .default('active'),
});

export type MonitorFields = z.infer<typeof monitorSchema>;
