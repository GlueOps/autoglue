import { type ComponentType } from "react"
import { motion } from "framer-motion"
import { Monitor, Moon, Sun } from "lucide-react"
import { useTheme } from "next-themes"

import { cn } from "@/lib/utils"

type ThemeValue = "light" | "dark" | "system"

const options: { id: ThemeValue; icon: ComponentType<{ className?: string }>; label: string }[] = [
  { id: "light", icon: Sun, label: "Light" },
  { id: "dark", icon: Moon, label: "Dark" },
  { id: "system", icon: Monitor, label: "System" },
]

interface ThemePillSwitcherProps {
  className?: string
  variant?: "pill" | "wide"
  ariaLabel?: string
}

export const ThemePillSwitcher = ({
  className = "",
  variant = "pill",
  ariaLabel = "Toggle theme",
}: ThemePillSwitcherProps) => {
  const { theme, setTheme } = useTheme()

  const currentTheme = (theme ?? "system") as ThemeValue
  const isPill = variant === "pill"
  return (
    <div
      className={cn(
        "inline-flex items-center",
        isPill && "bg-muted/70 rounded-full p-1 text-xs shadow-sm",
        !isPill && "gap-2",
        className
      )}
      aria-label={ariaLabel}
      role="radiogroup"
    >
      {options.map(({ id, icon: Icon, label }) => {
        const isActive = currentTheme === id

        return (
          <button
            key={id}
            type="button"
            role="radio"
            aria-checked={isActive}
            onClick={() => setTheme(id)}
            aria-label={isPill ? label : undefined}
            className={cn(
              "focus-visible:ring-ring focus-visible:ring-offset-background relative flex items-center justify-center rounded-full transition-colors focus-visible:ring-2 focus-visible:ring-offset-2 focus-visible:outline-none",
              isActive ? "text-foreground border" : "text-muted-foreground hover:text-foreground",

              // --- Conditional Classes ---
              // "pill" variant is a fixed 8x8 square
              isPill && "h-8 w-8",
              // "wide" variant has padding, a gap, and auto width
              !isPill && "h-8 gap-2 px-3 text-sm font-medium"
            )}
          >
            {isActive && (
              <motion.span
                layoutId="theme-switcher-pill"
                className="bg-background absolute inset-0 rounded-full shadow-sm"
                transition={{ type: "spring", stiffness: 350, damping: 26 }}
              />
            )}
            <Icon className="relative z-10 h-4 w-4" />
            {!isPill && <span className="relative z-10">{label}</span>}
          </button>
        )
      })}
    </div>
  )
}