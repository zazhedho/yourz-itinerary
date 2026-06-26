import { act, fireEvent, render, screen } from '@testing-library/react'
import { useState } from 'react'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import GooglePlaceAutocomplete from './GooglePlaceAutocomplete'

const fetchAutocompleteSuggestions = vi.fn()

const ControlledAutocomplete = () => {
  const [value, setValue] = useState('')

  return (
    <GooglePlaceAutocomplete
      isLoaded
      name="location_name"
      value={value}
      onChange={(event) => setValue(event.target.value)}
      onPlaceSelect={() => {}}
    />
  )
}

describe('GooglePlaceAutocomplete', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    fetchAutocompleteSuggestions.mockResolvedValue({ suggestions: [] })
    window.google = {
      maps: {
        importLibrary: vi.fn().mockResolvedValue({
          AutocompleteSessionToken: vi.fn(),
          AutocompleteSuggestion: { fetchAutocompleteSuggestions },
        }),
      },
    }
  })

  afterEach(() => {
    vi.useRealTimers()
    delete window.google
    vi.clearAllMocks()
  })

  it('does not open suggestions for an existing prefilled location', async () => {
    render(
      <GooglePlaceAutocomplete
        isLoaded
        name="location_name"
        value="Monas, Jakarta Pusat"
        onChange={() => {}}
        onPlaceSelect={() => {}}
      />,
    )

    await act(async () => {
      await vi.runAllTimersAsync()
    })

    expect(fetchAutocompleteSuggestions).not.toHaveBeenCalled()
  })

  it('searches only after the user types', async () => {
    render(<ControlledAutocomplete />)

    await act(async () => {
      await vi.runAllTimersAsync()
    })
    fireEvent.change(screen.getByRole('textbox'), { target: { value: 'Monas' } })
    await act(async () => {
      await vi.advanceTimersByTimeAsync(300)
    })

    expect(fetchAutocompleteSuggestions).toHaveBeenCalled()
  })
})
