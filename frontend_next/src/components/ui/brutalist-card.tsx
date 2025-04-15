import React from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { cn } from "@/lib/utils";

const brutalistCardVariants = cva("brutalist-card", {
  variants: {
    variant: {
      default: "",
      danger: "danger",
      success: "success",
      warning: "warning",
    },
  },
  defaultVariants: {
    variant: "default",
  },
});

export interface BrutalistCardProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof brutalistCardVariants> {
  asChild?: boolean;
}

const BrutalistCard = React.forwardRef<HTMLDivElement, BrutalistCardProps>(
  ({ className, variant, ...props }, ref) => {
    return (
      <div
        ref={ref}
        className={cn(brutalistCardVariants({ variant }), className)}
        {...props}
      />
    );
  }
);
BrutalistCard.displayName = "BrutalistCard";

const BrutalistCardHeader = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn("brutalist-card-header", className)}
    {...props}
  />
));
BrutalistCardHeader.displayName = "BrutalistCardHeader";

const BrutalistCardTitle = React.forwardRef<
  HTMLParagraphElement,
  React.HTMLAttributes<HTMLHeadingElement>
>(({ className, ...props }, ref) => (
  <h3
    ref={ref}
    className={cn("brutalist-card-title", className)}
    {...props}
  />
));
BrutalistCardTitle.displayName = "BrutalistCardTitle";

const BrutalistCardDescription = React.forwardRef<
  HTMLParagraphElement,
  React.HTMLAttributes<HTMLParagraphElement>
>(({ className, ...props }, ref) => (
  <p
    ref={ref}
    className={cn("brutalist-card-description", className)}
    {...props}
  />
));
BrutalistCardDescription.displayName = "BrutalistCardDescription";

const BrutalistCardContent = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn("brutalist-card-content", className)}
    {...props}
  />
));
BrutalistCardContent.displayName = "BrutalistCardContent";

const BrutalistCardFooter = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn("brutalist-card-footer", className)}
    {...props}
  />
));
BrutalistCardFooter.displayName = "BrutalistCardFooter";

export {
  BrutalistCard,
  BrutalistCardHeader,
  BrutalistCardTitle,
  BrutalistCardDescription,
  BrutalistCardContent,
  BrutalistCardFooter,
  brutalistCardVariants,
}; 