import { useEffect, useState } from "react"
import { orgStore } from "@/auth/org.ts"
import { Building2, Check, ChevronsUpDown } from "lucide-react"

import { cn } from "@/lib/utils.ts"
import { Button } from "@/components/ui/button.tsx"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command.tsx"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover.tsx"

type Org = { id: string; name: string }

export const OrgSwitcher = ({ orgs }: { orgs: Org[] }) => {
  const [open, setOpen] = useState(false)
  const [value, setValue] = useState(orgStore.get() ?? "")

  useEffect(() => {
    return orgStore.subscribe((id) => setValue(id ?? ""))
  }, [])

  const selected = orgs.find((o) => o.id === value)

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="ghost"
          className="h-9 w-full justify-between px-2"
          aria-label="Switch organization"
        >
          <span className="flex items-center gap-2 truncate">
            <Building2 className="h-4 w-4" />
            <span className="truncate">{selected?.name ?? "Select org"}</span>
          </span>
          <ChevronsUpDown className="ml-2 h-4 w-4 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-[280px] p-0" align="start">
        <Command>
          <CommandInput placeholder="Search orgs..." />
          <CommandList>
            <CommandEmpty>No orgs found.</CommandEmpty>
            <CommandGroup heading="Organizations">
              {orgs.map((org) => (
                <CommandItem
                  key={org.id}
                  value={org.id}
                  onSelect={(v) => {
                    orgStore.set(v)
                    setOpen(false)
                    // Reload page to refresh data for the new org
                    window.location.reload()
                  }}
                >
                  <Check
                    className={cn("mr-2 h-4 w-4", value === org.id ? "opacity-100" : "opacity-0")}
                  />
                  <span className="truncate">{org.name}</span>
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  )
}
