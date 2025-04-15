import * as React from "react"
import { Slot } from "@radix-ui/react-slot"
import { cva, type VariantProps } from "class-variance-authority"
import { cn } from "@/lib/utils"

const brutalistButtonVariants = cva(
  "inline-flex items-center justify-center whitespace-nowrap rounded-md font-mono text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50",
  {
    variants: {
      variant: {
        default: "bg-brutalist-bg-light text-brutalist-text-dark hover:bg-brutalist-bg-light/90",
        danger: "bg-brutalist-danger text-brutalist-text-light hover:bg-brutalist-danger/90",
        success: "bg-brutalist-success text-brutalist-text-light hover:bg-brutalist-success/90",
        warning: "bg-brutalist-warning text-brutalist-text-dark hover:bg-brutalist-warning/90",
        outline: "border-2 border-brutalist-border bg-transparent hover:bg-brutalist-bg-light/10",
      },
      size: {
        default: "h-10 px-4 py-2",
        sm: "h-8 px-3 py-1 text-xs",
        lg: "h-12 px-6 py-3 text-base",
        xl: "h-14 px-8 py-4 text-lg",
      },
      shadow: {
        none: "",
        sm: "shadow-sm",
        default: "shadow-md",
        lg: "shadow-lg",
        xl: "shadow-xl",
      },
      transform: {
        none: "",
        sm: "hover:-translate-y-0.5 active:translate-y-0.5",
        default: "hover:-translate-y-1 active:translate-y-1",
        lg: "hover:-translate-y-2 active:translate-y-2",
      },
      borderWidth: {
        none: "border-0",
        sm: "border-2",
        default: "border-4",
        lg: "border-[6px]",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
      shadow: "default",
      transform: "default",
      borderWidth: "none",
    },
  }
)

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof brutalistButtonVariants> {
  asChild?: boolean
}

const BrutalistButton = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, shadow, transform, borderWidth, asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot : "button"
    return (
      <Comp
        className={cn(
          brutalistButtonVariants({ 
            variant, 
            size, 
            shadow, 
            transform, 
            borderWidth, 
            className 
          }),
          variant !== 'outline' && borderWidth !== 'none' && "border-brutalist-border-dark",
          "transition-all duration-150"
        )}
        ref={ref}
        {...props}
      />
    )
  }
)
BrutalistButton.displayName = "BrutalistButton"

export { BrutalistButton, brutalistButtonVariants } 