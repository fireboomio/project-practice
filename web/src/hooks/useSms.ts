import { useCallback, useEffect, useRef, useState } from 'react'
import { StorageKey, createStorage } from '../utils'
import type { PromiseOr } from '../types'

const SEND_INTERVAL_TIME = 60

const storage = createStorage('localStorage')

export function useSmsCode() {
  const [countDown, setCountDown] = useState<number | null>(SEND_INTERVAL_TIME)
  const [loading, setLoading] = useState(false)
  const timer = useRef<number | null>(null)

  function stop() {
    if (timer.current) {
      clearInterval(timer.current)
      timer.current = null
    }
    setCountDown(SEND_INTERVAL_TIME)
    storage.removeItem(StorageKey.SmsSend)
  }

  const _startCountDown = useCallback(() => {
    timer.current = setInterval(() => {
      setCountDown((prev) => {
        if (prev && prev > 0) {
          return prev - 1
        } else {
          stop()
          return SEND_INTERVAL_TIME
        }
      })
    }, 1000)
  }, [])

  const startCountDown = useCallback(async (fn?: () => PromiseOr<boolean | undefined>) => {
    stop()
    if (fn) {
      setLoading(true)
      const ret = await fn()
      // 返回 false 可以中断
      if (ret === false) {
        setLoading(false)
        return
      }
      setLoading(false)
    }
    storage.setItem(StorageKey.SmsSend, new Date().toUTCString())
    _startCountDown()
  }, [_startCountDown])

  useEffect(() => {
    const stored = storage.getItem(StorageKey.SmsSend)
    if (stored) {
      const diffSecond = (+new Date() - (+new Date(stored))) / 1000
      if (diffSecond > SEND_INTERVAL_TIME) {
        storage.removeItem(StorageKey.SmsSend)
      } else {
        // 有剩余时间
        setCountDown(SEND_INTERVAL_TIME - diffSecond)
        startCountDown()
      }
    }
  }, [setCountDown, startCountDown])

  useEffect(() => {
    return () => {
      if (timer.current) {
        clearInterval(timer.current)
        timer.current = null
      }
    }
  }, [])

  return {
    loading,
    countDown,
    isCountingDown: countDown !== SEND_INTERVAL_TIME,
    startCountDown,
  }
}
