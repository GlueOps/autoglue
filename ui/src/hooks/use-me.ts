import { meApi } from "@/api/me.ts"
import { useQuery } from "@tanstack/react-query"

export function useMe() {
  return useQuery({
    queryKey: ["me"],
    queryFn: () => meApi.getMe(),
    staleTime: 5 * 60 * 1000, // cache a bit
  })
}
