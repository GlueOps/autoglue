import { withRefresh } from "@/api/with-refresh";
import type { DtoAttachBastionRequest, DtoAttachCaptainDomainRequest, DtoAttachLoadBalancerRequest, DtoAttachRecordSetRequest, DtoCreateClusterRequest, DtoSetKubeconfigRequest, DtoUpdateClusterRequest } from "@/sdk";
import { makeClusterApi, makeClusterRunsApi } from "@/sdkClient";





const clusters = makeClusterApi()
const clusterRuns = makeClusterRunsApi()

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
        createClusterRequest: body,
      })
    }),

  updateCluster: (id: string, body: DtoUpdateClusterRequest) =>
    withRefresh(async () => {
      return await clusters.updateCluster({
        clusterID: id,
        updateClusterRequest: body,
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
        setClusterKubeconfigRequest: body,
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
        attachCaptainDomainRequest: body,
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
        attachControlPlaneRecordSetRequest: body,
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
        attachAppsLoadBalancerRequest: body,
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
        attachAppsLoadBalancerRequest: body,
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
        attachBastionServerRequest: body,
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
        attachNodePoolRequest: { node_pool_id: nodePoolID },
      })
    }),

  detachNodePool: (clusterID: string, nodePoolID: string) =>
    withRefresh(async () => {
      return await clusters.detachNodePool({ clusterID, nodePoolID })
    }),

  // --- cluster runs / actions ---
  listClusterRuns: (clusterID: string) =>
    withRefresh(async () => {
      return await clusterRuns.listClusterRuns({ clusterID })
    }),

  getClusterRun: (clusterID: string, runID: string) =>
    withRefresh(async () => {
      return await clusterRuns.getClusterRun({ clusterID, runID })
    }),

  runClusterAction: (clusterID: string, actionID: string) =>
    withRefresh(async () => {
      return await clusterRuns.runClusterAction({ clusterID, actionID })
    }),
}
