# Grove Style Guide

Reference for all Grove example templates. Derived from the official logo and branding assets in `branding/`.

## Brand Colors

| Role | Hex | Usage |
|------|-----|-------|
| **Primary (Forest Green)** | `#2E6740` | Links, CTAs, brand highlights, active states |
| **Dark (Espresso)** | `#251917` | Nav bars, footers, headings, dark backgrounds |
| **Light (Cream)** | `#EEEBE3` | Nav text, text-on-dark, card backgrounds, page accents |
| **Page Background** | `#F7F5F0` | Body background (warm neutral, not cold gray) |
| **Body Text** | `#3D2E2A` | Primary readable text (warm dark brown) |
| **Muted Text** | `#7A6B66` | Dates, metadata, secondary info |
| **Border** | `#D9D3CB` | Card borders, dividers |

### Derived Accents

| Role | Hex | Usage |
|------|-----|-------|
| **Green Light** | `#E8F0EA` | Tag backgrounds, subtle highlights |
| **Green Hover** | `#245533` | Hover state for primary green |
| **Cream Dark** | `#DDD8CE` | Hover state for cream buttons/elements |

### Alert Colors

| Type | Background | Text | Border |
|------|-----------|------|--------|
| Info | `#E8F0EA` | `#2E6740` | `#2E6740` |
| Warning | `#FFF3CD` | `#6B5210` | `#E6C547` |
| Error | `#F8D7DA` | `#6B1D24` | `#E8A0A7` |
| Success | `#D4EDDA` | `#1B5E28` | `#A3D4AE` |

### Tag Colors

| Variant | Background | Text |
|---------|-----------|------|
| Green (default) | `#E8F0EA` | `#2E6740` |
| Brown | `#F0EAE4` | `#5C3D2E` |
| Red | `#FEE2E2` | `#991B1B` |
| Purple | `#EDE9FE` | `#5B21B6` |
| Orange | `#FFEDD5` | `#9A3412` |
| Gray | `#EDEBE7` | `#4A3F3A` |

## Typography

```css
font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
```

- Body text: `1rem`, `line-height: 1.6`
- Article body: `line-height: 1.7`
- Small/meta text: `0.85rem`
- Headings: bold, color `#251917`

## Buttons

| Variant | Background | Text | Border |
|---------|-----------|------|--------|
| Primary | `#2E6740` | `#EEEBE3` | `#2E6740` |
| Secondary | `#251917` | `#EEEBE3` | `#251917` |
| Outline | transparent | `#2E6740` | `#2E6740` |
| Default | `#7A6B66` | `#EEEBE3` | `#7A6B66` |

Shared: `padding: 0.5rem 1.25rem`, `border-radius: 6px`, `font-weight: 600`, `font-size: 0.9rem`, `border: 2px solid`.

Hover: darken background by one step (e.g., Primary hover `#245533`).

## Components

### Navigation
- Background: `#251917` (Espresso)
- Brand name: `#2E6740` (Forest Green), bold, `1.4rem`
- Nav links: `#EEEBE3` (Cream)
- Layout: flexbox, `space-between`

### Cards
- Background: `#EEEBE3` or `#fff`
- Border: `1px solid #D9D3CB`
- Border-radius: `8px`
- Padding: `1.5rem`
- Hover: `box-shadow: 0 2px 8px rgba(37, 25, 23, 0.1)` (warm shadow)

### Footer
- Background: `#251917`
- Text: `#B5ADA8`
- Accent links: `#2E6740`
- Padding: `2rem`, centered

### Tags
- Pill shape: `border-radius: 999px`
- Padding: `0.2rem 0.6rem`
- Font: `0.75rem`, weight `600`

### Alerts
- `border-left: 4px solid` (type color)
- Border-radius: `6px`
- Padding: `1rem 1.25rem`
- Icons: Unicode (info `i`, warning `!`, error `x`, success check)

## Layout

| Context | Max-width |
|---------|-----------|
| Blog content | `960px` |
| Store/product grid | `1080px` |
| Email | `600px` |

Grid gaps: `1.5rem`. Section padding: `2rem`.

## Guiding Principles

1. **Warm, not cold** -- use cream/brown tones instead of pure white/gray
2. **Green is the accent** -- forest green replaces any red/blue accent usage
3. **Match the logo** -- the three logo colors (`#2E6740`, `#251917`, `#EEEBE3`) are the foundation
4. **Generous whitespace** -- let content breathe with consistent spacing
5. **Subtle transitions** -- `transition: 0.2s` on interactive elements
