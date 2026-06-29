export const getTripAccess = (trip, user) => {
  const fallbackMember = (trip?.members || []).find((member) => member.user_id === user?.id)
  const fallbackRole = fallbackMember?.role || (trip?.owner_id === user?.id ? 'owner' : 'viewer')
  const currentRole = trip?.current_member_role || fallbackRole

  return {
    currentMemberId: trip?.current_member_id || fallbackMember?.id || '',
    currentRole,
    canEdit: trip?.can_edit ?? (currentRole === 'owner' || currentRole === 'editor'),
    canManageMembers: trip?.can_manage_members ?? currentRole === 'owner',
    canDelete: trip?.can_delete ?? currentRole === 'owner',
    canLeave: trip?.can_leave ?? currentRole !== 'owner',
  }
}

export const canAccessTripAction = (access, action) => {
  if (action === 'edit') return Boolean(access?.canEdit)
  if (action === 'manageMembers') return Boolean(access?.canManageMembers)
  if (action === 'delete') return Boolean(access?.canDelete)
  if (action === 'leave') return Boolean(access?.canLeave)
  return true
}
