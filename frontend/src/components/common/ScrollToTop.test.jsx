import { fireEvent, render, screen } from '@testing-library/react'
import { Link, MemoryRouter, Route, Routes } from 'react-router-dom'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import ScrollToTop from './ScrollToTop'

const Page = ({ label, to }) => (
  <div>
    <p>{label}</p>
    <Link to={to}>Next</Link>
  </div>
)

const renderRoutes = () => render(
  <MemoryRouter initialEntries={['/a']}>
    <ScrollToTop />
    <Routes>
      <Route path="/a" element={<Page label="Page A" to="/b" />} />
      <Route path="/b" element={<Page label="Page B" to="/a" />} />
    </Routes>
  </MemoryRouter>,
)

describe('ScrollToTop', () => {
  beforeEach(() => {
    window.scrollTo = vi.fn()
    document.documentElement.scrollTop = 240
    document.body.scrollTop = 240
  })

  it('resets scroll when pathname changes', () => {
    renderRoutes()
    window.scrollTo.mockClear()

    fireEvent.click(screen.getByRole('link', { name: /next/i }))

    expect(screen.getByText('Page B')).toBeInTheDocument()
    expect(window.scrollTo).toHaveBeenCalledWith({ top: 0, left: 0, behavior: 'auto' })
    expect(document.documentElement.scrollTop).toBe(0)
    expect(document.body.scrollTop).toBe(0)
  })
})
