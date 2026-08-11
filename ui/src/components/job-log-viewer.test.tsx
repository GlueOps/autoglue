import { JobLogViewer } from "@/components/job-log-viewer"
import type { DtoJobLogPage } from "@/sdk"
import { act, render, screen, waitFor } from "@testing-library/react"
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest"

type Chunk = { id: number; chunk: string; stream?: "stdout" | "system" }

/**
 * A stand-in for the real endpoint, which returns only chunks with `id > after`
 * and reports `done` once the reader has caught up on finished work.
 *
 * Modelling the cursor rather than replaying a fixed page matters: a mock that
 * ignores `after` re-serves the same chunks on every poll, which would make a
 * duplicate-appending bug in the component look like passing behaviour.
 */
function fakeLogApi(chunks: Chunk[], opts: { finished?: boolean; pageSize?: number } = {}) {
  const { finished = true, pageSize = 100 } = opts
  return vi.fn(async (after: number): Promise<DtoJobLogPage> => {
    const remaining = chunks.filter((c) => c.id > after)
    const items = remaining.slice(0, pageSize)
    return {
      items: items.map((c) => ({
        id: c.id,
        chunk: c.chunk,
        stream: c.stream ?? "stdout",
        created_at: new Date(),
      })),
      next_cursor: items.length ? items[items.length - 1].id : after,
      // `done` is withheld until the reader has caught up, mirroring the API.
      done: finished && items.length < pageSize,
    } as DtoJobLogPage
  })
}

describe("JobLogViewer", () => {
  // shouldAdvanceTime keeps testing-library's waitFor usable with fake timers.
  // Extra polls it causes are harmless here because fakeLogApi honours the
  // cursor and returns nothing new.
  beforeEach(() => vi.useFakeTimers({ shouldAdvanceTime: true }))
  afterEach(() => vi.useRealTimers())

  it("renders chunks from the first page", async () => {
    const fetchPage = fakeLogApi([{ id: 1, chunk: "hello\nworld\n" }])

    render(<JobLogViewer fetchPage={fetchPage} />)

    expect(await screen.findByText("hello")).toBeTruthy()
    expect(screen.getByText("world")).toBeTruthy()
    expect(fetchPage).toHaveBeenCalledWith(0)
  })

  it("advances the cursor instead of refetching from zero", async () => {
    // The whole point of the cursor: a long bootstrap must not re-download its
    // entire transcript on every poll.
    const fetchPage = fakeLogApi(
      [
        { id: 7, chunk: "first\n" },
        { id: 12, chunk: "second\n" },
      ],
      { pageSize: 1 },
    )

    render(<JobLogViewer fetchPage={fetchPage} />)

    await screen.findByText("first")
    await act(async () => {
      await vi.advanceTimersByTimeAsync(1600)
    })

    await waitFor(() => expect(screen.getByText("second")).toBeTruthy())
    // The invariant that matters: the first poll starts at 0 and every later
    // one resumes from a cursor it was given, never refetching from scratch.
    const cursors = fetchPage.mock.calls.map((c) => c[0] as number)
    expect(cursors[0]).toBe(0)
    expect(cursors).toContain(7)
    expect(cursors.slice(1).every((c) => c > 0)).toBe(true)
    expect(cursors).toEqual([...cursors].sort((a, b) => a - b))
  })

  it("stops polling once the API reports done", async () => {
    const fetchPage = fakeLogApi([{ id: 1, chunk: "only\n" }])

    render(<JobLogViewer fetchPage={fetchPage} />)
    await screen.findByText("only")

    const callsAfterFirst = fetchPage.mock.calls.length
    await act(async () => {
      await vi.advanceTimersByTimeAsync(10_000)
    })

    expect(fetchPage.mock.calls.length).toBe(callsAfterFirst)
    expect(screen.getByText("Finished")).toBeTruthy()
  })

  it("keeps polling while the job is still running", async () => {
    const fetchPage = fakeLogApi([{ id: 1, chunk: "tick\n" }], { finished: false })

    render(<JobLogViewer fetchPage={fetchPage} />)
    await screen.findByText("tick")

    await act(async () => {
      await vi.advanceTimersByTimeAsync(5000)
    })

    expect(fetchPage.mock.calls.length).toBeGreaterThan(1)
    expect(screen.getByText("Streaming")).toBeTruthy()
  })

  it("does not poll at all when live is false", async () => {
    const fetchPage = fakeLogApi([{ id: 1, chunk: "static\n" }], { finished: false })

    render(<JobLogViewer fetchPage={fetchPage} live={false} />)
    await screen.findByText("static")

    await act(async () => {
      await vi.advanceTimersByTimeAsync(10_000)
    })

    expect(fetchPage).toHaveBeenCalledTimes(1)
  })

  it("splits a batched chunk into lines and drops the trailing blank", async () => {
    // The API returns batched chunks, not lines. A naive renderer would show
    // one giant blob and an empty row for the trailing newline.
    const fetchPage = fakeLogApi([{ id: 1, chunk: "alpha\nbeta\ngamma\n" }])

    const { container } = render(<JobLogViewer fetchPage={fetchPage} />)
    await screen.findByText("alpha")

    const log = container.querySelector('[role="log"]')!
    expect(log.children.length).toBe(3)
    expect(screen.getByText("gamma")).toBeTruthy()
  })

  it("marks system narration differently from remote output", async () => {
    const fetchPage = fakeLogApi([
      { id: 1, chunk: "connecting to host\n", stream: "system" },
      { id: 2, chunk: "+ apt-get install\n", stream: "stdout" },
    ])

    render(<JobLogViewer fetchPage={fetchPage} />)

    const systemLine = await screen.findByText("connecting to host")
    const stdoutLine = screen.getByText("+ apt-get install")
    expect(systemLine.className).toContain("italic")
    expect(stdoutLine.className).not.toContain("italic")
  })

  it("surfaces a fetch failure without losing what it already rendered", async () => {
    const real = fakeLogApi([{ id: 1, chunk: "before failure\n" }], { finished: false })
    const fetchPage = vi
      .fn()
      .mockImplementationOnce(real)
      .mockRejectedValue(new Error("network is down"))

    render(<JobLogViewer fetchPage={fetchPage} />)
    await screen.findByText("before failure")

    await act(async () => {
      await vi.advanceTimersByTimeAsync(1600)
    })

    await waitFor(() => expect(screen.getByText("network is down")).toBeTruthy())
    // The earlier output must survive the error.
    expect(screen.getByText("before failure")).toBeTruthy()
  })

  it("shows the empty message when there is no output", async () => {
    const fetchPage = fakeLogApi([])

    render(<JobLogViewer fetchPage={fetchPage} emptyMessage="Nothing recorded." />)

    expect(await screen.findByText("Nothing recorded.")).toBeTruthy()
  })

  it("counts rendered lines", async () => {
    const fetchPage = fakeLogApi([{ id: 1, chunk: "a\nb\nc\n" }])

    render(<JobLogViewer fetchPage={fetchPage} />)
    await screen.findByText("a")

    expect(screen.getByText("3 lines")).toBeTruthy()
  })

  it("does not overlap requests when one is slow", async () => {
    // Two polls in flight against the same cursor would double-append.
    let resolveFirst: (p: DtoJobLogPage) => void = () => {}
    const rest = fakeLogApi([
      { id: 1, chunk: "first\n" },
      { id: 2, chunk: "second\n" },
    ])
    const fetchPage = vi
      .fn()
      .mockImplementationOnce(() => new Promise<DtoJobLogPage>((r) => (resolveFirst = r)))
      .mockImplementation(rest)

    render(<JobLogViewer fetchPage={fetchPage} />)

    // While the first request is outstanding, ticks must not start another.
    await act(async () => {
      await vi.advanceTimersByTimeAsync(5000)
    })
    expect(fetchPage).toHaveBeenCalledTimes(1)

    await act(async () => {
      resolveFirst({
        items: [{ id: 1, chunk: "first\n", stream: "stdout", created_at: new Date() }],
        next_cursor: 1,
        done: false,
      } as DtoJobLogPage)
      await vi.advanceTimersByTimeAsync(1600)
    })

    await waitFor(() => expect(screen.getByText("second")).toBeTruthy())
    expect(screen.getAllByText("first")).toHaveLength(1)
  })
})
