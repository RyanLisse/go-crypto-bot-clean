# Card

The Card component is a flexible container with rounded corners, border, and subtle shadow that provides a structured visual element for grouping related content.

## Usage

```tsx
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

export function CardDemo() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Card Title</CardTitle>
        <CardDescription>Card Description</CardDescription>
      </CardHeader>
      <CardContent>
        <p>Card Content</p>
      </CardContent>
      <CardFooter>
        <p>Card Footer</p>
      </CardFooter>
    </Card>
  )
}
```

## Features

- Modular structure with separate components for different card sections
- Consistent spacing and typography
- Built with React and fully typed with TypeScript
- Customizable through className props
- Semantic HTML structure
- Cohesive styling with your application's design system

## API Reference

### Card

The root container component that wraps all card elements.

| Prop | Type | Description |
|------|------|-------------|
| `className` | `string` | Additional CSS classes to apply |
| `...props` | `React.HTMLAttributes<HTMLDivElement>` | All standard div attributes |

### CardHeader

Container for the card title and description with appropriate spacing.

| Prop | Type | Description |
|------|------|-------------|
| `className` | `string` | Additional CSS classes to apply |
| `...props` | `React.HTMLAttributes<HTMLDivElement>` | All standard div attributes |

### CardTitle

The main heading of the card.

| Prop | Type | Description |
|------|------|-------------|
| `className` | `string` | Additional CSS classes to apply |
| `...props` | `React.HTMLAttributes<HTMLHeadingElement>` | All standard heading attributes |

### CardDescription

Supplementary text that provides more context about the card content.

| Prop | Type | Description |
|------|------|-------------|
| `className` | `string` | Additional CSS classes to apply |
| `...props` | `React.HTMLAttributes<HTMLParagraphElement>` | All standard paragraph attributes |

### CardContent

The main content area of the card.

| Prop | Type | Description |
|------|------|-------------|
| `className` | `string` | Additional CSS classes to apply |
| `...props` | `React.HTMLAttributes<HTMLDivElement>` | All standard div attributes |

### CardFooter

Container for actions or supplementary information at the bottom of the card.

| Prop | Type | Description |
|------|------|-------------|
| `className` | `string` | Additional CSS classes to apply |
| `...props` | `React.HTMLAttributes<HTMLDivElement>` | All standard div attributes |

## Styling

- Card: `rounded-lg border bg-card text-card-foreground shadow-sm`
- CardHeader: `flex flex-col space-y-1.5 p-6`
- CardTitle: `text-2xl font-semibold leading-none tracking-tight`
- CardDescription: `text-sm text-muted-foreground`
- CardContent: `p-6 pt-0`
- CardFooter: `flex items-center p-6 pt-0`

## Examples

### Basic Card

```tsx
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

export function BasicCard() {
  return (
    <Card className="w-[350px]">
      <CardHeader>
        <CardTitle>Notification</CardTitle>
        <CardDescription>You have a new message.</CardDescription>
      </CardHeader>
      <CardContent>
        <p>The contents of the notification go here.</p>
      </CardContent>
      <CardFooter>
        <p className="text-xs text-muted-foreground">Received 2 hours ago</p>
      </CardFooter>
    </Card>
  )
}
```

### Interactive Card with Actions

```tsx
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

export function CardWithActions() {
  return (
    <Card className="w-[350px]">
      <CardHeader>
        <CardTitle>Account Settings</CardTitle>
        <CardDescription>Update your account preferences.</CardDescription>
      </CardHeader>
      <CardContent>
        <form>
          <div className="grid w-full items-center gap-4">
            {/* Form fields would go here */}
          </div>
        </form>
      </CardContent>
      <CardFooter className="flex justify-between">
        <Button variant="outline">Cancel</Button>
        <Button>Save</Button>
      </CardFooter>
    </Card>
  )
}
```

### Card Grid Layout

```tsx
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"

export function CardGrid() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {[1, 2, 3].map((item) => (
        <Card key={item}>
          <CardHeader>
            <CardTitle>Card {item}</CardTitle>
            <CardDescription>Card description {item}</CardDescription>
          </CardHeader>
          <CardContent>
            <p>Content for card {item}</p>
          </CardContent>
          <CardFooter>
            <p>Card footer {item}</p>
          </CardFooter>
        </Card>
      ))}
    </div>
  )
}
```

## Accessibility

- Use appropriate heading levels within CardTitle for proper document outline
- Maintain sufficient color contrast between card background and content
- Consider card interactions for keyboard users when adding interactive elements

## Design Guidelines

- Use cards to group related content and provide visual separation
- Maintain consistent padding and spacing within cards throughout the application
- Consider appropriate width constraints (e.g., `max-width`) to prevent overly wide cards
- Use CardDescription to provide additional context when the purpose of a card may not be immediately obvious
- Keep content concise and focused within each card 