export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "" // e.g. "http://127.0.0.1:8080"

export class ApiError extends Error {
  status: number
  body?: unknown
  constructor(status: number, message: string, body?: unknown) {
    super(message)
    this.status = status
    this.body = body
  }
}

type Method = "GET" | "POST" | "PUT" | "PATCH" | "DELETE"

type RequestOpts = Omit<RequestInit, "method" | "body" | "headers"> & {
  auth?: boolean
  headers?: HeadersInit
}

function normalizeHeaders(h?: HeadersInit): Record<string, string> {
  const out: Record<string, string> = {}
  if (!h) return out
  if (h instanceof Headers) {
    h.forEach((v, k) => (out[k] = v))
  } else if (Array.isArray(h)) {
    for (const [k, v] of h) out[k] = v
  } else {
    Object.assign(out, h)
  }
  return out
}

function authHeaders(): Record<string, string> {
  const headers: Record<string, string> = {}
  const token = localStorage.getItem("access_token")
  if (token) headers.Authorization = `Bearer ${token}`
  return headers
}

function orgContextHeaders(): Record<string, string> {
  const id = localStorage.getItem("active_org_id")
  return id ? { "X-Org-ID": id } : {}
}

async function request<T>(
  path: string,
  method: Method,
  body?: unknown,
  opts: RequestOpts = {}
): Promise<T> {
  const baseHeaders: Record<string, string> = {
    "Content-Type": "application/json",
  }

  const merged: Record<string, string> = {
    ...baseHeaders,
    ...(opts.auth === false ? {} : authHeaders()),
    ...orgContextHeaders(),
    ...normalizeHeaders(opts.headers),
  }

  const res = await fetch(`${API_BASE_URL}${path}`, {
    method,
    headers: merged,
    body: body === undefined ? undefined : JSON.stringify(body),
    ...opts,
  })

  const ct = res.headers.get("content-type") || ""
  const isJSON = ct.includes("application/json")
  const payload = isJSON
    ? await res.json().catch(() => undefined)
    : await res.text().catch(() => "")

  if (!res.ok) {
    const msg =
      (isJSON &&
        payload &&
        typeof payload === "object" &&
        "error" in (payload as any) &&
        (payload as any).error) ||
      (isJSON &&
        payload &&
        typeof payload === "object" &&
        "message" in (payload as any) &&
        (payload as any).message) ||
      (typeof payload === "string" && payload) ||
      `HTTP ${res.status}`
    throw new ApiError(res.status, String(msg), payload)
  }

  console.debug("API ->", method, `${API_BASE_URL}${path}`, merged)

  return isJSON ? (payload as T) : (undefined as T)
}

export const api = {
  get: <T>(path: string, opts?: RequestOpts) => request<T>(path, "GET", undefined, opts),
  post: <T>(path: string, body?: unknown, opts?: RequestOpts) =>
    request<T>(path, "POST", body, opts),
  put: <T>(path: string, body?: unknown, opts?: RequestOpts) => request<T>(path, "PUT", body, opts),
  patch: <T>(path: string, body?: unknown, opts?: RequestOpts) =>
    request<T>(path, "PATCH", body, opts),
  delete: <T>(path: string, opts?: RequestOpts) => request<T>(path, "DELETE", undefined, opts),
}
