# Issue: Card Component Visual Inconsistencies

## Issue Details

- **Component**: Card
- **Priority**: High
- **Type**: Visual Inconsistency
- **Status**: Open
- **Identified**: April 16, 2023
- **Assigned**: Unassigned

## Description

During the UI consistency audit, significant visual inconsistencies were identified between the Vite and Next.js implementations of the Card component. These inconsistencies affect the overall design coherence of the application and should be addressed promptly.

## Specific Issues

1. **Shadow Implementation**:
   - Vite: Uses `box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1)`
   - Next.js: Uses `filter: drop-shadow(0 4px 12px rgba(0, 0, 0, 0.12))`
   - Impact: Different visual appearance, particularly at component edges

2. **Border Radius**:
   - Vite: `border-radius: 8px`
   - Next.js: `border-radius: 12px`
   - Impact: Noticeable difference in component shape, inconsistent with design system

3. **Card Footer Padding**:
   - Vite: `padding: 12px`
   - Next.js: `padding: 16px`
   - Impact: Different spacing between card content and footer elements

4. **Hover Effect**:
   - Vite: No hover effect
   - Next.js: Subtle elevation increase on hover
   - Impact: Inconsistent interaction patterns between implementations

## Additional Notes

The Next.js implementation includes a loading state skeleton not present in the Vite version. While this is an enhancement rather than an inconsistency, it should be documented as an intentional improvement.

## Recommended Solution

1. Standardize on the box-shadow approach for consistent shadow rendering
2. Align border radius to 8px as per the design system specification
3. Normalize padding values to 12px for consistency
4. Decide on whether to add hover effect to Vite or remove from Next.js based on design requirements

## Screenshots

[Placeholder for screenshots showing the differences]

## Implementation Estimate

- Estimated time: 3 hours
- Complexity: Medium
- Files to modify:
  - `frontend_next/src/components/ui/Card.tsx`
  - Associated styles/CSS files 