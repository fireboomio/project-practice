import type { RendererEnv, RenderOptions } from 'amis'
import { message, notification } from 'antd'
import _copy from 'copy-to-clipboard'

import { StorageKey, tokenStorage } from './storage'

export const theme = 'antd'

// @ts-ignore
export const fetcher: RenderOptions['fetcher'] = ({ url, method, data, config, headers }) => {
  const isGet = method !== 'post' && method !== 'put' && method !== 'patch'
  const _headers = config?.headers || headers || {}
  if (!headers['Content-Type']) {
    headers['Content-Type'] = 'application/json'
  }
  const token = tokenStorage.getItem(StorageKey.AccessToken)
  if (token) {
    _headers.Authorization = `Bearer ${token}`
  }
  return fetch((url as string) + (isGet ? `?${new URLSearchParams(data).toString()}` : ''), {
    method,
    headers: _headers,
    credentials: 'include',
    body: isGet ? undefined : JSON.stringify(data)
  })
}

export const notify: RendererEnv['notify'] = (type, msg) => {
  if (message[type]) {
    message[type](msg)
  } else {
    console.warn('[Notify]', type, msg)
  }
  console.log('[notify]', type, msg)
}

export const alert: RendererEnv['alert'] = (msg, title) => {
  notification.error({
    message: title,
    description: msg
  })
}

export const copy: RendererEnv['copy'] = (contents, options = {}) => {
  const ret = _copy(contents, options)
  ret && (!options || options.shutup !== true) && message.info('内容已拷贝到剪切板')
  return ret
}
