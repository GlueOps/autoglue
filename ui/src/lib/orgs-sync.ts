export const ACTIVE_ORG_KEY = "active_org_id"
export const EVT_ACTIVE_ORG_CHANGED = "active-org-changed"
export const EVT_ORGS_CHANGED = "orgs-changed"

export function getActiveOrgId(): string | null {
  return localStorage.getItem(ACTIVE_ORG_KEY)
}

export function setActiveOrgId(id: string | null) {
  if (id) localStorage.setItem(ACTIVE_ORG_KEY, id)
  else localStorage.removeItem(ACTIVE_ORG_KEY)
  window.dispatchEvent(new CustomEvent<string | null>(EVT_ACTIVE_ORG_CHANGED, { detail: id }))
}

export function emitOrgsChanged() {
  window.dispatchEvent(new Event(EVT_ORGS_CHANGED))
}
