# Frontend Implementation Plan (Demo, Industrial Density)

Last updated: 2026-02-18

Conventions:
- Dates: UK format `yyyy-MM-dd` and `yyyy-MM-dd HH:mm` (24h).
- Styling: Follow `docs/frontend_rules.md` (Tailwind utilities only, no arbitrary values, primary `bg-slate-700`, borders `border-slate-300`, `rounded-md` radius).
- Deletions: Always prompt with shadcn dialog before destructive actions.
- Global error handling: centralized toast + logging via API interceptors and query error handlers.

## Task Checklist

- [x] Foundation layout
  - Replace `src/App.vue` shell with router outlet, page container (`max-w-7xl mx-auto px-4 md:px-6 py-6`, `space-y-4`), no arbitrary values.
  - Ensure `src/assets/main.css` stays aligned to palette/radius; remove/avoid inline styles or hex overrides.
  - Add TanStack Query provider in `src/main.ts` (no inline styles).
- [x] UI kit & primitives
  - Install/config shadcn-vue primitives under `src/components/ui`.
  - Define button/input/badge variants via `class-variance-authority` (sizes `h-9`, `rounded-md`, primary `slate-700`).
  - Create reusable `ConfirmDialog` for delete flows.
- [x] Global utilities
  - Add UK date formatter helper `src/lib/date.ts` returning `yyyy-MM-dd` and `yyyy-MM-dd HH:mm`.
  - Add API client `src/lib/api.ts` with interceptors → global error handler/toasts; retries where safe.
  - Add `useGlobalError` composable to surface errors uniformly (hooks + toasts).
- [ ] Routing & auth
  - Configure routes: `login`, `dashboard`, `employees`, `employees/:id`, `profiles`, `events`, `certifications`, 404. (partially stubbed)
  - Implement Pinia auth store (role + employeeId); guard routes for admin vs employee self-scope. (auth store added; guards stubbed)
- [ ] Layout components
  - `AppSidebar` collapsible (`w-64`/`w-16`, icon-only collapsed).
  - `PageHeader` with title/actions; `KpiTile`, `StatusBadge`, `DataTable` (sortable, sticky header, skeleton rows), `Drawer/Sheet` wrappers, `GridLegend`.
- [ ] Dashboard view
  - KPI tiles (overall compliance, expiring, expired, total employees) using palette.
  - Coverage grid component (course/profile mode) with filters (site/department/profile/course) and drill-down drawer (employee list, export CSV client-side).
  - All dates displayed in UK format.
- [ ] Employees
  - List with search/filter (department/site/status); row action “View”.
  - Detail with demographics, assigned profiles, requirements matrix (status badges), transcript table; UK dates.
- [ ] Profiles
  - List + detail (requirements table, members list).
  - Sheets for manual add/exclude; destructive actions gated by `ConfirmDialog`.
- [ ] Events
  - Table + sheet to create event (courseId dropdown from courses API, site/date/instructor, roster multiselect).
  - Attendance confirm action; any delete/rollback uses confirmation dialog.
- [ ] Certifications
  - Form to add cert record (courseId dropdown for certification selection, issuer/issue/expiry/notes) with UK date pickers; list with filters.
- [ ] Courses Management
  - CoursesView with tabs for courses/certifications.
  - CourseFormSheet for create/edit with validation.
  - Delete confirmation dialog.
  - useCourses composable with TanStack Query.
- [ ] Demo/fixtures
  - Add fixtures matching `docs/demo_pack.md` in `src/data`; `VITE_USE_FIXTURES` flag to swap data source.
- [ ] Responsiveness & A11y
  - Sidebar collapse <1024px; sheets full-screen on mobile; tables horizontal scroll with locked first column.
  - Focus rings, aria labels on icon buttons; keyboard navigation for grid.
- [ ] Testing
  - Unit tests: date formatter (UK), status→badge mapper, API error interceptor.
  - Manual demo script aligned to Scenes A–D; verify delete confirmations and error toasts.

## Open Questions
- Awaiting backend endpoints for dashboards, events, certifications, auth. Until then, use fixtures/stubs.
- Confirm export format requirement (CSV vs PDF) for drill-down; default to CSV.
