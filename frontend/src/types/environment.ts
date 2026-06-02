export interface Environment {
  id: number
  project_id: number
  name: string
  is_active: boolean
  variables: EnvVariable[]
}

export interface EnvVariable {
  id: number
  key: string
  value: string
  enabled: boolean
}
