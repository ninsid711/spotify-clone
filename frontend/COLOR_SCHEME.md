# 🎨 OneDark Color Scheme

The application now uses the **OneDark** color palette, inspired by the popular Atom One Dark theme.

## Color Palette

### Primary Colors
| Color | Hex Code | Usage |
|-------|----------|-------|
| **Primary Blue** | `#61afef` | Main accent, buttons, links |
| **Secondary Green** | `#98c379` | Success states, accents |
| **Accent Purple** | `#c678dd` | Genre badges, highlights |
| **Warning Yellow** | `#e5c07b` | Warnings, notifications |
| **Error Red** | `#e06c75` | Error messages, danger |

### Background Colors
| Color | Hex Code | Usage |
|-------|----------|-------|
| **Background** | `#282c34` | Main app background |
| **Surface** | `#21252b` | Cards, navbar, surfaces |
| **Elevated Surface** | `#2c313a` | Input fields, hover states |
| **Border** | `#181a1f` | Borders, dividers |

### Text Colors
| Color | Hex Code | Usage |
|-------|----------|-------|
| **Text Bright** | `#ffffff` | Headings, important text |
| **Text Primary** | `#abb2bf` | Body text, labels |
| **Text Secondary** | `#5c6370` | Subtext, placeholders |

## Visual Examples

### Buttons
- **Primary Button**: Blue (`#61afef`) background with dark surface text
- **Secondary Button**: Transparent with blue border on hover
- **Danger Button**: Red (`#e06c75`) background

### Cards
- **Background**: Elevated surface (`#2c313a`)
- **Border**: Dark border (`#181a1f`)
- **Hover**: Border changes to primary blue

### Genre Badges
- **Background**: Purple (`#c678dd`)
- **Text**: Dark surface color
- **Hover**: Slightly lighter with blue border

### Play Button
- **Background**: Primary blue (`#61afef`)
- **Shadow**: Blue glow effect
- **Hover**: Slightly lighter blue with enhanced glow

## CSS Variables

All colors are defined as CSS variables in `App.css`:

```css
:root {
  /* OneDark Color Scheme */
  --primary-color: #61afef;
  --secondary-color: #98c379;
  --accent-color: #c678dd;
  --warning-color: #e5c07b;
  --error-color: #e06c75;
  --background-color: #282c34;
  --surface-color: #21252b;
  --elevated-surface: #2c313a;
  --text-primary: #abb2bf;
  --text-secondary: #5c6370;
  --text-bright: #ffffff;
  --hover-color: #2c313a;
  --border-color: #181a1f;
}
```

## Component Updates

### Updated Components
- ✅ Navbar - Blue branding, updated button colors
- ✅ TrackCard - Blue play button with glow, purple genre badges
- ✅ Auth Pages - Updated cards, inputs, and buttons
- ✅ Home Page - Blue search button, updated genre pills
- ✅ Playlists - Updated card backgrounds and borders
- ✅ PlaylistDetail - Blue gradient header

### Key Changes
1. **Primary color changed** from Spotify green (`#1db954`) to OneDark blue (`#61afef`)
2. **Background** changed from pure black to OneDark gray (`#282c34`)
3. **Surface colors** updated to match OneDark theme
4. **Genre badges** now use purple accent (`#c678dd`)
5. **Hover effects** include blue borders and glows
6. **Error states** use OneDark red (`#e06c75`)

## Theme Consistency

All components now follow the OneDark color scheme:
- Consistent blue accent throughout
- Proper contrast ratios for accessibility
- Subtle gradients and glows for depth
- Smooth transitions between states

## Developer Notes

To customize colors further:
1. Edit CSS variables in `App.css`
2. All components will automatically update
3. Maintain contrast ratios for accessibility
4. Test hover states for visibility

---

**Theme applied:** OneDark  
**Updated:** November 6, 2025  
**Status:** ✅ Complete

