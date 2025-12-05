import { withRefresh } from "@/api/with-refresh"
import type { DtoCreateLoadBalancerRequest, DtoUpdateLoadBalancerRequest } from "@/sdk"
import { makeLoadBalancerApi } from "@/sdkClient"

const loadBalancers = makeLoadBalancerApi()
export const loadBalancersApi = {
  listLoadBalancers: () =>
    withRefresh(async () => {
      return await loadBalancers.listLoadBalancers()
    }),
  getLoadBalancer: (id: string) =>
    withRefresh(async () => {
      return await loadBalancers.getLoadBalancers({ id })
    }),
  createLoadBalancer: (body: DtoCreateLoadBalancerRequest) =>
    withRefresh(async () => {
      return await loadBalancers.createLoadBalancer({
        dtoCreateLoadBalancerRequest: body,
      })
    }),
  updateLoadBalancer: (id: string, body: DtoUpdateLoadBalancerRequest) =>
    withRefresh(async () => {
      return await loadBalancers.updateLoadBalancer({
        id,
        dtoUpdateLoadBalancerRequest: body,
      })
    }),
  deleteLoadBalancer: (id: string) =>
    withRefresh(async () => {
      return await loadBalancers.deleteLoadBalancer({ id })
    }),
}
