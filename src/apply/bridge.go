package apply

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spicetify/cli/src/utils"
)

const bridgeMarker = "/* spicetify-bridge-injected */"

func InsertBridgeScript(content, script string) string {
	if strings.Contains(content, bridgeMarker) {
		return content
	}
	return bridgeMarker + "\n" + script + "\n" + content
}

func GenerateBridgeScript(sm *utils.SelectorMap) string {
	if sm == nil {
		return ""
	}

	selectorsJSON, _ := json.Marshal(sm.Selectors)

	return fmt.Sprintf(`
(function SpicetifyBridge() {
    'use strict';

    const SELECTOR_MAP = %s;
    const ATTR = 'data-spicetify';
    let stamped = new WeakSet();

    function stamp() {
        for (const [name, selector] of Object.entries(SELECTOR_MAP)) {
            if (!selector) continue; // Skip null mappings from scrape tool
            try {
                const elements = document.querySelectorAll(selector);
                elements.forEach(el => {
                    if (!el.hasAttribute(ATTR)) {
                        el.setAttribute(ATTR, name);
                        stamped.add(el);
                    }
                });
            } catch(e) {
                // Selector invalid for this Spotify version
            }
        }
    }

    function stampStableAnchors() {
        const stableMap = {
            'main-view':       '[data-testid="main-view"]',
            'now-playing-bar': '[data-testid="now-playing-widget"]',
            'topbar':          '[data-testid="topbar"]',
            'sidebar':         'nav[aria-label="Main"]',
        };
        for (const [name, selector] of Object.entries(stableMap)) {
            try {
                document.querySelectorAll(selector).forEach(el => {
                    if (!el.hasAttribute(ATTR)) {
                        el.setAttribute(ATTR, name);
                    }
                });
            } catch(e) {}
        }
    }

    const observer = new MutationObserver(() => {
        stamp();
        stampStableAnchors();
    });

    function startBridge() {
        if (!document.body) {
            document.addEventListener('DOMContentLoaded', startBridge, { once: true });
            return;
        }
        stamp();
        stampStableAnchors();
        observer.observe(document.body, {
            childList: true,
            subtree: true
        });
    }

    startBridge();
})();
`, string(selectorsJSON))
}

func GenerateCSSVariableAliasBlock(sm *utils.SelectorMap) string {
	if sm == nil || len(sm.CSSVarAliases) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(":root {\n")
	for spiceVar, encoreVar := range sm.CSSVarAliases {
		sb.WriteString(fmt.Sprintf("  %s: var(%s);\n", spiceVar, encoreVar))
	}
	sb.WriteString("}\n")
	return sb.String()
}
