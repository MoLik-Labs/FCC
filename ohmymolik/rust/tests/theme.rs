use oh_my_molik::theme::{find, THEMES};
use std::collections::HashSet;

fn luminance(hex: &str) -> f64 {
    assert_eq!(hex.len(), 7);
    assert!(hex.starts_with('#'));
    let channel = |i| {
        let c = u8::from_str_radix(&hex[i..i + 2], 16).unwrap() as f64 / 255.0;
        if c <= 0.04045 {
            c / 12.92
        } else {
            ((c + 0.055) / 1.055).powf(2.4)
        }
    };
    0.2126 * channel(1) + 0.7152 * channel(3) + 0.0722 * channel(5)
}
fn contrast(a: &str, b: &str) -> f64 {
    let (x, y) = (luminance(a), luminance(b));
    (x.max(y) + 0.05) / (x.min(y) + 0.05)
}
#[test]
fn palettes_are_distinct_and_readable_on_background_and_surface() {
    let mut names = HashSet::new();
    for t in THEMES {
        assert!(names.insert(t.name));
        for background in [t.background, t.surface] {
            for foreground in [t.text, t.muted, t.accent, t.success, t.warning, t.error] {
                assert!(
                    contrast(foreground, background) >= 4.5,
                    "{} {} on {} contrast {}",
                    t.name,
                    foreground,
                    background,
                    contrast(foreground, background)
                );
            }
        }
    }
    assert!(find("forest").is_some());
    assert!(find("catppuccin-mocha").is_some());
    assert!(find("FOREST").is_none());
}
