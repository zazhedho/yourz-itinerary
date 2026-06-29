import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'

import PageSkeleton from './PageSkeleton'

describe('PageSkeleton', () => {
  it('renders accessible loading skeleton rows', () => {
    render(<PageSkeleton label="Memuat trip" rows={3} />)

    expect(screen.getByRole('status', { name: /memuat trip/i })).toBeInTheDocument()
    expect(screen.getAllByTestId('skeleton-row')).toHaveLength(3)
  })
})
