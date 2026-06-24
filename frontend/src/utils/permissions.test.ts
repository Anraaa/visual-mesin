import { describe, it, expect, vi, beforeEach } from 'vitest'
import { can, canAny, canModule } from './permissions'
import { useAuthStore } from '../stores/authStore'

vi.mock('../stores/authStore', () => {
  const store = {
    user: null,
    getState: () => ({ user: null }),
  }
  return {
    useAuthStore: Object.assign(vi.fn(() => store), store),
  }
})

describe('can', () => {
  it('returns false when no user', () => {
    expect(can('view-dashboard')).toBe(false)
  })

  it('returns false when user has no permissions', () => {
    expect(can('view-dashboard', { permissions: [] } as any)).toBe(false)
  })

  it('returns true when user has exact permission', () => {
    expect(can('view-dashboard', { permissions: ['view-dashboard'] } as any)).toBe(true)
  })

  it('returns true when user has view-any permission', () => {
    expect(can('view-dashboard', { permissions: ['view-any-dashboard'] } as any)).toBe(true)
  })

  it('returns false when user lacks permission', () => {
    expect(can('view-dashboard', { permissions: ['view-other'] } as any)).toBe(false)
  })

  it('returns true for view-any-* when checking specific view-*', () => {
    expect(can('view-user', { permissions: ['view-any-user'] } as any)).toBe(true)
  })

  it('handles multi-word permission names', () => {
    expect(can('view-resource-connection', { permissions: ['view-any-resource-connection'] } as any)).toBe(true)
    expect(can('view-resource-connection', { permissions: ['view-resource-connection'] } as any)).toBe(true)
  })
})

describe('canAny', () => {
  it('returns true if any permission matches', () => {
    expect(canAny(['view-dashboard', 'view-user'], { permissions: ['view-user'] } as any)).toBe(true)
  })

  it('returns false if no permission matches', () => {
    expect(canAny(['view-dashboard', 'view-user'], { permissions: ['view-other'] } as any)).toBe(false)
  })

  it('returns false for empty list', () => {
    expect(canAny([], { permissions: ['view-dashboard'] } as any)).toBe(false)
  })
})

describe('canModule', () => {
  it('returns true when user has view-module', () => {
    expect(canModule('dashboard', { permissions: ['view-dashboard'] } as any)).toBe(true)
  })

  it('returns true when user has view-any-module', () => {
    expect(canModule('dashboard', { permissions: ['view-any-dashboard'] } as any)).toBe(true)
  })

  it('returns false when user lacks permission', () => {
    expect(canModule('dashboard', { permissions: ['view-other'] } as any)).toBe(false)
  })

  it('works with modules containing hyphens', () => {
    expect(canModule('resource-connection', { permissions: ['view-resource-connection'] } as any)).toBe(true)
    expect(canModule('activity-log', { permissions: ['view-any-activity-log'] } as any)).toBe(true)
  })
})
