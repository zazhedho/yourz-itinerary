import { useEffect, useState } from 'react'

import { getErrorMessage, getResponseData } from '../services/api'
import authService from '../services/authService'

const useRegisterStatus = () => {
  const [enabled, setEnabled] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    authService
      .getRegisterStatus()
      .then((response) => setEnabled(Boolean(getResponseData(response)?.enabled)))
      .catch((err) => {
        setEnabled(false)
        setError(getErrorMessage(err, 'Gagal mengecek status registrasi'))
      })
      .finally(() => setLoading(false))
  }, [])

  return { enabled, loading, error }
}

export default useRegisterStatus
