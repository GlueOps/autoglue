import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function truncateMiddle(str: string, keep = 24) {
  if (!str || str.length <= keep * 2 + 3) return str
  return `${str.slice(0, keep)}…${str.slice(-keep)}`
}
