export interface Environment {
  id: number
  project_id: number
  name: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface EnvVariable {
  id: number
  key: string
  value: string
  enabled: boolean
}
