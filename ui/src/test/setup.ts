import { cleanup } from "@testing-library/react"
import { afterEach, vi } from "vitest"

afterEach(() => {
  cleanup()
})

// jsdom has no layout engine, so scrollHeight/clientHeight are always 0 and
// scrollTop never moves on its own. Components that follow the bottom of a
// scroll container need these to be settable, or the "am I at the bottom?"
// check is meaningless in tests.
export function stubScrollMetrics(el: HTMLElement, metrics: { scrollHeight: number; clientHeight: number }) {
  Object.defineProperty(el, "scrollHeight", { value: metrics.scrollHeight, configurable: true })
  Object.defineProperty(el, "clientHeight", { value: metrics.clientHeight, configurable: true })
}

// Silence React's act() noise from timer-driven polling; the tests assert on
// rendered output rather than on warnings.
const origError = console.error
console.error = (...args: unknown[]) => {
  if (typeof args[0] === "string" && args[0].includes("not wrapped in act")) return
  origError(...args)
}

vi.stubGlobal("scrollTo", () => {})
