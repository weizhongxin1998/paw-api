export interface HttpResponse {
  status: number
  time: number
  size: number
  headers: Record<string, string>
  cookies: { name: string; value: string; domain: string; path: string }[]
  body: string
  rawRequest: string
  curlCommand: string
}
