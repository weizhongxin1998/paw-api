export interface KvItem {
  id: string
  key: string
  value: string
  description: string
  enabled: boolean
}

export interface RequestBody {
  type: 'none' | 'form-data' | 'x-www-form-urlencoded' | 'raw' | 'binary'
  formData?: (KvItem & { fieldType: 'text' | 'file' })[]
  urlEncoded?: KvItem[]
  raw?: { subType: string; content: string }
  binary?: { fileName: string; filePath: string }
}

export interface AuthConfig {
  type: 'none' | 'bearer' | 'basic' | 'apikey'
  token?: string
  username?: string
  password?: string
  apiKey?: string
  apiValue?: string
  addTo?: 'header' | 'query'
}
