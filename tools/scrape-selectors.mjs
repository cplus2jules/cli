import { readFileSync, existsSync } from 'fs';
import { resolve } from 'path';

if (process.argv.length < 3) {
  console.error('Usage: node scrape-selectors.mjs <path-to-xpui.js>');
  process.exit(1);
}

let targetPath = resolve(process.argv[2]);

if (!existsSync(targetPath)) {
  console.error(`Error: File not found at ${targetPath}`);
  process.exit(1);
}

// If user points to a folder, try to find a valid bundle inside
if (existsSync(targetPath) && !targetPath.endsWith('.js')) {
    const candidate1 = resolve(targetPath, 'xpui.js');
    const candidate2 = resolve(targetPath, 'xpui-snapshot.js');
    if (existsSync(candidate1)) targetPath = candidate1;
    else if (existsSync(candidate2)) targetPath = candidate2;
    else {
        console.error(`Error: Could not find xpui.js or xpui-snapshot.js in ${targetPath}`);
        process.exit(1);
    }
}

const bundle = readFileSync(targetPath, 'utf8');
console.log(`Analyzing bundle: ${targetPath}...`);

// Known stable structural strings we can regex match
const patterns = {
  'main-view':       /Root__main-view|main-view-container/,
  'now-playing-bar': /Root__now-playing-bar|now-playing-widget|now-playing-bar/,
  'sidebar':         /Desktop_LeftSidebar_Id|LeftSidebar/,
  'topbar':          /main-topBar-container|Root__top-bar|topbar/,
  'right-sidebar':   /Root__right-sidebar/,
  'lyrics-cinema':   /lyrics-cinema/,
  'buddy-feed':      /BuddyFeed|getBuddyFeedAPI/
};

const found = {};

for (const [name, pattern] of Object.entries(patterns)) {
  const match = bundle.match(new RegExp(`([\\w-]*${pattern.source}[\\w-]*)`, 'g'));
  if (match) {
    // 1. Prioritize matches starting with "Root__"
    // 2. Otherwise, pick the SHORTEST match (usually the top-level container)
    const sorted = match.sort((a, b) => {
        const aRoot = a.startsWith('Root__');
        const bRoot = b.startsWith('Root__');
        if (aRoot && !bRoot) return -1;
        if (!aRoot && bRoot) return 1;
        return a.length - b.length;
    });
    found[name] = '.' + sorted[0];
  } else {
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
