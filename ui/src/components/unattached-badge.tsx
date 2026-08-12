import { Badge } from "@/components/ui/badge.tsx"
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip.tsx"

/**
 * Marks an SSH key that no server references.
 *
 * `server_count` is optional in the API: responses describing key material
 * rather than attachment (reveal, download) omit it entirely. Absent must not
 * read as unattached, so this tests for exactly zero — `!serverCount` would
 * badge every key returned by those endpoints.
 */
export function UnattachedBadge({ serverCount }: { serverCount?: number }) {
  if (serverCount !== 0) return null

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <Badge variant="outline" className="whitespace-nowrap">
          Unattached
        </Badge>
      </TooltipTrigger>
      <TooltipContent className="max-w-xs">
        <p>
          No server in autoglue uses this key. That is not the same as unused — it may still be in
          an authorized_keys on a host autoglue does not track.
        </p>
      </TooltipContent>
    </Tooltip>
  )
}
