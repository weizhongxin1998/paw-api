import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Request } from '../types/request'

export interface ResponseData {
  status: number
  status_text: string
  headers: Record<string, string[]>
  body: string
  duration_ms: number
}

export const useRequestStore = defineStore('request', () => {
  const requests = ref<Request[]>([])
  const currentRequest = ref<Request | null>(null)
  const lastResponse = ref<ResponseData | null>(null)

  function setCurrentRequest(r: Request) {
    currentRequest.value = r
  }

  function setLastResponse(resp: ResponseData) {
    lastResponse.value = resp
  }

  return { requests, currentRequest, lastResponse, setCurrentRequest, setLastResponse }
})
