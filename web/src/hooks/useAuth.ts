import { create } from 'zustand'
import { message } from 'antd'

import { StorageKey, TokenStorageKey, tokenStorage, setTokenStorage } from '../utils'

export type LoginType = 'password' | 'sms'

async function fetchMe() {
  const resp = await client.query({
    operationName: 'user/me'
  })
  if (!resp.error) {
    const user = resp.data!.data!
    // 处理用户被删除但是仍然能返回的情况
    if (user.id) {
      return { user, error: null }
    }
  }
  return { user: null, error: resp.error }
}

export type AuthState<User> = {
  user: User | null
  ready: boolean
  loading: boolean
  computed: {
    logged: boolean
    username?: string
    avatar?: string
  }
  init: () => Promise<User | null>
  login: (
    type: LoginType,
    args: {
      username?: string
      password?: string
      phone?: string
      code?: string
      rememberMe?: boolean
    }
  ) => Promise<[boolean, User | Error]>
  logout: (showMsg?: boolean) => void
  refreshMe: () => Promise<User | null>
  updateMyInfo: (args: Omit<User, 'userId'>) => Promise<boolean>
  updateMyPassword: (newPassword: string) => Promise<boolean>
}

export const useAuth = create<AuthState>((set, get) => ({
  user: null,
  ready: false,
  loading: false,
  computed: {
    get logged() {
      return Boolean(get().user)
    },
    get username() {
      return get().user?.name
    },
    get avatar() {
      return get().user?.avatar
    }
  },
  async init() {
    set({ loading: true })
    // 先读取token存在哪里
    const tokenStoreType = localStorage.getItem(StorageKey.TokenStorageType)
    if (tokenStoreType === TokenStorageKey.SessionStorage) {
      setTokenStorage('sessionStorage')
    }
    const accessToken = tokenStorage.getItem(StorageKey.AccessToken)
    if (accessToken) {
      client.setAuthorizationToken(accessToken)
      const { user } = await fetchMe()
      if (user) {
        set({ accessToken, user, ready: true, loading: false })
        return user
      }
      tokenStorage.removeItem(StorageKey.AccessToken)
      client.unsetAuthorization()
    }
    set({ ready: true, loading: false })
    return null
  },
  async login(loginType, { rememberMe = false, ...args }) {
    const { error, data } = await client.mutate({
      operationName: 'user/casdoor/login',
      input: {
        loginType,
        ...args
      }
    })
    if (!error) {
      const { accessToken, refreshToken } = data!.data!.data!
      client.setAuthorizationToken(accessToken!)
      const { user, error } = await fetchMe()
      if (user) {
        set({ accessToken, user })
        // 记住我用localStorage存，不记住用sessionStorage
        if (rememberMe) {
          localStorage.setItem(StorageKey.TokenStorageType, TokenStorageKey.LocalStorage)
          localStorage.setItem(StorageKey.AccessToken, accessToken!)
          // 记住refreshToken
          localStorage.setItem(StorageKey.RefreshToken, refreshToken!)
        } else {
          localStorage.setItem(StorageKey.TokenStorageType, TokenStorageKey.SessionStorage)
          sessionStorage.setItem(StorageKey.AccessToken, accessToken!)
        }
        return [true, user]
      }
      client.unsetAuthorization()
      return [false, error]
    } else {
      return [false, error]
    }
  },
  async refreshMe() {
    const { user } = await fetchMe()
    if (user) {
      set({ user })
      return user
    }
    // 退出
    get().logout(false)
    return null
  },
  logout(showMsg = true) {
    tokenStorage.removeItem(StorageKey.AccessToken)
    tokenStorage.removeItem(StorageKey.RefreshToken)
    set({ accessToken: null, user: null })
    if (showMsg) {
      message.success('你已成功退出')
    }
    window.location.href = Path.Login
  },
  async updateMyInfo(args) {
    if (!get().computed.logged) {
      return false
    }
    const { user } = get()
    const { error } = await client.mutate({
      operationName: 'user/casdoor/updateUser',
      input: {
        userId: user!.userId!,
        ...args
      }
    })
    if (!error) {
      set({ user: { ...user!, ...args } })
      message.success('个人信息更新成功')
      return true
    }
    return false
  },
  async updateMyPassword(newPwd) {
    const { user, logout } = get()
    const { error } = await client.mutate({
      operationName: 'user/casdoor/updateUser',
      input: {
        userId: user!.userId!,
        password: newPwd
      }
    })
    if (!error) {
      message.success('密码已修改，请重新登录')
      setTimeout(() => {
        logout(false)
      }, 2000)
      return true
    } else {
      message.error('密码修改失败')
      return false
    }
  }
}))

// token 过期
export function expiredLogin() {
  window.location.href = Path.Login + '?expired=1'
  tokenStorage.removeItem(StorageKey.AccessToken)
  tokenStorage.removeItem(StorageKey.RefreshToken)
}
