import * as React from "react"

const MOBILE_BREAKPOINT = 768
const MOBILE_QUERY = `(max-width: ${MOBILE_BREAKPOINT - 1}px)`

function subscribe(onStoreChange: () => void) {
  const mql = window.matchMedia(MOBILE_QUERY)
  mql.addEventListener("change", onStoreChange)
  return () => mql.removeEventListener("change", onStoreChange)
}

function getSnapshot() {
  return window.innerWidth < MOBILE_BREAKPOINT
}

// There is no window while server-rendering or pre-hydration. The previous
// implementation started as `undefined` and coerced to false, so keep false.
function getServerSnapshot() {
  return false
}

// matchMedia is an external store, so subscribing to it with
// useSyncExternalStore is the direct expression of what this hook does.
// Reading it into state from an effect instead cost an extra render on mount
// and reported the wrong value for that first paint.
export function useIsMobile() {
  return React.useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot)
}
