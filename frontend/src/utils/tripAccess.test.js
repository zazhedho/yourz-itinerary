import { describe, expect, it } from 'vitest'

import { canAccessTripAction, getTripAccess } from './tripAccess'

describe('tripAccess', () => {
  it('uses backend access flags when present', () => {
    const access = getTripAccess({
      current_member_id: 'member-1',
      current_member_role: 'viewer',
      can_edit: true,
      can_manage_members: false,
      can_delete: false,
      can_leave: true,
    }, { id: 'user-1' })

    expect(access).toMatchObject({
      currentMemberId: 'member-1',
      currentRole: 'viewer',
      canEdit: true,
      canManageMembers: false,
      canDelete: false,
      canLeave: true,
    })
  })

  it('falls back to trip members when backend flags are missing', () => {
    const access = getTripAccess({
      owner_id: 'owner-1',
      members: [{ id: 'member-2', user_id: 'user-2', role: 'editor' }],
    }, { id: 'user-2' })

    expect(access.canEdit).toBe(true)
    expect(access.canManageMembers).toBe(false)
    expect(canAccessTripAction(access, 'edit')).toBe(true)
    expect(canAccessTripAction(access, 'manageMembers')).toBe(false)
  })
})
