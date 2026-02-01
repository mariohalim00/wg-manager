
You are working on the frontend of an existing WireGuard management web application.

Your task is to refine, redesign, and polish the user interface to match a modern, production-grade infrastructure dashboard experience. based on the design : ![alt text](../001-frontend-implementation/design/stitch_vpn_management_dashboard/vpn_management_dashboard_1/screen.png). Also 

## Important constraints:
- Do NOT modify any backend logic, APIs, data models, or state management.
- Do NOT change existing business logic.
- Only improve layout, visual hierarchy, spacing, color usage, responsiveness, and UI consistency.
- Assume all data, actions, and behaviors already work correctly.

## UI goals:
- Provide a clear app-shell layout with a persistent sidebar, contextual header, and card-based content.
- Emphasize clarity, scannability, and operational confidence.
- Use only two main color palettes (primary and secondary), supported by neutral grays and limited accents.
- Support both dark mode (default) and light mode with consistent semantics.
- Ensure responsive behavior across mobile, tablet, and desktop breakpoints.
- Maintain a restrained, professional aesthetic suitable for DevOps and infrastructure tooling.


## design specs (JSONc format)
- the values in the jsonC below are hardcoded, however, some are meant to be dynamic like peersCount etc. Just keep that in mind
```jsonc
{
  "product": {
    "name": "WireGuard Manager Web UI",
    "purpose": "Visual management dashboard for a WireGuard interface, peers, and traffic stats",
    "designTone": {
      "darkMode": "Primary mode; modern, muted, professional, infrastructure-focused",
      "lightMode": "Secondary mode; clean, low-contrast, enterprise-neutral"
    }
  },

  "layout": {
    "globalStructure": {
      "type": "App shell",
      "regions": [
        "Left sidebar (persistent navigation)",
        "Top header (context + controls)",
        "Main content area (cards + tables)"
      ]
    },

    "sidebar": {
      "width": {
        "default": "w-64",
        "collapsed": "w-16"
      },
      "contents": [
        "App logo + title",
        "Primary navigation links",
        "Usage indicator",
        "Primary CTA (Add New Peer)"
      ],
      "navigationItems": [
        "Dashboard (active)",
        "Peers",
        "Settings",
        "Logs"
      ],
      "behavior": {
        "desktop": "Collapsible, the behaviour is the same like right now",
        "tablet": "Collapsible, the behaviour is the same like right now",
        "mobile": "Off-canvas drawer"
      }
    },

    "header": {
      "left": [
        "Page title: Interface name (wg0)",
        "Sub-label: Interface context"
      ],
      "right": [
        "Search peers input",
        "Notification icon",
        "User avatar / profile menu"
      ],
      "height": "h-16",
      "position": "sticky"
    }
  },

  "contentSections": {
    "statusCards": {
      "grid": {
        "desktop": "4 columns",
        "tablet": "2 columns",
        "mobile": "1 column"
      },
      "cards": [
        {
          "label": "Status",
          "value": "Active",
          "indicator": "Green dot",
          "emphasis": "Binary state"
        },
        {
          "label": "Public Key",
          "value": "Truncated key with copy affordance",
          "interaction": "Click to copy"
        },
        {
          "label": "Listening Port",
          "value": "51820"
        },
        {
          "label": "Subnet",
          "value": "10.0.0.0/24"
        }
      ]
    },

    "trafficCharts": {
      "layout": "Two equal cards side-by-side",
      "cards": [
        {
          "title": "Total Received",
          "value": "14.2 GB",
          "delta": "+12%",
          "chart": {
            "type": "Smooth line chart",
            "style": "Minimal grid, no axes labels",
            "accent": "Primary color"
          }
        },
        {
          "title": "Total Sent",
          "value": "8.4 GB",
          "delta": "+5%",
          "chart": {
            "type": "Smooth line chart",
            "style": "Minimal grid",
            "accent": "Secondary color"
          }
        }
      ]
    },

    "peersTable": {
      "header": {
        "title": "Active Peers",
        "count": 12,
        "actions": [
          "Filter",
          "Export"
        ]
      },
      "columns": [
        "Status",
        "Peer Name",
        "Internal IP",
        "Transfer (Up / Down)",
        "Last Handshake",
        "Actions"
      ],
      "rowStates": {
        "online": "Green dot + brighter text",
        "offline": "Muted text + gray dot"
      },
      "rowContent": {
        "peerName": {
          "icon": "Device-type icon",
          "subtitle": "OS / device descriptor"
        },
        "transfer": {
          "format": "Up arrow + Down arrow with values",
          "colors": {
            "upload": "Green",
            "download": "Blue"
          }
        }
      },
      "footer": "View All Peers link"
    }
  },

  "colorSystem": {
    "palettes": {
      "primary": {
        "intent": "Status, highlights, charts",
        "approx": "Tailwind blue / indigo range"
      },
      "secondary": {
        "intent": "Success, connectivity, active states",
        "approx": "Tailwind emerald / green range"
      }
    },

    "grays": {
      "background": [
        "gray-950",
        "gray-900",
        "gray-800"
      ],
      "surfaces": [
        "gray-850 (custom)",
        "gray-800"
      ],
      "borders": "gray-700",
      "text": {
        "primary": "gray-100",
        "secondary": "gray-400",
        "muted": "gray-500"
      }
    },

    "accents": {
      "success": "green-500",
      "warning": "amber-400",
      "danger": "red-500",
      "info": "blue-400"
    },

    "lightModeOverrides": {
      "background": "gray-50",
      "surface": "white",
      "text": {
        "primary": "gray-900",
        "secondary": "gray-600"
      },
      "borders": "gray-200"
    }
  },

  "typography": {
    "scale": {
      "pageTitle": "text-xl font-semibold",
      "sectionTitle": "text-lg font-medium",
      "cardValue": "text-2xl font-semibold",
      "body": "text-sm",
      "caption": "text-xs"
    },
    "numericalData": {
      "font": "Tabular numerals",
      "weight": "Medium"
    }
  },

  "interactionPatterns": {
    "hover": "Subtle background elevation",
    "focus": "Soft ring using primary color",
    "active": "Slight inset / darker surface",
    "transitions": "Fast, ease-out, non-distracting"
  },

  "responsiveBreakpoints": {
    "sm": {
      "behavior": [
        "Sidebar becomes drawer",
        "Cards stack vertically",
        "Charts stack"
      ]
    },
    "md": {
      "behavior": [
        "2-column card layouts",
        "Table becomes horizontally scrollable"
      ]
    },
    "lg": {
      "behavior": [
        "Full sidebar",
        "4-column stat grid",
        "Side-by-side charts"
      ]
    },
    "xl": {
      "behavior": [
        "Increased max-width",
        "More whitespace",
        "Denser data visibility"
      ]
    }
  }
}
```

## Design tokens:
- Primary palette: used for charts, focus states, and informational emphasis.
- Secondary palette: used for success states and connectivity indicators.
- Grays: define background layers, surfaces, borders, and text hierarchy.
- Accents: reserved for status indicators and directional metrics (up/down traffic).

## Focus on:
- Visual consistency
- Information hierarchy
- Subtle interaction feedback
- Readability of numerical data
- Clear distinction between active, inactive, and disabled states

Do not introduce new features, controls, or workflows.
Your output should be a purely visual and structural enhancement of the existing UI.
If something is missing and is blocking the progress like missing APIs:

- use mock data (the ones in jsonc config will do just fine)
- add the backlog in this document
