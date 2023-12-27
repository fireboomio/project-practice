export function deepClone<T = any>(value: T) {
  // 优先使用浏览器自带的
  if (typeof structuredClone === 'function') {
    return structuredClone
  } else {
    try {
      // 降级使用 json
      return JSON.parse(JSON.stringify(value))
    } catch (error) {
      // 降级
      try {
        return Object.assign({}, value)
      } catch (error) {
        // 处理不了的
        return value
      }
    }
  }
}