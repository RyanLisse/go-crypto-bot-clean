import React from "react";
import { VariantProps, cva } from "class-variance-authority";
import { cn } from "@/lib/utils";

const brutalistButtonVariants = cva(
  "brutalist-element inline-flex items-center justify-center whitespace-nowrap border-2 border-black font-medium transition-all focus-visible:outline-none disabled:opacity-50",
  {
    variants: {
      variant: {
        default: "bg-white text-black hover:bg-black hover:text-white",
        primary: "bg-[var(--brutalist-primary)] text-white border-black hover:bg-[var(--brutalist-primary-hover)]", 
        secondary: "bg-[var(--brutalist-secondary)] text-white border-black hover:bg-[var(--brutalist-secondary-hover)]",
        success: "bg-[var(--brutalist-success)] text-white border-black hover:bg-[var(--brutalist-success-hover)]",
        warning: "bg-[var(--brutalist-warning)] text-black border-black hover:bg-[var(--brutalist-warning-hover)]",
        danger: "bg-[var(--brutalist-danger)] text-white border-black hover:bg-[var(--brutalist-danger-hover)]",
        ghost: "bg-transparent border-black text-black hover:bg-gray-100",
        link: "bg-transparent border-transparent text-black underline-offset-4 hover:underline",
      },
      size: {
        default: "h-12 px-6 py-3",
        sm: "h-9 px-3 py-2 text-sm",
        lg: "h-14 px-8 py-4 text-lg",
        icon: "h-10 w-10",
      },
      shadow: {
        none: "",
        default: "brutalist-shadow",
        sm: "brutalist-shadow-sm",
        lg: "brutalist-shadow-lg"
      },
      transform: {
        none: "",
        default: "active:translate-x-[4px] active:translate-y-[4px] active:shadow-none",
        sm: "active:translate-x-[2px] active:translate-y-[2px] active:shadow-none",
        lg: "active:translate-x-[6px] active:translate-y-[6px] active:shadow-none",
      }
    },
    defaultVariants: {
      variant: "default",
      size: "default",
      shadow: "default",
      transform: "default"
    },
  }
);

export interface BrutalistButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof brutalistButtonVariants> {
  asChild?: boolean;
}

/**
 * BrutalistButton - A button component with a brutalist design aesthetic
 * 
 * Features thick borders, sharp edges, and interactive shadows that create
 * a tactile, physical feeling when pressed.
 */
const BrutalistButton = React.forwardRef<HTMLButtonElement, BrutalistButtonProps>(
  ({ className, variant, size, shadow, transform, ...props }, ref) => {
    return (
      <button
        className={cn(brutalistButtonVariants({ variant, size, shadow, transform, className }))}
        ref={ref}
        {...props}
      />
    );
  }
);

BrutalistButton.displayName = "BrutalistButton";

export { BrutalistButton, brutalistButtonVariants }; 