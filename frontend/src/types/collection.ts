export interface TreeItem {
  id: number
  name: string
  type: 'folder' | 'request'
  method?: string
  url?: string
  children: TreeItem[]
  sort_order: number
}
