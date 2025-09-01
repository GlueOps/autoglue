import {
    BoxesIcon,
    BrainCogIcon,
    Building2Icon,
    BuildingIcon,
    ComponentIcon,
    FileKey2Icon,
    HomeIcon,
    KeyIcon,
    ListTodoIcon,
    LockKeyholeIcon,
    ServerIcon,
    SettingsIcon,
    SprayCanIcon,
    TagsIcon,
    UsersIcon,
} from "lucide-react";
import { AiOutlineCluster } from "react-icons/ai";

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
                to: "/core/cluster",
                icon: AiOutlineCluster,
            },
            {
                label: "Node Pools",
                icon: BoxesIcon,
                to: "/core/node-pools",
            },
            {
                label: "Labels",
                icon: TagsIcon,
                to: "/core/labels",
            },
            {
                label: "Roles",
                icon: ComponentIcon,
                to: "/core/roles",
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
                icon: Building2Icon,
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
                ],
            },
        ],
    },
];
