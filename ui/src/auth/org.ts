const KEY = "autoglue.org"

let cache: string | null = localStorage.getItem(KEY)

export const orgStore = {
  get(): string | null {
    return cache
  },
  set(id: string) {
    cache = id
    localStorage.setItem(KEY, id)
    window.dispatchEvent(new CustomEvent("autoglue:org-change", { detail: id }))
  },
  subscribe(fn: (id: string | null) => void) {
    const onCustom = (e: Event) => fn((e as CustomEvent<string>).detail ?? null)
    const onStorage = (e: StorageEvent) => {
      if (e.key === KEY) {
        cache = e.newValue
        fn(cache)
      }
    }
    window.addEventListener("autoglue:org-change", onCustom as EventListener)
    window.addEventListener("storage", onStorage)
    return () => {
      window.removeEventListener("autoglue:org-change", onCustom as EventListener)
      window.removeEventListener("storage", onStorage)
    }
  },
}
