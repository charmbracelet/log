# Upgrade Guide: Log v2

This guide covers everything you need to upgrade from Log v1 to v2.

> [!TIP]
> Most upgrades require just two changes: updating your import path and running `go get`. The API remains largely the same.

## Table of Contents

- [Quick Start](#quick-start)
- [Import Path Changes](#import-path-changes)
- [Dependency Updates](#dependency-updates)
- [Breaking Changes](#breaking-changes)
- [Migration Checklist](#migration-checklist)
- [Common Issues](#common-issues)

## Quick Start

For most projects, upgrading is straightforward:

1. Update your import path
2. Update dependencies
3. Fix any type references (if you customize styles)

**Estimated time:** 5-10 minutes for most codebases.

## Import Path Changes

The import path has changed to use the Charm vanity domain.

### Before (v1)

```go
import "github.com/charmbracelet/log"
```

### After (v2)

```go
import "charm.land/log/v2"
```

**Update in your project:**

```bash
# Find all files that need updating
grep -r "github.com/charmbracelet/log" .

# Use your editor's find-and-replace to update imports
# Or use a tool like gofmt with rewrites
```

## Dependency Updates

Update your `go.mod` file:

```bash
go get charm.land/log/v2@latest
go mod tidy
```

### New Dependencies

Log v2 brings these updated dependencies:

- **charm.land/lipgloss/v2** — Lip Gloss v2 for styling
- **github.com/charmbracelet/colorprofile** — Replaces termenv for color profile detection

### Removed Dependencies

- **github.com/muesli/termenv** — No longer needed
- **github.com/aymanbagabas/go-osc52/v2** — Removed (handled by Lip Gloss v2)

## Breaking Changes

### 1. Color Profile Type Change

The `SetColorProfile` method now accepts `colorprofile.Profile` instead of `termenv.Profile`.

#### Before (v1)

```go
import (
    "github.com/charmbracelet/log"
    "github.com/muesli/termenv"
)

logger := log.New(os.Stderr)
logger.SetColorProfile(termenv.TrueColor)
```

#### After (v2)

```go
import (
    "charm.land/log/v2"
    "github.com/charmbracelet/colorprofile"
)

logger := log.New(os.Stderr)
logger.SetColorProfile(colorprofile.TrueColor)
```

**Migration:**

- Replace `termenv.Profile` imports with `colorprofile.Profile`
- Update profile constants:
  - `termenv.TrueColor` → `colorprofile.TrueColor`
  - `termenv.ANSI256` → `colorprofile.ANSI256`
  - `termenv.ANSI` → `colorprofile.ANSI`
  - `termenv.Ascii` → `colorprofile.Ascii`
  - `termenv.NoTTY` → `colorprofile.NoTTY`

### 2. Styles Type Changes

All Lip Gloss style types in the `Styles` struct have changed to use Lip Gloss v2.

#### Before (v1)

```go
import (
    "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/log"
)

styles := log.DefaultStyles()
styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
    Background(lipgloss.Color("204"))
```

#### After (v2)

```go
import (
    "charm.land/lipgloss/v2"
    "charm.land/log/v2"
)

styles := log.DefaultStyles()
styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
    Background(lipgloss.Color("204"))
```

**Changed fields in `Styles` struct:**

- `Caller` — Now `lipgloss/v2.Style`
- `Key` — Now `lipgloss/v2.Style`
- `Keys` — Now `map[string]lipgloss/v2.Style`
- `Levels` — Now `map[Level]lipgloss/v2.Style`
- `Message` — Now `lipgloss/v2.Style`
- `Prefix` — Now `lipgloss/v2.Style`
- `Separator` — Now `lipgloss/v2.Style`
- `Timestamp` — Now `lipgloss/v2.Style`
- `Value` — Now `lipgloss/v2.Style`
- `Values` — Now `map[string]lipgloss/v2.Style`

**Migration:**

If you're using custom styles, update your Lip Gloss import:

```go
// Before
import "github.com/charmbracelet/lipgloss"

// After
import "charm.land/lipgloss/v2"
```

The Lip Gloss v2 API is mostly the same. See the [Lip Gloss v2 upgrade guide][lg-upgrade] for details on any style-specific changes.

[lg-upgrade]: https://charm.land/lipgloss

## Migration Checklist

Use this checklist to ensure a smooth upgrade:

- [ ] **Update import paths** from `github.com/charmbracelet/log` to `charm.land/log/v2`
- [ ] **Update Lip Gloss imports** from `github.com/charmbracelet/lipgloss` to `charm.land/lipgloss/v2` (if using custom styles)
- [ ] **Replace termenv usage** with `colorprofile` (if calling `SetColorProfile`)
- [ ] **Run `go get charm.land/log/v2@latest`**
- [ ] **Run `go mod tidy`** to clean up dependencies
- [ ] **Build your project** with `go build` or `go test`
- [ ] **Run your tests** to verify everything works
- [ ] **Check your logs visually** in different terminals (optional but recommended)

## Common Issues

### Issue: "cannot find package"

**Symptom:**

```
cannot find package "github.com/charmbracelet/log" in any of:
```

**Solution:**

You missed updating an import path. Search your codebase:

```bash
grep -r "github.com/charmbracelet/log" .
```

Update all occurrences to `charm.land/log/v2`.

---

### Issue: "cannot use termenv.Profile as colorprofile.Profile"

**Symptom:**

```
cannot use termenv.TrueColor (type termenv.Profile) as type colorprofile.Profile
```

**Solution:**

Replace `termenv` imports and profile constants with `colorprofile`:

```go
// Before
import "github.com/muesli/termenv"
logger.SetColorProfile(termenv.TrueColor)

// After
import "github.com/charmbracelet/colorprofile"
logger.SetColorProfile(colorprofile.TrueColor)
```

---

### Issue: "lipgloss.Style type mismatch"

**Symptom:**

```
cannot use lipgloss.NewStyle() (type "github.com/charmbracelet/lipgloss".Style)
as type "charm.land/lipgloss/v2".Style
```

**Solution:**

Update Lip Gloss imports to v2:

```go
// Before
import "github.com/charmbracelet/lipgloss"

// After
import "charm.land/lipgloss/v2"
```

---

### Issue: "module declares its path as X but was required as Y"

**Symptom:**

```
module declares its path as: charm.land/log/v2
        but was required as: github.com/charmbracelet/log
```

**Solution:**

Run `go mod tidy` to sync your `go.mod` and `go.sum` files:

```bash
go mod tidy
```

If the issue persists, clear your module cache:

```bash
go clean -modcache
go mod tidy
```

---

### Issue: Colors look wrong after upgrade

**Symptom:**

Logs display incorrectly or with garbled colors in some terminals.

**Solution:**

Log v2 automatically detects color profiles. If you were manually setting a profile, verify you're using the correct `colorprofile` constant:

```go
import "github.com/charmbracelet/colorprofile"

// Explicitly set if needed
logger.SetColorProfile(colorprofile.TrueColor) // or ANSI256, ANSI, etc.
```

In most cases, you can remove manual profile setting—Log v2 handles detection automatically.

---

## Need Help?

If you run into issues not covered here:

1. Check the [Log v2 release notes][release]
2. Visit our [Discord community](https://charm.land/chat)
3. Open an issue on [GitHub](https://github.com/charmbracelet/log/issues)

We're here to help!

[release]: https://github.com/charmbracelet/log/releases

---

Part of [Charm](https://charm.land).

<a href="https://charm.land/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge.jpg" width="400"></a>

Charm热爱开源 • Charm loves open source • نحنُ نحب المصادر المفتوحة
