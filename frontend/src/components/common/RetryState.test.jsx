import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'

import RetryState from './RetryState'

describe('RetryState', () => {
  it('shows message and calls retry action', async () => {
    const user = userEvent.setup()
    const onRetry = vi.fn()

    render(<RetryState message="Internet lambat" onRetry={onRetry} />)

    expect(screen.getByText('Internet lambat')).toBeInTheDocument()
    await user.click(screen.getByRole('button', { name: /coba lagi/i }))
    expect(onRetry).toHaveBeenCalledTimes(1)
  })
})
