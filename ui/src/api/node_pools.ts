import { withRefresh } from "@/api/with-refresh.ts"
import type {
  DtoAttachAnnotationsRequest,
  DtoAttachLabelsRequest,
  DtoAttachServersRequest,
  DtoAttachTaintsRequest,
  DtoCreateNodePoolRequest,
  DtoUpdateNodePoolRequest,
} from "@/sdk"
import { makeNodePoolApi } from "@/sdkClient.ts"

const nodePools = makeNodePoolApi()
export const canAttachToPool = (poolRole: string | undefined, serverRole: string | undefined) => {
  if (!poolRole) return true
  return poolRole === serverRole
}

export const nodePoolsApi = {
  listNodePools: () =>
    withRefresh(async () => {
      return await nodePools.listNodePools({})
    }),
  createNodePool: (body: DtoCreateNodePoolRequest) =>
    withRefresh(async () => {
      return await nodePools.createNodePool({ createNodePoolRequest: body })
    }),
  getNodePool: (id: string) =>
    withRefresh(async () => {
      return await nodePools.getNodePool({ id })
    }),
  deleteNodePool: (id: string) =>
    withRefresh(async () => {
      await nodePools.deleteNodePool({ id })
    }),
  updateNodePool: (id: string, body: DtoUpdateNodePoolRequest) =>
    withRefresh(async () => {
      return await nodePools.updateNodePool({ id, updateNodePoolRequest: body })
    }),
  // Servers
  listNodePoolServers: (id: string) =>
    withRefresh(async () => {
      return await nodePools.listNodePoolServers({ id })
    }),
  attachNodePoolServer: (id: string, body: DtoAttachServersRequest) =>
    withRefresh(async () => {
      return await nodePools.attachNodePoolServers({ id, attachNodePoolServersRequest: body })
    }),
  detachNodePoolServers: (id: string, serverId: string) =>
    withRefresh(async () => {
      return await nodePools.detachNodePoolServer({ id, serverId })
    }),
  // Taints
  listNodePoolTaints: (id: string) =>
    withRefresh(async () => {
      return await nodePools.listNodePoolTaints({ id })
    }),
  attachNodePoolTaints: (id: string, body: DtoAttachTaintsRequest) =>
    withRefresh(async () => {
      return await nodePools.attachNodePoolTaints({ id, attachNodePoolTaintsRequest: body })
    }),
  detachNodePoolTaints: (id: string, taintId: string) =>
    withRefresh(async () => {
      return await nodePools.detachNodePoolTaint({ id, taintId })
    }),
  // Labels
  listNodePoolLabels: (id: string) =>
    withRefresh(async () => {
      return await nodePools.listNodePoolLabels({ id })
    }),
  attachNodePoolLabels: (id: string, body: DtoAttachLabelsRequest) =>
    withRefresh(async () => {
      return await nodePools.attachNodePoolLabels({ id, attachNodePoolLabelsRequest: body })
    }),
  detachNodePoolLabels: (id: string, labelId: string) =>
    withRefresh(async () => {
      return await nodePools.detachNodePoolLabel({ id, labelId })
    }),
  // Annotations
  listNodePoolAnnotations: (id: string) =>
    withRefresh(async () => {
      return await nodePools.listNodePoolAnnotations({ id })
    }),
  attachNodePoolAnnotations: (id: string, body: DtoAttachAnnotationsRequest) =>
    withRefresh(async () => {
      return await nodePools.attachNodePoolAnnotations({ id, attachNodePoolAnnotationsRequest: body })
    }),
  detachNodePoolAnnotations: (id: string, annotationId: string) =>
    withRefresh(async () => {
      return await nodePools.detachNodePoolAnnotation({ id, annotationId })
    }),
}
