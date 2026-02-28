import { withRefresh } from "@/api/with-refresh.ts"
import type { DtoCreateAnnotationRequest, DtoUpdateAnnotationRequest } from "@/sdk"
import { makeAnnotationsApi } from "@/sdkClient.ts"

const annotations = makeAnnotationsApi()
export const annotationsApi = {
  listAnnotations: () =>
    withRefresh(async () => {
      return await annotations.listAnnotations()
    }),
  createAnnotation: (body: DtoCreateAnnotationRequest) =>
    withRefresh(async () => {
      return await annotations.createAnnotation({
        createAnnotationRequest: body,
      })
    }),
  getAnnotation: (id: string) =>
    withRefresh(async () => {
      return await annotations.getAnnotation({ id })
    }),
  deleteAnnotation: (id: string) =>
    withRefresh(async () => {
      await annotations.deleteAnnotation({ id })
    }),
  updateAnnotation: (id: string, body: DtoUpdateAnnotationRequest) =>
    withRefresh(async () => {
      return await annotations.updateAnnotation({
        id,
        updateAnnotationRequest: body,
      })
    }),
}
