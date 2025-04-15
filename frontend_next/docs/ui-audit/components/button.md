# Button

The Button component is a versatile and customizable button element that supports various styles, sizes, and states. It's built on top of Radix UI's Slot primitive for composition flexibility.

## Usage

```tsx
import { Button } from "@/components/ui/button"

export function ButtonDemo() {
  return (
    <div className="flex flex-wrap gap-4">
      <Button>Default</Button>
      <Button variant="destructive">Destructive</Button>
      <Button variant="outline">Outline</Button>
      <Button variant="secondary">Secondary</Button>
      <Button variant="ghost">Ghost</Button>
      <Button variant="link">Link</Button>
    </div>
  )
}
```

## Features

- Built with React and fully typed with TypeScript
- Customizable variants and sizes
- Supports asChild pattern for enhanced composition
- Responsive design with proper focus states
- Compatible with icons and text content
- Disabled state styling

## API Reference

### Button

The main button component that provides interactive elements for user actions.

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `variant` | `'default' \| 'destructive' \| 'outline' \| 'secondary' \| 'ghost' \| 'link'` | `'default'` | Controls the visual style of the button |
| `size` | `'default' \| 'sm' \| 'lg' \| 'icon'` | `'default'` | Determines the size of the button |
| `asChild` | `boolean` | `false` | When true, the component will render as its child |
| `className` | `string` | - | Additional CSS classes to apply |
| `...props` | `ButtonHTMLAttributes<HTMLButtonElement>` | - | All standard button attributes |

## Variants

### Default
The primary button style with solid background color and high contrast.

### Destructive
Used for actions with destructive or irreversible consequences, such as delete actions.

### Outline
A button with a border and transparent background, useful for secondary actions.

### Secondary
An alternative style with a different background color than the primary button.

### Ghost
A button that only shows its background on hover, useful for toolbar actions.

### Link
Styled to appear as a text link with appropriate hover underline effect.

## Sizes

- `default`: Standard size for most use cases (h-10, px-4, py-2)
- `sm`: Smaller size for compact UIs (h-9, px-3)
- `lg`: Larger size for emphasized actions (h-11, px-8)
- `icon`: Square button for icon-only usage (h-10, w-10)

## Examples

### With Icon

```tsx
import { Button } from "@/components/ui/button"
import { PlusIcon } from "@radix-ui/react-icons"

export function IconButtonDemo() {
  return (
    <Button>
      <PlusIcon className="mr-2" />
      Add item
    </Button>
  )
}
```

### Icon-only Button

```tsx
import { Button } from "@/components/ui/button"
import { PlusIcon } from "@radix-ui/react-icons"

export function IconOnlyButtonDemo() {
  return (
    <Button size="icon" aria-label="Add item">
      <PlusIcon />
    </Button>
  )
}
```

### Loading Button

```tsx
import { useState } from "react"
import { Button } from "@/components/ui/button"
import { ReloadIcon } from "@radix-ui/react-icons"

export function LoadingButtonDemo() {
  const [isLoading, setIsLoading] = useState(false)
  
  return (
    <Button 
      onClick={() => setIsLoading(true)} 
      disabled={isLoading}
    >
      {isLoading ? (
        <>
          <ReloadIcon className="mr-2 animate-spin" />
          Loading...
        </>
      ) : (
        "Click me"
      )}
    </Button>
  )
}
```

## Accessibility

- Buttons maintain accessible contrast ratios for all variants
- Proper focus states for keyboard navigation
- Disabled styles for non-interactive states
- Icons are set to `pointer-events: none` to ensure consistent click behavior

## Design Guidelines

- Use the default variant for primary actions
- Use destructive for dangerous actions that may have irreversible consequences
- Use outline or ghost variants for secondary actions
- Maintain consistent button usage patterns throughout the application
- Add appropriate loading states for actions that aren't instantaneous 