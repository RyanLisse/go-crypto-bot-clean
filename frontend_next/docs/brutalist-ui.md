# Brutalist UI Component System

## Introduction

The Brutalist UI Component System provides a collection of React components that follow the brutalist design aesthetic. Characterized by raw, exposed elements, bold geometry, and an emphasis on functionality over decoration, this system offers a distinctive and impactful visual language for web applications.

## Design Philosophy

Brutalism in UI design draws inspiration from the architectural movement of the same name. It embraces:

- **Honesty in Materials**: Raw, unadorned elements with visible structure
- **Bold Geometry**: Strong shapes, sharp edges, and pronounced shadows
- **Function Over Form**: Prioritizing usability while making a visual statement
- **Contrast**: Stark differences between elements to create visual hierarchy
- **Minimalism**: Removing unnecessary decoration to focus on core functionality

Our Brutalist UI system translates these principles into web components that are visually distinctive while maintaining accessibility and usability.

## Component Library

### BrutalistCard

A card component with pronounced borders and drop shadows, available in multiple variants.

#### Subcomponents

- `BrutalistCard`: The main container component
- `BrutalistCardHeader`: Top section of the card (typically containing title and description)
- `BrutalistCardTitle`: Card title element
- `BrutalistCardDescription`: Secondary text below the title
- `BrutalistCardContent`: Main content area of the card
- `BrutalistCardFooter`: Bottom section of the card (typically containing actions)

#### Variants

- `default`: Standard black border and shadow
- `success`: Green border and shadow, indicating positive outcomes
- `warning`: Orange border and shadow, signaling caution
- `danger`: Red border and shadow, highlighting critical actions

#### Example Usage

```tsx
import {
  BrutalistCard,
  BrutalistCardHeader,
  BrutalistCardTitle,
  BrutalistCardDescription,
  BrutalistCardContent,
  BrutalistCardFooter
} from "@/components/ui/brutalist-card";

export function MyComponent() {
  return (
    <BrutalistCard>
      <BrutalistCardHeader>
        <BrutalistCardTitle>Card Title</BrutalistCardTitle>
        <BrutalistCardDescription>Card description text</BrutalistCardDescription>
      </BrutalistCardHeader>
      <BrutalistCardContent>
        <p>Main content goes here</p>
      </BrutalistCardContent>
      <BrutalistCardFooter>
        <BrutalistButton>Action</BrutalistButton>
      </BrutalistCardFooter>
    </BrutalistCard>
  );
}
```

### BrutalistButton

A button component with bold borders, sharp corners, and interactive shadows.

#### Variants

- `default`: Standard black border and shadow
- `destructive`: Red background with matching shadow
- `outline`: Transparent background with black border
- `secondary`: Gray background and border
- `ghost`: No border or background until hovered
- `link`: Underlined text without borders or shadows

#### Sizes

- `default`: Standard size
- `sm`: Smaller button
- `lg`: Larger button
- `icon`: Square button for icons

#### Example Usage

```tsx
import { BrutalistButton } from "@/components/ui/brutalist-button";
import { Trash } from "lucide-react";

export function ButtonExample() {
  return (
    <div className="space-y-4">
      <BrutalistButton>Default Button</BrutalistButton>
      
      <BrutalistButton variant="destructive">
        <Trash className="h-4 w-4 mr-2" />
        Delete Item
      </BrutalistButton>
      
      <BrutalistButton variant="outline" size="lg">
        Large Outline Button
      </BrutalistButton>
      
      <BrutalistButton size="icon">
        <Trash className="h-4 w-4" />
      </BrutalistButton>
    </div>
  );
}
```

## Utility Classes

The Brutalist UI system also provides CSS utility classes that can be applied to any element:

- `.brutalist-card`: Applies brutalist card styling
- `.typewriter`: Applies a monospace typewriter font style
- `.brutal-btn`: Applies brutalist button styling
- `.brutal-border`: Applies a bold black border

## Dark Mode Support

All Brutalist UI components support dark mode through the `.dark` class applied to a parent element. In dark mode, the color scheme is inverted, with light borders on dark backgrounds.

## CSS Variables

The Brutalist UI system uses a set of CSS variables that can be customized:

```css
/* Light mode */
:root {
  --brutal-border-color: black;
  --brutal-shadow-color: black;
  --brutal-background: white;
}

/* Dark mode */
.dark {
  --brutal-border-color: white;
  --brutal-shadow-color: white;
  --brutal-background: #121212;
}
```

## Accessibility Considerations

While brutalist design emphasizes visual impact, our components maintain accessibility standards:

- All components maintain proper color contrast ratios
- Interactive elements have appropriate hover and focus states
- Semantic HTML is used throughout the component structure
- Components support keyboard navigation

## Best Practices

### Do:

- Use brutalist components as focal points in your UI
- Maintain adequate spacing between brutalist elements
- Consider using brutalist styles for primary actions and important information
- Combine with more subdued elements to create visual hierarchy

### Don't:

- Overuse brutalist elements (they compete for attention)
- Sacrifice usability for visual impact
- Use brutalist styles for complex forms or dense information displays
- Mix brutalist elements with conflicting design styles

## Example Implementation

For a complete implementation example, visit the brutalist examples page in your application:

```
/examples/brutalist
```

This page demonstrates all available brutalist components and their variants in context.

## Technical Implementation

The Brutalist UI system is implemented as React components using:

- TypeScript for type safety
- Tailwind CSS for styling fundamentals
- CSS custom properties for theme consistency
- Radix UI primitives for accessibility
- Class Variance Authority (CVA) for variant management

## Contributing

When extending the Brutalist UI system, follow these guidelines:

1. Maintain the core brutalist aesthetic principles
2. Ensure all components are fully accessible
3. Implement dark mode support
4. Document all new components and variants
5. Add examples to the brutalist examples page 