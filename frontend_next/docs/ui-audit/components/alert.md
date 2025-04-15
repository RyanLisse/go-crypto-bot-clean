# Alert Component

The Alert component is used to display important messages or feedback to users. It provides context through colors, icons, and variant styles to communicate status and severity.

## Usage

```tsx
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert"
import { Info } from "lucide-react"

export function AlertDemo() {
  return (
    <Alert>
      <Info className="h-4 w-4" />
      <AlertTitle>Information</AlertTitle>
      <AlertDescription>
        This is an informational message.
      </AlertDescription>
    </Alert>
  )
}
```

## Features

- Built with React and fully typed with TypeScript
- Two variants for communicating different statuses
- Supports custom icons for enhanced visual cues
- Accessible design with appropriate ARIA attributes
- Composable structure with title and description components
- Customizable through className props

## API Reference

### Alert

The root alert component that serves as a container.

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `variant` | `'default' \| 'destructive'` | `'default'` | Controls the visual style of the alert |
| `className` | `string` | - | Additional CSS classes to apply |
| `...props` | `React.HTMLAttributes<HTMLDivElement>` | - | All standard div attributes |

### AlertTitle

The title component for the alert that provides a concise description of the message.

| Prop | Type | Description |
|------|------|-------------|
| `className` | `string` | Additional CSS classes to apply |
| `...props` | `React.HTMLAttributes<HTMLHeadingElement>` | All standard heading attributes |

### AlertDescription

The description component that provides additional details about the alert.

| Prop | Type | Description |
|------|------|-------------|
| `className` | `string` | Additional CSS classes to apply |
| `...props` | `React.HTMLAttributes<HTMLParagraphElement>` | All standard paragraph attributes |

## Variants

### Default

The standard alert style with a neutral appearance for informational content.

```tsx
<Alert>
  <AlertTitle>Information</AlertTitle>
  <AlertDescription>
    This is a standard informational alert.
  </AlertDescription>
</Alert>
```

### Destructive

A variant with a more critical appearance for warnings, errors, or destructive actions.

```tsx
<Alert variant="destructive">
  <AlertTitle>Error</AlertTitle>
  <AlertDescription>
    Your session has expired. Please log in again.
  </AlertDescription>
</Alert>
```

## Styling

The Alert component uses Tailwind CSS for styling with the following default classes:

- Alert: `relative w-full rounded-lg border p-4 [&>svg~*]:pl-7 [&>svg+div]:translate-y-[-3px] [&>svg]:absolute [&>svg]:left-4 [&>svg]:top-4 [&>svg]:text-foreground`
- AlertTitle: `mb-1 font-medium leading-none tracking-tight`
- AlertDescription: `text-sm [&_p]:leading-relaxed`

## Examples

### With Icon

Icons can be added to provide additional visual context:

```tsx
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert"
import { AlertCircle } from "lucide-react"

export function AlertWithIcon() {
  return (
    <Alert variant="destructive">
      <AlertCircle className="h-4 w-4" />
      <AlertTitle>Warning</AlertTitle>
      <AlertDescription>
        Your account is about to reach its usage limit.
      </AlertDescription>
    </Alert>
  )
}
```

### Information Alert

For general information and notifications:

```tsx
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert"
import { Info } from "lucide-react"

export function InformationAlert() {
  return (
    <Alert>
      <Info className="h-4 w-4" />
      <AlertTitle>Information</AlertTitle>
      <AlertDescription>
        We've just released a new feature. Check it out in your dashboard.
      </AlertDescription>
    </Alert>
  )
}
```

### Success Alert

For success messages and confirmations:

```tsx
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert"
import { CheckCircle } from "lucide-react"

export function SuccessAlert() {
  return (
    <Alert className="border-green-500 text-green-700 dark:text-green-400">
      <CheckCircle className="h-4 w-4 text-green-700 dark:text-green-400" />
      <AlertTitle>Success</AlertTitle>
      <AlertDescription>
        Your changes have been successfully saved.
      </AlertDescription>
    </Alert>
  )
}
```

## Accessibility

- Uses the `role="alert"` attribute to ensure proper screen reader announcement
- Provides clear visual and semantic structure with title and description
- Maintains appropriate color contrast ratios for all variants
- Icon placement supports both visual design and assistive technology

## Design Guidelines

- Use alerts sparingly to avoid overwhelming users
- Select the appropriate variant for the message:
  - Use default for non-critical information
  - Use destructive for errors or warnings
- Provide clear, concise messaging that explains:
  - What happened
  - Why it happened (if relevant)
  - What the user should do next (if action is required)
- Include an icon when it enhances understanding of the message
- Position alerts in a consistent location in the interface
- Consider using animations for alerts that appear dynamically, but ensure they don't interfere with accessibility

## Implementation Comparison

| Aspect | Vite Implementation | Next.js Implementation | Status |
|--------|---------------------|------------------------|--------|
| Base styling | Identical | Identical | ✅ Match |
| Variants | Default, Destructive | Default, Destructive | ✅ Match |
| Component structure | Alert, Title, Description | Alert, Title, Description | ✅ Match |
| ARIA attributes | role="alert" | role="alert" | ✅ Match |
| Icon support | Yes, with similar positioning | Yes, with similar positioning | ✅ Match |
| Typography | Font size and weights match | Font size and weights match | ✅ Match |

## Recommended Actions

No significant inconsistencies were found between the Vite and Next.js implementations of the Alert component. Both provide the same functionality, API, and visual appearance. 