import type { VariantProps } from "class-variance-authority"
import { cva } from "class-variance-authority"

export { default as Badge } from "./Badge.vue"

export const badgeVariants = cva(
  "inline-flex items-center justify-center rounded border px-2 py-0.5 text-xs font-medium w-fit whitespace-nowrap shrink-0 [&>svg]:size-3 gap-1 [&>svg]:pointer-events-none transition-colors",
  {
    variants: {
      variant: {
        default: "border-slate-200 bg-slate-50 text-slate-600",
        secondary: "border-slate-200 bg-slate-50 text-slate-600",
        destructive: "border-rose-200 bg-rose-50 text-rose-600",
        outline: "border-slate-200 text-slate-600",
        success: "border-emerald-200 bg-emerald-50 text-emerald-600",
        warning: "border-amber-200 bg-amber-50 text-amber-700",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  },
)
export type BadgeVariants = VariantProps<typeof badgeVariants>
