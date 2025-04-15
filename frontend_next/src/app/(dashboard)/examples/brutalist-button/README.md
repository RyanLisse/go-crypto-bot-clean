# BrutalistButton Component

A bold, brutal-style button component with customizable variants, sizes, shadows, transforms, and border widths.

## Features

- **Multiple Variants**: Default, Danger, Success, Warning, and Outline styles
- **Customizable Sizes**: From small to extra large
- **Shadow Effects**: Various shadow intensities
- **Transform Animations**: Dynamic movement on hover and click
- **Border Width Options**: From none to large borders
- **Accessibility**: Supports keyboard navigation and screen readers
- **Composition**: Can be composed with other components using `asChild`

## Usage Examples

### Basic Button

```tsx
import { BrutalistButton } from "@/components/ui/brutalist-button"

export default function MyComponent() {
  return <BrutalistButton>Click Me</BrutalistButton>
}
```

### With Variants

```tsx
<BrutalistButton variant="default">Default</BrutalistButton>
<BrutalistButton variant="danger">Danger</BrutalistButton>
<BrutalistButton variant="success">Success</BrutalistButton>
<BrutalistButton variant="warning">Warning</BrutalistButton>
<BrutalistButton variant="outline">Outline</BrutalistButton>
```

### With Different Sizes

```tsx
<BrutalistButton size="sm">Small</BrutalistButton>
<BrutalistButton size="default">Default</BrutalistButton>
<BrutalistButton size="lg">Large</BrutalistButton>
<BrutalistButton size="xl">Extra Large</BrutalistButton>
```

### With Shadows

```tsx
<BrutalistButton shadow="none">No Shadow</BrutalistButton>
<BrutalistButton shadow="sm">Small Shadow</BrutalistButton>
<BrutalistButton shadow="default">Default Shadow</BrutalistButton>
<BrutalistButton shadow="lg">Large Shadow</BrutalistButton>
<BrutalistButton shadow="xl">Extra Large Shadow</BrutalistButton>
```

### With Transform Effects

```tsx
<BrutalistButton transform="none">No Transform</BrutalistButton>
<BrutalistButton transform="sm">Small Transform</BrutalistButton>
<BrutalistButton transform="default">Default Transform</BrutalistButton>
<BrutalistButton transform="lg">Large Transform</BrutalistButton>
```

### With Border Width Options

```tsx
<BrutalistButton borderWidth="none">No Border</BrutalistButton>
<BrutalistButton borderWidth="sm">Small Border</BrutalistButton>
<BrutalistButton borderWidth="default">Default Border</BrutalistButton>
<BrutalistButton borderWidth="lg">Large Border</BrutalistButton>
```

### Combined Customizations

```tsx
<BrutalistButton 
  variant="success"
  size="lg"
  shadow="lg"
  transform="lg"
  borderWidth="default"
>
  Submit
</BrutalistButton>
```

### Using with Other Components (asChild)

```tsx
import { BrutalistButton } from "@/components/ui/brutalist-button"
import Link from "next/link"

export default function MyComponent() {
  return (
    <BrutalistButton asChild>
      <Link href="/dashboard">Go to Dashboard</Link>
    </BrutalistButton>
  )
}
```

## Props Reference

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `variant` | `"default" \| "danger" \| "success" \| "warning" \| "outline"` | `"default"` | Determines the visual style of the button |
| `size` | `"default" \| "sm" \| "lg" \| "xl"` | `"default"` | Controls the size of the button |
| `shadow` | `"none" \| "sm" \| "default" \| "lg" \| "xl"` | `"default"` | Sets the shadow intensity |
| `transform` | `"none" \| "sm" \| "default" \| "lg"` | `"default"` | Controls the hover/active animation effect |
| `borderWidth` | `"none" \| "sm" \| "default" \| "lg"` | `"none"` | Sets the border width |
| `asChild` | `boolean` | `false` | Merges props onto the immediate child element |
| `className` | `string` | - | Additional CSS classes to apply |
| `[...props]` | `ButtonHTMLAttributes<HTMLButtonElement>` | - | All standard button attributes |

## Styling Guidelines

### Theme Integration

The BrutalistButton uses CSS variables for its theme colors. Make sure your `globals.css` includes these variables:

```css
:root {
  --brutalist-bg-light: #ffffff;
  --brutalist-text-dark: #000000;
  --brutalist-text-light: #ffffff;
  --brutalist-danger: #ff4d4f;
  --brutalist-success: #52c41a;
  --brutalist-warning: #faad14;
  --brutalist-border: #000000;
  --brutalist-border-dark: #000000;
}

.dark {
  --brutalist-bg-light: #1f1f1f;
  --brutalist-text-dark: #ffffff;
  --brutalist-text-light: #ffffff;
  --brutalist-border-dark: #ffffff;
}
```

### Custom Styling

You can extend the button's styling with the `className` prop:

```tsx
<BrutalistButton className="my-custom-class">
  Custom Styled Button
</BrutalistButton>
```

## Accessibility

- Supports keyboard navigation
- Maintains proper contrast ratios for text readability
- Preserves focus visibility for keyboard users
- When using `asChild` with links, ensure proper ARIA attributes are added

## Best Practices

1. Use appropriate variants to convey purpose (e.g., "danger" for destructive actions)
2. Maintain consistent styling within your application
3. Choose transform effects that complement your UI motion design
4. Consider mobile usability when selecting sizes and interactive effects
5. Use appropriate shadow effects based on your design's elevation system

## Implementation Details

The BrutalistButton is built using:
- React's `forwardRef` for proper ref forwarding
- Radix UI's `Slot` for component composition (`asChild` prop)
- Class Variance Authority (CVA) for variant management
- Tailwind CSS for styling 