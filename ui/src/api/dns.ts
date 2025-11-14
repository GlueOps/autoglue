import { withRefresh } from "@/api/with-refresh.ts"
import type {
  DtoCreateDomainRequest,
  DtoCreateRecordSetRequest,
  DtoUpdateDomainRequest,
  DtoUpdateRecordSetRequest,
} from "@/sdk"
import { makeDnsApi } from "@/sdkClient.ts"

const dns = makeDnsApi()

export const dnsApi = {
  listDomains: () =>
    withRefresh(async () => {
      return await dns.listDomains()
    }),
  getDomain: (id: string) =>
    withRefresh(async () => {
      return await dns.getDomain({ id })
    }),
  createDomain: async (body: DtoCreateDomainRequest) =>
    withRefresh(async () => {
      return await dns.createDomain({ body })
    }),
  updateDomain: async (id: string, body: DtoUpdateDomainRequest) =>
    withRefresh(async () => {
      return await dns.updateDomain({ id, body })
    }),
  deleteDomain: async (id: string) =>
    withRefresh(async () => {
      return await dns.deleteDomain({ id })
    }),
  listRecordSetsByDomain: async (domainId: string) =>
    withRefresh(async () => {
      return await dns.listRecordSets({ domainId })
    }),
  createRecordSetsByDomain: async (domainId: string, body: DtoCreateRecordSetRequest) =>
    withRefresh(async () => {
      return await dns.createRecordSet({ domainId, body })
    }),
  updateRecordSetsByDomain: async (id: string, body: DtoUpdateRecordSetRequest) =>
    withRefresh(async () => {
      return await dns.updateRecordSet({ id, body })
    }),
  deleteRecordSetsByDomain: async (id: string) =>
    withRefresh(async () => {
      return await dns.deleteRecordSet({ id })
    }),
}