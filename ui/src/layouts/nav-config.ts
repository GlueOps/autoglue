import type { ComponentType } from "react"
import {
  BoxesIcon,
  Building2,
  ComponentIcon,
  FileKey2Icon,
  KeyRound,
  ServerIcon,
  SprayCanIcon,
  TagsIcon,
  User2,
  Users,
} from "lucide-react"
import { AiOutlineCluster } from "react-icons/ai"
import { GrUserWorker } from "react-icons/gr"

export type NavItem = {
  to: string
  label: string
  icon: ComponentType<{ className?: string }>
}

export const mainNav: NavItem[] = [
  { to: "/clusters", label: "Clusters", icon: AiOutlineCluster },
  { to: "/node-pools", label: "Node Pools", icon: BoxesIcon },
  { to: "/annotations", label: "Annotations", icon: ComponentIcon },
  { to: "/labels", label: "Labels", icon: TagsIcon },
  { to: "/taints", label: "Taints", icon: SprayCanIcon },
  { to: "/servers", label: "Servers", icon: ServerIcon },
  { to: "/ssh", label: "SSH Keys", icon: FileKey2Icon },
]

export const orgNav: NavItem[] = [
  { to: "/org/members", label: "Members", icon: Users },
  { to: "/org/api-keys", label: "Org API Keys", icon: KeyRound },
  { to: "/org/settings", label: "Org Settings", icon: Building2 },
]

export const userNav: NavItem[] = [{ to: "/me", label: "Profile", icon: User2 }]

export const adminNav: NavItem[] = [
  { to: "/admin/users", label: "Users Admin", icon: Users },
  { to: "/admin/jobs", label: "Jobs Admin", icon: GrUserWorker },
]
