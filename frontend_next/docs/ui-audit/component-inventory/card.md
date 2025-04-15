# Card Component Inventory

## Basic Information
- **Component Name**: Card
- **Category**: Core Component
- **Description**: A container component that groups related content and UI elements with consistent styling and spacing.
- **Usage Frequency**: High (used across most dashboard views, portfolio sections, and asset details)

## Implementation Details

### Next.js Implementation
- **File Location**: `frontend_next/src/components/ui/Card.tsx`
- **Style Implementation**: Tailwind CSS with custom classes
- **Responsive Behavior**: Adapts to container width, maintains consistent padding
- **Animation/Transitions**: Subtle hover elevation effect
- **Variants**:
  - Default (with rounded corners and shadow)
  - Borderless (no border, just shadow)
  - Elevated (more pronounced shadow)
  - Interactive (pointer events, hover effects)

### Vite Implementation  
- **File Location**: `src/components/ui/Card/Card.tsx`
- **Style Implementation**: Styled Components with theme variables
- **Responsive Behavior**: Similar to Next.js version, fluid width
- **Animation/Transitions**: None
- **Variants**:
  - Default (with rounded corners and shadow)
  - Flat (no shadow)
  - Bordered (explicit border instead of shadow)

## Props & API

### Common Props (Both Implementations)
- `children`: ReactNode - Content to render inside the card
- `className`: string - Additional CSS classes
- `variant`: string - Card style variant

### Next.js Specific Props
- `isLoading`: boolean - Show skeleton loading state
- `hoverEffect`: boolean - Enable/disable hover elevation
- `as`: ElementType - Render as different HTML element

### Vite Specific Props
- `fullWidth`: boolean - Expand to 100% of container width
- `noPadding`: boolean - Remove default padding

## Visual Comparison

### Style Differences
- **Shadow**: 
  - Next.js: `filter: drop-shadow(0 4px 12px rgba(0, 0, 0, 0.12))`
  - Vite: `box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1)`
- **Border Radius**:
  - Next.js: 12px
  - Vite: 8px
- **Default Padding**:
  - Next.js: 16px
  - Vite: 12px
- **Hover Effects**:
  - Next.js: Subtle elevation increase
  - Vite: None

## Accessibility
- Both implementations use appropriate semantic HTML elements
- Both use sufficient color contrast for borders and backgrounds
- Next.js version includes proper ARIA attributes for skeleton loaders

## Usage Examples

### Next.js Example
```tsx
<Card className="w-full max-w-md">
  <CardHeader>
    <CardTitle>Portfolio Summary</CardTitle>
    <CardDescription>Your assets at a glance</CardDescription>
  </CardHeader>
  <CardContent>
    <p>Current value: $10,243.55</p>
  </CardContent>
  <CardFooter>
    <Button variant="outline">View Details</Button>
  </CardFooter>
</Card>
```

### Vite Example
```tsx
<Card variant="default">
  <CardHeader>
    <CardTitle>Portfolio Summary</CardTitle>
    <CardDescription>Your assets at a glance</CardDescription>
  </CardHeader>
  <CardContent>
    <p>Current value: $10,243.55</p>
  </CardContent>
  <CardFooter>
    <Button variant="outline">View Details</Button>
  </CardFooter>
</Card>
```

## Recommended Improvements
1. Standardize shadow implementation (preferably `box-shadow` for better performance)
2. Align border-radius values to design system specification (8px recommended)
3. Normalize padding across implementations
4. Decide on consistent hover behavior
5. Merge prop APIs to support both implementation patterns

## Additional Notes
- The Next.js implementation includes loading state skeleton not present in Vite version
- Card child components (Header, Title, Description, Content, Footer) should also be compared separately 