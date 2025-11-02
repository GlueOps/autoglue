import { withRefresh } from "@/api/with-refresh.ts"
import type { DtoCreateSSHRequest, DtoSshResponse, DtoSshRevealResponse } from "@/sdk"
import { makeSshApi } from "@/sdkClient.ts"

const ssh = makeSshApi()
export type SshDownloadPart = "public" | "private" | "both"

export const sshApi = {
  listSshKeys: () =>
    withRefresh(async (): Promise<DtoSshResponse[]> => {
      return await ssh.listPublicSshKeys()
    }),

  createSshKey: (body: DtoCreateSSHRequest) =>
    withRefresh(async (): Promise<DtoSshResponse> => {
      // SDK expects { body }
      return await ssh.createSSHKey({ body })
    }),

  getSshKeyById: (id: string) =>
    withRefresh(async (): Promise<DtoSshResponse> => {
      return await ssh.getSSHKey({ id })
    }),

  revealSshKeyById: (id: string) =>
    withRefresh(async (): Promise<DtoSshRevealResponse> => {
      return await ssh.getSSHKey({ id, reveal: true as any })
      // Note: TS fetch generator often models query params as part of params bag.
      // If your generated client uses a different shape, change to:
      // return await ssh.getSSHKeyRaw({ id, reveal: true }).then(r => r.value())
    }),

  deleteSshKey: (id: string) =>
    withRefresh(async (): Promise<void> => {
      await ssh.deleteSSHKey({ id })
    }),

  // 1) JSON mode: returns structured JSON with filenames & (optionally) base64 zip
  downloadJson: (id: string, part: SshDownloadPart) =>
    withRefresh(async () => {
      const url = new URL(`/api/v1/ssh/${id}/download`, window.location.origin)
      url.searchParams.set("part", part)
      url.searchParams.set("mode", "json")

      const res = await fetch(url.toString())
      if (!res.ok) throw new Error(`Download failed: ${res.statusText}`)
      return (await res.json()) as {
        id: string
        name: string | null
        fingerprint: string
        filenames: string[]
        publicKey?: string | null
        privatePEM?: string | null
        zipBase64?: string | null
      }
    }),

  // 2) Attachment mode: returns a Blob (public/private file or a .zip)
  downloadBlob: (id: string, part: SshDownloadPart) =>
    withRefresh(async (): Promise<{ filename: string; blob: Blob }> => {
      const url = new URL(`/api/v1/ssh/${id}/download`, window.location.origin)
      url.searchParams.set("part", part)

      const res = await fetch(url.toString())
      if (!res.ok) throw new Error(`Download failed: ${res.statusText}`)

      // Parse filename from Content-Disposition
      const cd = res.headers.get("Content-Disposition") || ""
      const match = /filename="([^"]+)"/i.exec(cd)
      const filename = match?.[1] ?? "ssh-key-download"

      const blob = await res.blob()
      return { filename, blob }
    }),
}
