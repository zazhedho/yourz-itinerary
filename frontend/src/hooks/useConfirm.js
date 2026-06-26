import { useCallback, useState } from 'react'

export const useConfirm = () => {
  const [state, setState] = useState({ open: false })

  const confirm = useCallback((options) => {
    return new Promise((resolve) => {
      setState({
        open: true,
        ...options,
        onCancel: () => {
          setState({ open: false })
          resolve(false)
        },
        onConfirm: () => {
          setState({ open: false })
          resolve(true)
        },
      })
    })
  }, [])

  return { confirm, dialogProps: state }
}
