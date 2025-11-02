export type TokenPair = {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
}

const KEY = "autoglue.tokens"
const EVT = "autoglue.auth-change"

let cache: TokenPair | null = read()

function read(): TokenPair | null {
  try {
    const raw = localStorage.getItem(KEY)
    return raw ? (JSON.parse(raw) as TokenPair) : null
  } catch {
    return null
  }
}

function write(tokens: TokenPair | null) {
  if (tokens) localStorage.setItem(KEY, JSON.stringify(tokens))
  else localStorage.removeItem(KEY)
}

function emit(tokens: TokenPair | null) {
  // include payload for convenience
  window.dispatchEvent(new CustomEvent<TokenPair | null>(EVT, { detail: tokens }))
}

export const authStore = {
  /** Current tokens (from in-memory cache). */
  get(): TokenPair | null {
    return cache
  },

  /** Set tokens; updates memory, localStorage, broadcasts event. */
  set(tokens: TokenPair | null) {
    cache = tokens
    write(tokens)
    emit(tokens)
  },

  /** Fresh read from storage (useful if you suspect out-of-band changes). */
  reload(): TokenPair | null {
    cache = read()
    return cache
  },

  /** Is there an access token at all? (not checking expiry) */
  isAuthed(): boolean {
    return !!cache?.access_token
  },

  /** Convenience accessor */
  getAccessToken(): string | null {
    return cache?.access_token ?? null
  },

  /** Decode JWT exp and check expiry (no clock skew handling here). */
  isExpired(nowSec = Math.floor(Date.now() / 1000)): boolean {
    const exp = decodeExp(cache?.access_token)
    return exp !== null ? nowSec >= exp : true
  },

  /** Will expire within `thresholdSec` (default 60s). */
  willExpireSoon(thresholdSec = 60, nowSec = Math.floor(Date.now() / 1000)): boolean {
    const exp = decodeExp(cache?.access_token)
    return exp !== null ? exp - nowSec <= thresholdSec : true
  },

  logout() {
    authStore.set(null)
  },

  /** Subscribe to changes (pairs well with useSyncExternalStore). */
  subscribe(fn: (tokens: TokenPair | null) => void): () => void {
    const onCustom = (e: Event) => fn((e as CustomEvent<TokenPair | null>).detail ?? null)
    const onStorage = (e: StorageEvent) => {
      if (e.key === KEY) {
        cache = read()
        fn(cache)
      }
    }

    window.addEventListener(EVT, onCustom as EventListener)
    window.addEventListener("storage", onStorage)
    return () => {
      window.removeEventListener(EVT, onCustom as EventListener)
      window.removeEventListener("storage", onStorage)
    }
  },
}

// --- helpers ---
function decodeExp(jwt?: string): number | null {
  if (!jwt) return null
  const parts = jwt.split(".")
  if (parts.length < 2) return null
  try {
    const json = JSON.parse(atob(base64urlToBase64(parts[1])))
    const exp = typeof json?.exp === "number" ? json.exp : null
    return exp ?? null
  } catch {
    return null
  }
}

function base64urlToBase64(s: string) {
  return s.replace(/-/g, "+").replace(/_/g, "/") + "==".slice((2 - ((s.length * 3) % 4)) % 4)
}
