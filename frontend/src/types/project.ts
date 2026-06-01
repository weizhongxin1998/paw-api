export interface Project {
  id: string
  name: string
  description: string
  created_at: string
  updated_at: string
}

export interface Collection {
  id: string
  project_id: string
  parent_id?: string
  name: string
  sort_order: number
  created_at: string
  updated_at: string
}
