# Dialog Component

The Dialog component is a modal window that appears in front of app content to provide critical information or require user decisions. Dialogs are purposefully interruptive, so they should be used sparingly.

## Overview

The Dialog component is built on top of Radix UI's Dialog primitive and provides a set of composable parts that can be combined to create accessible modal dialogs. The implementation includes animations, styling, and proper accessibility features.

## Usage

### Basic Dialog

```tsx
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
  DialogClose,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"

export function BasicDialog() {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline">Open Dialog</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Edit Profile</DialogTitle>
          <DialogDescription>
            Make changes to your profile here. Click save when you're done.
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          {/* Dialog content goes here */}
        </div>
        <DialogFooter>
          <Button type="submit">Save changes</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
```

### Controlled Dialog

```tsx
import { useState } from "react"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"

export function ControlledDialog() {
  const [open, setOpen] = useState(false)

  return (
    <>
      <Button onClick={() => setOpen(true)}>Open Dialog</Button>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Confirmation</DialogTitle>
            <DialogDescription>
              Are you sure you want to perform this action?
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setOpen(false)}>
              Cancel
            </Button>
            <Button onClick={() => {
              // Perform action
              setOpen(false)
            }}>
              Confirm
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </>
  )
}
```

### Alert Dialog

```tsx
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"

export function AlertDialog() {
  return (
    <Dialog>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Delete Account</DialogTitle>
          <DialogDescription>
            This action cannot be undone. This will permanently delete your account
            and remove your data from our servers.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline">Cancel</Button>
          <Button variant="destructive">Delete</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
```

## API Reference

### Dialog

The root dialog component.

```tsx
import { Dialog } from "@/components/ui/dialog"
```

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `open` | `boolean` | `undefined` | Whether the dialog is open. Use with `onOpenChange` for controlled usage. |
| `onOpenChange` | `(open: boolean) => void` | `undefined` | Callback fired when the open state changes. |
| `modal` | `boolean` | `true` | Whether the dialog is modal (blocks interaction with the rest of the page). |

### DialogTrigger

The button that opens the dialog.

```tsx
import { DialogTrigger } from "@/components/ui/dialog"
```

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `asChild` | `boolean` | `false` | When `true`, the component will render its child instead of a default button. |

### DialogContent

The content area of the dialog.

```tsx
import { DialogContent } from "@/components/ui/dialog"
```

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `className` | `string` | `""` | Additional CSS classes for styling. |
| `children` | `React.ReactNode` | `undefined` | The content of the dialog. |
| `forceMount` | `boolean` | `false` | Forces the dialog to mount even when closed. |

### DialogHeader

A layout component for the dialog header.

```tsx
import { DialogHeader } from "@/components/ui/dialog"
```

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `className` | `string` | `""` | Additional CSS classes for styling. |

### DialogTitle

The title of the dialog.

```tsx
import { DialogTitle } from "@/components/ui/dialog"
```

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `className` | `string` | `""` | Additional CSS classes for styling. |

### DialogDescription

A description for the dialog.

```tsx
import { DialogDescription } from "@/components/ui/dialog"
```

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `className` | `string` | `""` | Additional CSS classes for styling. |

### DialogFooter

A layout component for the dialog footer.

```tsx
import { DialogFooter } from "@/components/ui/dialog"
```

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `className` | `string` | `""` | Additional CSS classes for styling. |

### DialogClose

The button that closes the dialog.

```tsx
import { DialogClose } from "@/components/ui/dialog"
```

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `asChild` | `boolean` | `false` | When `true`, the component will render its child instead of a default button. |

## Styling

The Dialog component uses Tailwind CSS for styling. The default styles provide a clean, modern appearance, but can be customized as needed:

- `DialogOverlay`: Controls the overlay background and animation
- `DialogContent`: Controls the main dialog container, positioning, and animations
- `DialogHeader`: Controls the header spacing and alignment
- `DialogFooter`: Controls the footer spacing and alignment
- `DialogTitle`: Controls the title typography
- `DialogDescription`: Controls the description typography

## Accessibility

The Dialog component follows WAI-ARIA guidelines for dialog accessibility:

- It uses proper ARIA roles and attributes
- Handles focus management automatically (traps focus within dialog)
- Provides keyboard navigation support (Escape to close)
- Includes a visually hidden close button label for screen readers
- When opened, focus is automatically moved to the first focusable element in the dialog

## Best Practices

1. **Use sparingly**: Dialogs interrupt the user flow, so use them only when necessary.
2. **Keep it simple**: Limit the number of actions in a dialog.
3. **Clear titles**: Use concise, descriptive titles that clearly communicate the purpose.
4. **Descriptive buttons**: Action buttons should clearly indicate what happens when clicked.
5. **Responsive design**: Ensure the dialog works well on all screen sizes.
6. **Keyboard accessibility**: All interactions should be possible with keyboard alone.
7. **Escape key**: Always allow users to dismiss non-critical dialogs with the Escape key.

## Design Guidelines

- **Size**: Generally, dialogs should be kept as small as possible while accommodating their content. Common widths are 425px for forms and 500px for confirmations.
- **Position**: Dialogs are centered by default to draw attention.
- **Animation**: The default animations help users understand that the dialog is a temporary, modal layer.
- **Close button**: Always include a visible close mechanism, either through a dedicated close button or cancel action.

## Examples

### Form Dialog

A dialog containing a form for user input:

```tsx
import { useState } from "react"
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"

export function FormDialog() {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button>Edit Profile</Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Edit profile</DialogTitle>
          <DialogDescription>
            Make changes to your profile here. Click save when you're done.
          </DialogDescription>
        </DialogHeader>
        <form>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="name" className="text-right">
                Name
              </Label>
              <Input
                id="name"
                defaultValue="Alex Smith"
                className="col-span-3"
              />
            </div>
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="username" className="text-right">
                Username
              </Label>
              <Input
                id="username"
                defaultValue="alexsmith"
                className="col-span-3"
              />
            </div>
          </div>
          <DialogFooter>
            <Button type="submit">Save changes</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
} 