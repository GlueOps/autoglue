import { render, screen } from "@testing-library/react"
import { describe, expect, it } from "vitest"

import { UnattachedBadge } from "@/components/unattached-badge"

describe("UnattachedBadge", () => {
  it("badges a key no server references", () => {
    render(<UnattachedBadge serverCount={0} />)
    expect(screen.getByText("Unattached")).toBeTruthy()
  })

  it("stays hidden for a key in use", () => {
    render(<UnattachedBadge serverCount={3} />)
    expect(screen.queryByText("Unattached")).toBeNull()
  })

  // The case that makes this component worth having. server_count is absent on
  // responses that never computed it, and `!serverCount` would treat that as
  // zero and badge every revealed or downloaded key as unattached.
  it("stays hidden when the count was never computed", () => {
    render(<UnattachedBadge />)
    expect(screen.queryByText("Unattached")).toBeNull()
  })
})
