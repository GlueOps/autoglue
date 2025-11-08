export const metaApi = {
  footer: async () => {
    const res = await fetch("/api/v1/version", { cache: "no-store" })
    if (!res.ok) throw new Error("failed to fetch version")
    return (await res.json()) as {
      built: string
      builtBy: string
      commit: string
      go: string
      goArch: string
      goOS: string
      version: string
    }
  },
}
