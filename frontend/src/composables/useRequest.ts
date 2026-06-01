export function useRequest() {
  function parseParams(url: string): Record<string, string> {
    const idx = url.indexOf('?')
    if (idx === -1) return {}
    const params: Record<string, string> = {}
    new URLSearchParams(url.slice(idx)).forEach((v, k) => { params[k] = v })
    return params
  }

  return { parseParams }
}
