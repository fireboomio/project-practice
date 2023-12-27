export function isIOS() {
  const userAgent = navigator.userAgent.toLowerCase()
  return /iphone|ipad|ipod/.test(userAgent)
}

export function isMobileWidth() {
  return window.innerWidth < 768
}

// 根据ua 和设备尺寸判断是否是微信浏览器
export function isInWeixin() {
  return (/micromessenger/.test(navigator.userAgent.toLowerCase())) && isMobileWidth()
}