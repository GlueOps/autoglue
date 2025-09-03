import type { ComponentType } from "react"
import {
  BoxesIcon,
  BrainCogIcon,
  BuildingIcon,
  ComponentIcon,
  FileKey2Icon,
  HomeIcon,
  KeyIcon,
  ListTodoIcon,
  LockKeyholeIcon,
  ServerIcon,
  SettingsIcon,
  ShieldCheckIcon,
  SprayCanIcon,
  TagsIcon,
  UserIcon,
  UsersIcon,
} from "lucide-react"
import { AiOutlineCluster } from "react-icons/ai"

export type NavItem = {
  label: string
  icon: ComponentType<{ className?: string }>
  to?: string
  items?: NavItem[]
  requiresAdmin?: boolean
  requiresOrgAdmin?: boolean
}

export const items = [
  {
    label: "Dashboard",
    icon: HomeIcon,
    to: "/dashboard",
  },
  {
    label: "Core",
    icon: BrainCogIcon,
    items: [
      {
        label: "Cluster",
        to: "/core/clusters",
        icon: AiOutlineCluster,
      },
      {
        label: "Node Pools",
        icon: BoxesIcon,
        to: "/core/nodepools",
      },
      {
        label: "Annotations",
        icon: ComponentIcon,
        to: "/core/annotations",
      },
      {
        label: "Labels",
        icon: TagsIcon,
        to: "/core/labels",
      },
      {
        label: "Taints",
        icon: SprayCanIcon,
        to: "/core/taints",
      },
      {
        label: "Servers",
        icon: ServerIcon,
        to: "/core/servers",
      },
    ],
  },
  {
    label: "Security",
    icon: LockKeyholeIcon,
    items: [
      {
        label: "Keys & Tokens",
        icon: KeyIcon,
        to: "/security/keys",
      },
      {
        label: "SSH Keys",
        to: "/security/ssh",
        icon: FileKey2Icon,
      },
    ],
  },
  {
    label: "Tasks",
    icon: ListTodoIcon,
    items: [],
  },
  {
    label: "Settings",
    icon: SettingsIcon,
    items: [
      {
        label: "Organizations",
        to: "/settings/orgs",
        icon: BuildingIcon,
      },
      {
        label: "Members",
        to: "/settings/members",
        icon: UsersIcon,
      },
      {
        label: "Profile",
        to: "/settings/me",
        icon: UserIcon,
      },
    ],
  },
  {
    label: "Admin",
    icon: ShieldCheckIcon,
    requiresAdmin: true,
    items: [
      {
        label: "Users",
        to: "/admin/users",
        icon: UsersIcon,
        requiresAdmin: true,
      },
    ],
  },
]
