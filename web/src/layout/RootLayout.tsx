import { Outlet, useLocation, useNavigate } from '@modern-js/runtime/router'
import { ConfigProvider } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import { useEffect, useState } from 'react'

import FullScreenLoading from '@/components/FullScreenLoading'
import { useAuth } from '@/hooks'

export default function Layout() {
  const location = useLocation()
  const navigate = useNavigate()
  const {
    init: initUser,
    ready: initUserReady,
    computed: { logged }
  } = useAuth()
  const [ready, setReady] = useState(false)

  useEffect(() => {
    if (initUserReady && !logged && location.pathname !== Path.Login) {
      navigate(Path.Login, { replace: true })
    }
  }, [initUserReady, location, logged, navigate])

  useEffect(() => {
    initUser().then(async user => {
      if (user) {
        // do something else
      } else {
        navigate(Path.Login, { replace: true })
      }
      setReady(true)
    })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return ready ? (
    <ConfigProvider
      locale={zhCN}
      theme={{
        cssVar: true,
        hashed: false
      }}
    >
      <Outlet />
    </ConfigProvider>
  ) : (
    <FullScreenLoading />
  )
}
