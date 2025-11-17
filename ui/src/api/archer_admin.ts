import { withRefresh } from "@/api/with-refresh.ts"
import type { AdminListArcherJobsRequest } from "@/sdk"
import { makeArcherAdminApi } from "@/sdkClient.ts"

const archerAdmin = makeArcherAdminApi()

export const archerAdminApi = {
  listJobs: (params: AdminListArcherJobsRequest = {}) => {
    return withRefresh(async () => {
      return await archerAdmin.adminListArcherJobs(params)
    })
  },
  enqueue: (body: {
    queue: string
    type: string
    payload?: object | undefined
    run_at?: string
  }) => {
    return withRefresh(async () => {
      return await archerAdmin.adminEnqueueArcherJob({ body })
    })
  },
  retryJob: (id: string) => {
    return withRefresh(async () => {
      return await archerAdmin.adminRetryArcherJob({ id })
    })
  },
  cancelJob: (id: string) => {
    return withRefresh(async () => {
      return await archerAdmin.adminCancelArcherJob({ id })
    })
  },
  listQueues: () => {
    return withRefresh(async () => {
      return await archerAdmin.adminListArcherQueues()
    })
  },
}
