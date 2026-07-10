import { describe, expect, it } from 'vitest'

import { loginSchema, setupSchema } from '@/validations/auth'
import { monitorSchema } from '@/validations/monitor'
import { notificationChannelSchema } from '@/validations/notification-channel'
import { settingSchema } from '@/validations/setting'
import { statusPageSchema } from '@/validations/status-page'
import { userSchema } from '@/validations/user'

describe('validation schemas', () => {
  it('validates login credentials', () => {
    expect(loginSchema.safeParse({ email: 'user@example.com', password: 'secret1' }).success).toBe(
      true,
    )

    const result = loginSchema.safeParse({ email: 'invalid', password: '123' })

    expect(result.success).toBe(false)
    if (!result.success) {
      expect(result.error.flatten().fieldErrors.email).toContain(
        'Please enter a valid email address',
      )
      expect(result.error.flatten().fieldErrors.password).toContain(
        'Password must be at least 6 characters',
      )
    }
  })

  it('validates setup data', () => {
    expect(
      setupSchema.safeParse({
        name: 'Jane Doe',
        email: 'jane@example.com',
        password: 'secret1',
      }).success,
    ).toBe(true)

    const result = setupSchema.safeParse({
      name: '',
      email: 'bad-email',
      password: '123',
    })

    expect(result.success).toBe(false)
    if (!result.success) {
      const errors = result.error.flatten().fieldErrors
      expect(errors.name).toContain('Full name is required')
      expect(errors.email).toContain('Please enter a valid email address')
      expect(errors.password).toContain('Password must be at least 6 characters')
    }
  })

  it('validates monitor configuration and applies the active status default', () => {
    const result = monitorSchema.safeParse({
      name: 'API',
      url: 'https://example.com/health',
      type: 'http',
      interval: 60,
      timeout: 5,
    })

    expect(result.success).toBe(true)
    if (result.success) {
      expect(result.data.status).toBe('active')
    }

    const invalid = monitorSchema.safeParse({
      name: '',
      url: 'not-a-url',
      type: '',
      interval: 5,
      timeout: 0,
    })

    expect(invalid.success).toBe(false)
    if (!invalid.success) {
      const errors = invalid.error.flatten().fieldErrors
      expect(errors.name).toContain('Friendly name is required')
      expect(errors.url).toContain('Please enter a valid target URL')
      expect(errors.type).toContain('Check type is required')
      expect(errors.interval).toContain('Interval must be at least 10 seconds')
      expect(errors.timeout).toContain('Timeout must be at least 1 second')
    }
  })

  it('validates notification channel configuration and defaults enabled', () => {
    const result = notificationChannelSchema.safeParse({
      name: 'Ops',
      type: 'webhook',
      config: '{}',
    })

    expect(result.success).toBe(true)
    if (result.success) {
      expect(result.data.enabled).toBe(true)
    }

    const invalid = notificationChannelSchema.safeParse({
      name: '',
      type: '',
      config: '',
    })

    expect(invalid.success).toBe(false)
    if (!invalid.success) {
      const errors = invalid.error.flatten().fieldErrors
      expect(errors.name).toContain('Channel name is required')
      expect(errors.type).toContain('Channel type is required')
      expect(errors.config).toContain('Config is required')
    }
  })

  it('validates status page fields and defaults monitor ids', () => {
    const result = statusPageSchema.safeParse({
      name: 'Public Status',
      slug: 'public-status',
      description: '',
    })

    expect(result.success).toBe(true)
    if (result.success) {
      expect(result.data.monitor_ids).toEqual([])
    }

    const invalid = statusPageSchema.safeParse({
      name: '',
      slug: 'Invalid Slug!',
      description: 'x'.repeat(501),
    })

    expect(invalid.success).toBe(false)
    if (!invalid.success) {
      const errors = invalid.error.flatten().fieldErrors
      expect(errors.name).toContain('Status page name is required')
      expect(errors.slug).toContain(
        'Slug can only contain lowercase letters, numbers, dashes, and underscores',
      )
      expect(errors.description).toContain('Description must be less than 500 characters')
    }
  })

  it('validates settings while allowing empty optional integration values', () => {
    const result = settingSchema.safeParse({
      system_name: 'Ping Uptime',
      admin_email: '',
      smtp_host: '',
      smtp_port: '',
      smtp_username: '',
      smtp_password: '',
      smtp_sender: '',
      smtp_encryption: '',
      discord_bot_token: '',
      telegram_bot_token: '',
      discord_client_id: '',
    })

    expect(result.success).toBe(true)
    if (result.success) {
      expect(result.data.allow_registration).toBe('true')
    }

    const invalid = settingSchema.safeParse({
      system_name: '',
      admin_email: 'not-email',
    })

    expect(invalid.success).toBe(false)
    if (!invalid.success) {
      const errors = invalid.error.flatten().fieldErrors
      expect(errors.system_name).toContain('System name is required')
      expect(errors.admin_email).toContain('Please enter a valid email address')
    }
  })

  it('validates users and allows blank optional passwords', () => {
    const result = userSchema.safeParse({
      name: 'Admin User',
      email: 'admin@example.com',
      password: '',
      role: 'admin',
    })

    expect(result.success).toBe(true)
    if (result.success) {
      expect(result.data.is_blocked).toBe(false)
    }

    const invalid = userSchema.safeParse({
      name: '',
      email: 'bad-email',
      password: '123',
      role: '',
    })

    expect(invalid.success).toBe(false)
    if (!invalid.success) {
      const errors = invalid.error.flatten().fieldErrors
      expect(errors.name).toContain('Full name is required')
      expect(errors.email).toContain('Please enter a valid email address')
      expect(errors.password).toContain('Password must be at least 6 characters')
      expect(errors.role).toContain('Access privilege is required')
    }
  })
})
