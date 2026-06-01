import { useEnvironmentStore, type EnvVariable } from '../stores/environment'

export function useVariableResolver() {
  const envStore = useEnvironmentStore()

  function resolve(text: string, variables: EnvVariable[]): string {
    let result = text
    for (const v of variables) {
      if (!v.enabled) continue
      const pattern = new RegExp(`\\{\\{\\s*${escapeRegex(v.key)}\\s*\\}\\}`, 'g')
      result = result.replace(pattern, v.value)
    }
    return result
  }

  function escapeRegex(str: string): string {
    return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  }

  return { resolve }
}
