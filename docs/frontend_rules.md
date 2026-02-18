# Frontend Development Rules & UI Standard

**Version:** 2.1
**Enforcement:** STRICT - DEVIATIONS ARE NOT PERMITTED
**Theme:** "Industrial Density" (Compact, Data-Heavy, Serious)

---

## 1. Technology Stack
- **Framework:** Vue 3 (Composition API, `<script setup>`)
- **Styling:** Tailwind CSS v4.0
- **Components:** shadcn-vue (Radix Vue based)
- **State/Data:** TanStack Query v5
- **Icons:** Lucide Vue (Stroke width: 1.5px ALWAYS)

---

## 2. STRICT CSS & Styling Rules

### 2.1 No Hardcoded CSS Classes in Components
**CRITICAL RULE:** All CSS classes must be defined centrally in `src/assets/main.css`. Do NOT hardcode CSS classes directly in Vue components.

- ❌ **FORBIDDEN:** Hardcoding Tailwind utility classes directly in component templates (e.g., `class="bg-slate-700 text-white p-4"`).
- ❌ **FORBIDDEN:** Arbitrary values (e.g., `w-[350px]`, `bg-[#1e293b]`, `h-[37px]`).
- ❌ **FORBIDDEN:** Inline `style="..."` attributes.
- ❌ **FORBIDDEN:** `<style>` blocks in components (unless absolutely necessary for complex animations not possible with Tailwind).
- ✅ **REQUIRED:** Define all CSS classes in `main.css` using Tailwind's `@apply` directive or custom CSS.
- ✅ **REQUIRED:** Reference classes by semantic names in components (e.g., `class="card-header"` instead of `class="bg-slate-100 text-slate-700"`).

**Example (main.css):**
```css
.card-header {
  @apply bg-slate-100 text-slate-700 uppercase text-xs font-semibold;
}

.btn-primary {
  @apply bg-slate-700 text-white hover:bg-slate-800 rounded-md h-9 px-4;
}
```

**Example (Component):**
```vue
<!-- ✅ CORRECT -->
<div class="card-header">Header</div>
<button class="btn-primary">Submit</button>

<!-- ❌ WRONG -->
<div class="bg-slate-100 text-slate-700 uppercase text-xs font-semibold">Header</div>
<button class="bg-slate-700 text-white hover:bg-slate-800 rounded-md h-9 px-4">Submit</button>
```

### 2.2 Component Variants
Use `class-variance-authority` (cva) for component variants to map props to Tailwind classes.
```typescript
const buttonVariants = cva(
  "inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50",
  {
    variants: {
      variant: {
        default: "bg-primary text-primary-foreground hover:bg-primary/90",
        destructive: "bg-destructive text-destructive-foreground hover:bg-destructive/90",
        outline: "border border-input bg-background hover:bg-accent hover:text-accent-foreground",
        secondary: "bg-secondary text-secondary-foreground hover:bg-secondary/80",
        ghost: "hover:bg-accent hover:text-accent-foreground",
        link: "text-primary underline-offset-4 hover:underline",
      },
      size: {
        default: "h-9 px-4 py-2",
        sm: "h-8 rounded-md px-3",
        lg: "h-10 rounded-md px-8",
        icon: "h-9 w-9",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  }
)
```

---

## 3. Component Usage Rules

### 3.1 No Native HTML Components - Only shadCN
**CRITICAL RULE:** Use ONLY shadcn-vue components. Native HTML elements are **FORBIDDEN** for interactive/form elements.

| Native HTML (❌ FORBIDDEN) | shadCN Equivalent (✅ REQUIRED) |
|---------------------------|--------------------------------|
| `<button>` | `<Button>` from `@/components/ui/button` |
| `<input>` | `<Input>` from `@/components/ui/input` |
| `<select>` / `<option>` | `<Select>`, `<SelectContent>`, `<SelectItem>` |
| `<textarea>` | `<Textarea>` from `@/components/ui/textarea` |
| `<checkbox>` | `<Checkbox>` from `@/components/ui/checkbox` |
| `<radio>` | `<RadioGroup>`, `<RadioGroupItem>` |
| `<dialog>` | `<Dialog>`, `<DialogContent>` |
| `<form>` | `<Form>` from `@/components/ui/form` (with vee-validate) |
| `<label>` | `<Label>` from `@/components/ui/label` |
| `<table>` / `<tr>` / `<td>` | `<Table>`, `<TableRow>`, `<TableCell>` |
| `<a>` (navigation) | `<RouterLink>` or shadCN navigation components |

**Exceptions:**
- Semantic HTML for content structure is allowed: `<div>`, `<span>`, `<p>`, `<h1>`-`<h6>`, `<ul>`, `<ol>`, `<li>`, `<section>`, `<article>`, `<header>`, `<footer>`, `<nav>`, `<main>`, `<aside>`.

**Example:**
```vue
<!-- ❌ WRONG - Native HTML -->
<button onclick="handleClick">Submit</button>
<input type="text" v-model="value" />
<select v-model="selected"><option value="a">A</option></select>

<!-- ✅ CORRECT - shadCN Components -->
<Button @click="handleClick">Submit</Button>
<Input v-model="value" type="text" />
<Select v-model="selected">
  <SelectContent>
    <SelectItem value="a">A</SelectItem>
  </SelectContent>
</Select>
```

---

## 4. Visual Language & Theming

### 4.1 Color Palette (Slate Monochromatic)
- **Primary Action:** `bg-slate-700` (Variables: `--primary`). **NEVER** use `slate-900` for buttons.
- **Backgrounds:**
    - App: `bg-slate-50/50` (Variables: `--background`)
    - Cards/Surface: `bg-white` (Variables: `--card`)
- **Text:**
    - Headings: `text-slate-900`
    - Body: `text-slate-700`
    - Muted: `text-slate-500`
- **Borders:** `border-slate-300` (Variables: `--border`) for ALL inputs and cards.

### 4.2 Typography & Radius
- **Font:** Inter (or system sans-serif).
- **Radius:** `rounded-md` (Standard).
    - ❌ **PROHIBITED:** `rounded-xl`, `rounded-2xl`, `rounded-full` (except avatars).
- **Density:** Compact padding. Reduce default shadcn spacing by ~30%.

### 4.3 Layout Structure
- **Sidebar:** Collapsible left sidebar (`w-64` -> `w-16`). Icons must constitute the *only* content in collapsed mode and be centered (`px-2`).
- **Page Container:** `max-w-7xl mx-auto px-4 md:px-6 py-6`.
- **Spacing:** `space-y-4` for vertical rhythm.

---

## 5. Component Patterns

### 5.1 Data Tables
- **Header:** `bg-slate-100`, `text-slate-700`, uppercase `text-xs font-semibold`.
- **Rows:** `bg-white`, `hover:bg-slate-50`.
- **Actions:** Dropdown menu at row end.
- **Loading:** Skeleton rows or centered spinner.

### 5.2 Forms
- **Input Height:** `h-9` (Compact).
- **Labels:** `text-sm font-medium text-slate-700`.
- **Required:** Red asterisk `text-rose-600`.
- **Validation:** Inline `text-rose-600 text-xs`.

### 5.3 Modals vs. Drawers
- **Complex Forms (Create/Edit):** Right-side Sheet (`Sheet` component).
    - Width: `w-full sm:max-w-lg` (default).
    - Footer: Fixed at bottom with "Cancel" and "Save" buttons.
- **Simple Confirmations:** Centered Dialog.

### 5.4 Sheet Component Standard

All sheets in the application must follow this consistent structure:

#### Sheet Widths
| Class | Width | Usage |
|-------|-------|-------|
| `sheet-content` | `sm:max-w-lg` (512px) | Simple forms, confirmations |
| `sheet-content-wide` | `sm:max-w-4xl` (896px) | Complex data, tables, multi-tab content |

#### Standard Sheet Structure
```vue
<Sheet v-model:open="open">
  <SheetContent class="sheet-content-wide">
    <!-- Header -->
    <SheetHeader class="sheet-header">
      <SheetTitle class="sheet-title">Title</SheetTitle>
      <SheetDescription class="sheet-description">Description text</SheetDescription>
    </SheetHeader>

    <!-- Body -->
    <div class="sheet-body">
      <!-- With tabs -->
      <Tabs defaultValue="tab1" class="sheet-tabs">
        <TabsList class="sheet-tabs-header">
          <TabsTrigger value="tab1">Tab 1</TabsTrigger>
          <TabsTrigger value="tab2">Tab 2</TabsTrigger>
        </TabsList>
        
        <TabsContent value="tab1" class="sheet-tab-content">
          <div class="tab-panel">
            <!-- Content -->
          </div>
        </TabsContent>
      </Tabs>
      
      <!-- Or without tabs -->
      <div class="sheet-body-section">
        <!-- Content -->
      </div>
    </div>

    <!-- Footer -->
    <SheetFooter class="sheet-footer">
      <div class="sheet-footer-actions">
        <Button variant="outline" class="sheet-footer-btn">Cancel</Button>
        <Button class="sheet-footer-btn">Save</Button>
      </div>
    </SheetFooter>
  </SheetContent>
</Sheet>
```

#### Sheet CSS Classes Reference

| Class | Purpose |
|-------|---------|
| `sheet-content` | Standard width sheet (sm:max-w-lg) |
| `sheet-content-wide` | Wide sheet for tables/data (sm:max-w-4xl) |
| `sheet-header` | Header section with border |
| `sheet-title` | Title styling |
| `sheet-description` | Description text styling |
| `sheet-body` | Main content area with padding and scroll |
| `sheet-body-section` | Spaced content section |
| `sheet-tabs` | Tabs container |
| `sheet-tabs-header` | Sticky tab header |
| `sheet-tab-content` | Tab content padding |
| `tab-panel` | Bordered content panel inside tabs |
| `sheet-footer` | Footer with border |
| `sheet-footer-actions` | Footer button container |
| `sheet-footer-btn` | Stretch button in footer |

#### Identity Header (for entity sheets)
Use `sheet-header-identity` for sheets showing entity details:

```vue
<div class="sheet-header-identity">
  <div class="sheet-identity-row">
    <div class="sheet-identity-icon">
      <Package class="h-4 w-4" />
    </div>
    <div class="sheet-identity-info">
      <p class="sheet-identity-title">{{ entity.name }}</p>
      <p class="sheet-identity-desc">{{ entity.description }}</p>
    </div>
    <Badge class="sheet-identity-badge">{{ entity.id }}</Badge>
  </div>
</div>
```

### 5.5 Badge System
Use subtle, desaturated backgrounds (Sonner-style).
- **Success:** Emerald (Subtle)
- **Warning:** Amber (Subtle)
- **Danger:** Rose (Subtle)
- **Info:** Blue/Slate (Subtle)
- **Neutral:** Slate/Gray (Subtle)

---

## 6. Directory Structure
```
src/
├── assets/
│   └── main.css             # Central Tailwind v4 imports & @theme config (ALL CSS GOES HERE)
├── components/
│   ├── ui/                  # Shadcn primitives (tailored to theme)
│   ├── common/              # App-wide shared (AppSidebar, PageHeader)
│   └── domain/              # Feature-specific components
├── composables/             # Business logic & Query hooks
├── layouts/                 # DashboardLayout.vue
└── views/                   # Route targets
```

---

## 7. Common Mistakes to Avoid
1.  **Using native HTML elements:** Use shadCN components, NOT `<button>`, `<input>`, `<select>`, etc.
2.  **Hardcoding CSS in components:** Define all classes in `main.css`, use semantic class names.
3.  **Using `slate-900` for buttons:** It is too dark. Use `slate-700` / `primary` variant.
4.  **Hardcoding colors:** Use semantic names (e.g., `text-muted-foreground`) or standard Tailwind colors (e.g., `text-slate-500`), NOT hex codes.
5.  **Inconsistent borders:** Always use `border-slate-300` / `border-border` for inputs.
6.  **Arbitrary sizing:** Do not use `w-[123px]`. Use the nearest grid step (e.g., `w-32`).

---
**Enforcement:** Code reviews will reject any PR containing:
- Hardcoded CSS classes in components (must be in `main.css`)
- Native HTML form/interactive elements (must use shadCN components)
- Hardcoded arbitrary CSS values
- Any deviations from this standard
