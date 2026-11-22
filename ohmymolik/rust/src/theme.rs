//! Named UI color tokens. These palettes do not yet affect an agent TUI.

#[derive(Clone, Copy, Debug, PartialEq, Eq)]
pub struct Theme {
    pub name: &'static str,
    pub background: &'static str,
    pub surface: &'static str,
    pub text: &'static str,
    pub muted: &'static str,
    pub accent: &'static str,
    pub success: &'static str,
    pub warning: &'static str,
    pub error: &'static str,
}

pub const THEMES: &[Theme] = &[
    Theme {
        name: "forest",
        background: "#0B211B",
        surface: "#18362B",
        text: "#E4F4E8",
        muted: "#A8C7B3",
        accent: "#A9E6C0",
        success: "#A7DA9F",
        warning: "#F1D395",
        error: "#F4B1A8",
    },
    // Catppuccin Mocha palette: https://github.com/catppuccin/catppuccin (MIT).
    Theme {
        name: "catppuccin-mocha",
        background: "#1E1E2E",
        surface: "#313244",
        text: "#CDD6F4",
        muted: "#A6ADC8",
        accent: "#CBA6F7",
        success: "#A6E3A1",
        warning: "#F9E2AF",
        error: "#F38BA8",
    },
    Theme {
        name: "ocean",
        background: "#102332",
        surface: "#203C4B",
        text: "#E2F1F5",
        muted: "#A9C5CE",
        accent: "#8ED5E0",
        success: "#A9DDBE",
        warning: "#F3D596",
        error: "#F3B1AD",
    },
    Theme {
        name: "dusk",
        background: "#272333",
        surface: "#393449",
        text: "#F0EAF5",
        muted: "#C7B8D0",
        accent: "#D8B8E7",
        success: "#B8DDB9",
        warning: "#EFD199",
        error: "#F2B0B4",
    },
    Theme {
        name: "dawn",
        background: "#F9F6F0",
        surface: "#EDE7DC",
        text: "#252E29",
        muted: "#4D5E53",
        accent: "#176048",
        success: "#24633B",
        warning: "#805100",
        error: "#9C3D3D",
    },
];

pub fn find(name: &str) -> Option<&'static Theme> {
    THEMES.iter().find(|theme| theme.name == name)
}
