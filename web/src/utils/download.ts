import { showNotification } from './notification'

/**
 * 将 Markdown 文本内容转化成 docx 并下载
 * @param content Markdown 内容
 * @param filename 希望保存的文件名
 */
export async function downloadDocxFromMarkdown(content: string, filename?: string) {
  let blob: Blob
  try {
    const resp = await fetch('https://pandoc.org/cgi-bin/pandoc-server.cgi', {
      method: 'post',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        text: content,
        to: 'docx',
        from: 'markdown',
        standalone: true,
        'embed-resources': false,
        'table-of-contents': false,
        'number-sections': false,
        citeproc: false,
        'html-math-method': 'plain',
        wrap: 'auto',
        'highlight-style': 'pygments',
        template: null
      })
    })
    blob = await resp.blob()
  } catch (error) {
    showNotification('warning', { message: '转换word文档失败，下载源文档' })
    blob = new Blob([content], { type: 'text/plain' })
  }
  const downloadLink = document.createElement('a')
  downloadLink.href = URL.createObjectURL(blob)
  downloadLink.download = filename ? `${filename}.docx` : 'generated.docx'
  downloadLink.click()
  downloadLink.remove()
}

/**
 * 下载图片文件
 * @param img 图片 url
 * @param filename 希望保存的文件名
 */
export function downloadImage(img: string, filename?: string) {
  // fetch(img).then(res => res.blob()).then(blob => {
  //   const downloadLink = document.createElement('a')
  //   downloadLink.href = URL.createObjectURL(blob)
  //   downloadLink.download = filename ? `${filename}.jpeg` : 'download.jpeg'
  //   downloadLink.click()
  //   downloadLink.remove()
  // })
  const downloadLink = document.createElement('a')
  downloadLink.href = img
  downloadLink.download = filename ? `${filename}.jpeg` : 'download.jpeg'
  downloadLink.click()
  downloadLink.remove()
}
