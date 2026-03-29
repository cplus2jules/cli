# Spicetify Theme Selectors Guide

Because Spicetify modifies an active React app, classnames (like `.main-view-container__a1b2c`) frequently change whenever Spotify pushes an update, silently breaking thousands of active custom themes.

Going forward, theme authors should exclusively use the newly injected **Stable Selectors Pipeline**.

## Data Attributes

Spicetify natively tracks Spotify's DOM and attaches structural identifiers (`data-spicetify="xxx"`) to components dynamically securely avoiding cross-version breakages!

### Supported Views:
Instead of styling `.Root__main-view` or whatever its future name becomes, leverage our stable semantic wrapper:
```css
/* Old way (Fragile) */
.main-view-container__a1b2c {
  background: black !important;
}

/* New way (Stable) */
[data-spicetify="main-view"] {
  background: black !important;
}
```

Current guaranteed wrappers:
- `[data-spicetify="main-view"]`
- `[data-spicetify="now-playing-bar"]`
- `[data-spicetify="sidebar"]`
- `[data-spicetify="topbar"]`
- `[data-spicetify="right-sidebar"]`

## CSS Tokens Strategy

We now auto-alias Spotify's internal `--encore-*` design tokens back into `--spice-*` variables. Changing values inside your structural `color.ini` reliably cascades through all of Spotify components natively again. To override the mapped blocks, simply do so within your `user.css` file as the script prepends aliases prior to resolving your actual file contents.
