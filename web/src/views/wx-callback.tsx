import { FullScreenLoading } from '@/components/FullScreenLoading'
import { WX_LOGIN_POST_MESSAGE } from '../utils'
import { useEffect, useState } from 'react'
import { useSearchParams } from 'react-router-dom'

export default function WXLoginCallbackPage() {
  const [params] = useSearchParams()
  const [err, setErr] = useState('')

  useEffect(() => {
    const code = params.get('code')
    const state = params.get('state')
    if (code) {
      window.parent.postMessage({ code, state, type: WX_LOGIN_POST_MESSAGE }, '*')
    } else {
      setErr('无效的code')
    }
  }, [params])

  return err ? <p>{err}</p> : <FullScreenLoading />
}
