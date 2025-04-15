# Badge Component

The Badge component is a small UI element used to display statuses, categories, counts, or labels. It's typically displayed as a small pill-shaped element with minimal content.

## Usage

```tsx
import { Badge } from "@/components/ui/badge"

export function BadgeDemo() {
  return (
    <div className="flex flex-wrap gap-2">
      <Badge>Default</Badge>
      <Badge variant="secondary">Secondary</Badge>
      <Badge variant="destructive">Destructive</Badge>
      <Badge variant="outline">Outline</Badge>
    </div>
  )
}
```

## Features

- Built with React and fully typed with TypeScript
- Four variant styles for different contexts
- Appropriate styling for small textual indicators
- Consistent sizing and padding
- Focus styles for keyboard navigation
- Seamless integration with other UI components

## API Reference

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `variant` | `'default' \| 'secondary' \| 'destructive' \| 'outline'` | `'default'` | Controls the visual style of the badge |
| `className` | `string` | - | Additional CSS classes to apply |
| `...props` | `React.HTMLAttributes<HTMLDivElement>` | - | All standard div attributes |

## Variants

### Default
The primary badge style with solid background color.

```tsx
<Badge>Default</Badge>
```

### Secondary
An alternative style with a different background color.

```tsx
<Badge variant="secondary">Secondary</Badge>
```

### Destructive
Used for indicating errors, warnings, or destructive actions.

```tsx
<Badge variant="destructive">Destructive</Badge>
```

### Outline
A subtle badge with only a border and no background fill.

```tsx
<Badge variant="outline">Outline</Badge>
```

## Styling

The Badge component uses Tailwind CSS for styling with the following default classes:

- Base: `inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2`
- Default variant: `border-transparent bg-primary text-primary-foreground hover:bg-primary/80`
- Secondary variant: `border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80`
- Destructive variant: `border-transparent bg-destructive text-destructive-foreground hover:bg-destructive/80`
- Outline variant: `text-foreground`

## Examples

### Status Badges

Use badges to indicate status:

```tsx
import { Badge } from "@/components/ui/badge"

export function StatusBadges() {
  return (
    <div className="flex flex-wrap gap-2">
      <Badge variant="outline" className="border-green-500 text-green-500">Active</Badge>
      <Badge variant="outline" className="border-yellow-500 text-yellow-500">Pending</Badge>
      <Badge variant="outline" className="border-red-500 text-red-500">Closed</Badge>
    </div>
  )
}
```

### Count Indicator

Use badges to display numerical values:

```tsx
import { Badge } from "@/components/ui/badge"
import { Bell } from "lucide-react"

export function NotificationIndicator() {
  return (
    <div className="relative">
      <Bell className="h-6 w-6" />
      <Badge className="absolute -top-2 -right-2 h-5 w-5 flex items-center justify-center p-0">
        5
      </Badge>
    </div>
  )
}
```

### With Icons

Combine badges with icons for enhanced visual cues:

```tsx
import { Badge } from "@/components/ui/badge"
import { Check, X } from "lucide-react"

export function BadgesWithIcons() {
  return (
    <div className="flex flex-wrap gap-2">
      <Badge className="gap-1">
        <Check className="h-3 w-3" />
        Completed
      </Badge>
      <Badge variant="destructive" className="gap-1">
        <X className="h-3 w-3" />
        Failed
      </Badge>
    </div>
  )
}
```

## Accessibility

- Use appropriate color contrast for all badge variants
- Ensure badge content is concise and readable
- When used as a status indicator, consider including additional context for screen readers
- When badges convey important information, ensure they aren't the only means of communicating that information

## Design Guidelines

- Keep badge content short (1-2 words or numbers)
- Use consistently across the interface for similar types of information
- Choose appropriate variants based on context:
  - Use default or secondary for neutral information
  - Use destructive for errors or warnings
  - Use outline for more subtle indicators
- Consider placement carefully to ensure badges don't disrupt the overall layout
- Maintain adequate spacing between badges when displaying multiple instances

## Implementation Comparison

| Aspect | Vite Implementation | Next.js Implementation | Status |
|--------|---------------------|------------------------|--------|
| Base styling | Identical | Identical | ✅ Match |
| Variants | Default, Secondary, Destructive, Outline | Default, Secondary, Destructive, Outline | ✅ Match |
| Visual Appearance | Rounded pill shape with consistent padding | Rounded pill shape with consistent padding | ✅ Match |
| Focus states | Focus ring styling matches | Focus ring styling matches | ✅ Match |
| Typography | Text size, weight match | Text size, weight match | ✅ Match |
| HTML Element | `<div>` | `<div>` | ✅ Match |

## Recommended Actions

No significant inconsistencies were found between the Vite and Next.js implementations of the Badge component. Both provide the same functionality, API, and visual appearance. 