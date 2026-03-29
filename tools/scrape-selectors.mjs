import { readFileSync, existsSync } from 'fs';
import { resolve } from 'path';

if (process.argv.length < 3) {
  console.error('Usage: node scrape-selectors.mjs <path-to-xpui.js>');
  process.exit(1);
}

const targetPath = resolve(process.argv[2]);
if (!existsSync(targetPath)) {
  console.error(`Error: File not found at ${targetPath}`);
  process.exit(1);
}

const bundle = readFileSync(targetPath, 'utf8');

// Known stable structural strings we can regex match
const patterns = {
  'main-view': /Root__main-view/,
  'now-playing-bar': /now-playing-widget/,
  'sidebar': /LeftSidebar/,
  'topbar': /main-topBar/,
  'right-sidebar': /Root__right-sidebar/,
  'lyrics-cinema': /lyrics-cinema/,
  'buddy-feed': /BuddyFeed/
};

const found = {};

for (const [name, pattern] of Object.entries(patterns)) {
  // We match standard CSS classes attached to React components logic like:
  // "className:"Root__main-view"" or ".Root__main-view {" or similar react-hash artifacts
  const match = bundle.match(new RegExp(`([\\w-]+${pattern.source}[\\w-]*)`, 'g'));
  if (match) {
    found[name] = '.' + match[0];
  } else {
    // If we can't find a selector, explicitly set it to null so maintainers see what's missing
    found[name] = null;
  }
}

// Map encore variable bindings
const cssAliases = {
  "--spice-main": "--encore-base-color-black",
  "--spice-sidebar": "--encore-tinted-base-color",
  "--spice-player": "--encore-base-color-black",
  "--spice-text": "--encore-text-color-on-dark-base",
  "--spice-subtext": "--encore-text-color-on-dark-subdued",
  "--spice-button": "--encore-essential-positive-set-color",
  "--spice-button-active": "--encore-essential-positive-set-highlight-color"
};

const output = {
  "version": "1.X.XX",
  "spotifyVersionRange": "1.X.XX.0 - 1.X.XX.9999",
  "selectors": found,
  "cssVarAliases": cssAliases
}

console.log(JSON.stringify(output, null, 2));
