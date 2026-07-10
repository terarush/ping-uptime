import { describe, expect, it } from 'vitest'

import { cn } from '@/lib/utils'

describe('cn', () => {
  it('joins conditional class names', () => {
    expect(cn('inline-flex', false && 'hidden', ['items-center'], { 'font-medium': true })).toBe(
      'inline-flex items-center font-medium',
    )
  })

  it('merges conflicting Tailwind classes with the latest value winning', () => {
    expect(cn('px-2 py-1 text-sm', 'px-4 text-lg')).toBe('py-1 px-4 text-lg')
  })
})
