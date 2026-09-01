import { ref } from 'vue'

// Yandex Geosuggest API (HTTP, JSON) -- a dedicated suggest-as-you-type
// product, separate from the classic Maps JS API's `ymaps.SuggestView`
// widget. No script injection, no visual map, just a fetch call -- lighter
// weight, and this project's key has a 1000/day quota vs. 100/day for the
// JS API key. Contract confirmed against https://yandex.ru/dev/geosuggest/
// (2026-07-14): GET https://suggest-maps.yandex.ru/v1/suggest?apikey&text
// -> { results: [{ title: { text }, address: { formatted_address }, uri }] }.
// `uri` is meant to be handed to the Geocoder for a precise point -- more
// reliable than re-searching by the formatted text.
const SUGGEST_URL = 'https://suggest-maps.yandex.ru/v1/suggest'
const GEOCODE_URL = 'https://geocode-maps.yandex.ru/1.x/'

export interface AddressSuggestion {
  text: string
  uri?: string
}

export interface ResolvedAddress {
  value: string
  latitude?: number
  longitude?: number
}

export function useYandexSuggest() {
  const suggestApiKey = import.meta.env.VITE_YANDEX_GEOSUGGEST_API_KEY as string | undefined
  const geocoderApiKey = import.meta.env.VITE_YANDEX_GEOCODER_API_KEY as string | undefined
  const available = Boolean(suggestApiKey)

  const suggestions = ref<AddressSuggestion[]>([])
  const loading = ref(false)

  let debounceTimer: ReturnType<typeof setTimeout> | null = null
  let requestSeq = 0

  async function fetchSuggestions(text: string) {
    if (!suggestApiKey || text.trim().length < 3) {
      suggestions.value = []
      return
    }
    const seq = ++requestSeq
    loading.value = true
    try {
      const url = `${SUGGEST_URL}?apikey=${suggestApiKey}&text=${encodeURIComponent(text)}&print_address=1&results=5`
      const res = await fetch(url)
      if (!res.ok) throw new Error(`suggest request failed: ${res.status}`)
      const data = await res.json()
      if (seq !== requestSeq) return // a newer keystroke already superseded this response
      const results: Array<{ title?: { text?: string }; address?: { formatted_address?: string }; uri?: string }> =
        data?.results ?? []
      suggestions.value = results
        .map((r) => ({ text: r.address?.formatted_address || r.title?.text || '', uri: r.uri }))
        .filter((s) => s.text)
    } catch {
      if (seq === requestSeq) suggestions.value = []
    } finally {
      if (seq === requestSeq) loading.value = false
    }
  }

  // Debounced so every keystroke doesn't burn a request against the daily quota.
  function onInput(text: string) {
    if (debounceTimer) clearTimeout(debounceTimer)
    debounceTimer = setTimeout(() => fetchSuggestions(text), 300)
  }

  function clear() {
    suggestions.value = []
  }

  async function resolve(suggestion: AddressSuggestion): Promise<ResolvedAddress> {
    if (!geocoderApiKey) return { value: suggestion.text }
    try {
      // A suggestion's `uri` is a more precise geocoder lookup key than its
      // display text when present; fall back to text-based geocoding otherwise.
      const geocode = suggestion.uri || suggestion.text
      const url = `${GEOCODE_URL}?apikey=${geocoderApiKey}&format=json&geocode=${encodeURIComponent(geocode)}`
      const res = await fetch(url)
      if (!res.ok) return { value: suggestion.text }
      const data = await res.json()
      const pos: string | undefined =
        data?.response?.GeoObjectCollection?.featureMember?.[0]?.GeoObject?.Point?.pos
      if (!pos) return { value: suggestion.text }
      const [longitude, latitude] = pos.split(' ').map(Number)
      return { value: suggestion.text, latitude, longitude }
    } catch {
      return { value: suggestion.text } // best-effort: text alone is still a valid address
    }
  }

  return { available, suggestions, loading, onInput, clear, resolve }
}
