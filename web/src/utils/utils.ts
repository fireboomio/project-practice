/**
 * 生成随机 id
 * @returns 随机 id
 */
export function randomId() {
  return Math.random().toString(36).substring(2);
}

/**
 * 获取 css 变量名对应的值
 * @param varName css 变量名
 * @returns 返回 css 变量名对应的值
 */
export function getCSSVar(varName: string) {
  return getComputedStyle(document.body).getPropertyValue(varName).trim();
}