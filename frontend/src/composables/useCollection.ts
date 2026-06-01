export function useCollection() {
  function flattenTree(items: Array<{ id: string; parent_id: string | null }>): string[] {
    const ids: string[] = []
    function walk(items: Array<{ id: string; parent_id: string | null }>, parentId: string | null) {
      for (const item of items) {
        if (item.parent_id === parentId) {
          ids.push(item.id)
          walk(items, item.id)
        }
      }
    }
    walk(items, null)
    return ids
  }

  return { flattenTree }
}
