export interface HistoryItem {
  id: number
  project_id: number
  request_id: number | null
  method: string
  url: string
  request_headers: string
  request_body: string
  response_status: number
  response_headers: string
  response_body: string
  duration_ms: number
  created_at: string
}
