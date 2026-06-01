export interface Request {
  id: string
  collection_id: string
  name: string
  method: string
  url: string
  headers: string
  params: string
  body: string
  auth: string
  script: string
  sort_order: number
  created_at: string
  updated_at: string
}
