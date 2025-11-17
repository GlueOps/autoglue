import { withRefresh } from "@/api/with-refresh"
import type {
  DtoAttachBastionRequest,
  DtoAttachCaptainDomainRequest,
  DtoAttachLoadBalancerRequest,
  DtoAttachRecordSetRequest,
  DtoCreateClusterRequest,
  DtoSetKubeconfigRequest,
  DtoUpdateClusterRequest,
} from "@/sdk"
import { makeClusterApi } from "@/sdkClient"

const clusters = makeClusterApi()

export const clustersApi = {
  // --- basic CRUD ---

  listClusters: (q?: string) =>
    withRefresh(async () => {
      return await clusters.listClusters(q ? { q } : {})
    }),

  getCluster: (id: string) =>
    withRefresh(async () => {
      return await clusters.getCluster({ clusterID: id })
    }),

  createCluster: (body: DtoCreateClusterRequest) =>
    withRefresh(async () => {
      return await clusters.createCluster({
        dtoCreateClusterRequest: body,
      })
    }),

  updateCluster: (id: string, body: DtoUpdateClusterRequest) =>
    withRefresh(async () => {
      return await clusters.updateCluster({
        clusterID: id,
        dtoUpdateClusterRequest: body,
      })
    }),

  deleteCluster: (id: string) =>
    withRefresh(async () => {
      return await clusters.deleteCluster({ clusterID: id })
    }),

  // --- kubeconfig ---

  setKubeconfig: (clusterID: string, body: DtoSetKubeconfigRequest) =>
    withRefresh(async () => {
      return await clusters.setClusterKubeconfig({
        clusterID,
        dtoSetKubeconfigRequest: body,
      })
    }),

  clearKubeconfig: (clusterID: string) =>
    withRefresh(async () => {
      return await clusters.clearClusterKubeconfig({ clusterID })
    }),

  // --- captain domain ---

  attachCaptainDomain: (clusterID: string, body: DtoAttachCaptainDomainRequest) =>
    withRefresh(async () => {
      return await clusters.attachCaptainDomain({
        clusterID,
        dtoAttachCaptainDomainRequest: body,
      })
    }),

  detachCaptainDomain: (clusterID: string) =>
    withRefresh(async () => {
      return await clusters.detachCaptainDomain({ clusterID })
    }),

  // --- control plane record set ---

  attachControlPlaneRecordSet: (clusterID: string, body: DtoAttachRecordSetRequest) =>
    withRefresh(async () => {
      return await clusters.attachControlPlaneRecordSet({
        clusterID,
        dtoAttachRecordSetRequest: body,
      })
    }),

  detachControlPlaneRecordSet: (clusterID: string) =>
    withRefresh(async () => {
      return await clusters.detachControlPlaneRecordSet({ clusterID })
    }),

  // --- load balancers ---

  attachAppsLoadBalancer: (clusterID: string, body: DtoAttachLoadBalancerRequest) =>
    withRefresh(async () => {
      return await clusters.attachAppsLoadBalancer({
        clusterID,
        dtoAttachLoadBalancerRequest: body,
      })
    }),

  detachAppsLoadBalancer: (clusterID: string) =>
    withRefresh(async () => {
      return await clusters.detachAppsLoadBalancer({ clusterID })
    }),

  attachGlueOpsLoadBalancer: (clusterID: string, body: DtoAttachLoadBalancerRequest) =>
    withRefresh(async () => {
      return await clusters.attachGlueOpsLoadBalancer({
        clusterID,
        dtoAttachLoadBalancerRequest: body,
      })
    }),

  detachGlueOpsLoadBalancer: (clusterID: string) =>
    withRefresh(async () => {
      return await clusters.detachGlueOpsLoadBalancer({ clusterID })
    }),

  // --- bastion ---

  attachBastion: (clusterID: string, body: DtoAttachBastionRequest) =>
    withRefresh(async () => {
      return await clusters.attachBastionServer({
        clusterID,
        dtoAttachBastionRequest: body,
      })
    }),

  detachBastion: (clusterID: string) =>
    withRefresh(async () => {
      return await clusters.detachBastionServer({ clusterID })
    }),

  // -- node-pools

  attachNodePool: (clusterID: string, nodePoolID: string) =>
    withRefresh(async () => {
      return await clusters.attachNodePool({
        clusterID,
        dtoAttachNodePoolRequest: { node_pool_id: nodePoolID },
      })
    }),

  detachNodePool: (clusterID: string, nodePoolID: string) =>
    withRefresh(async () => {
      return await clusters.detachNodePool({ clusterID, nodePoolID })
    }),
}