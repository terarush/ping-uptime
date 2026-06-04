/**
 * Notification Channel Validation Schema
 * Client-side Zod validation rules for notification integrations
 */
import { z } from 'zod';

export const notificationChannelSchema = z.object({
  name: z
    .string({ required_error: 'Channel name is required' })
    .min(1, 'Channel name is required')
    .max(100, 'Channel name must be less than 100 characters'),
  type: z
    .string({ required_error: 'Channel type is required' })
    .min(1, 'Channel type is required'),
  enabled: z
    .boolean()
    .default(true),
  config: z
    .string({ required_error: 'Config JSON is required' })
    .min(2, 'Config is required'), // Must at least be empty JSON object "{}"
});

export type NotificationChannelFields = z.infer<typeof notificationChannelSchema>;
