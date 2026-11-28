# Oh My Molik identity

**Oh My Molik** is the independent Rust agent project by MoLik Labs, formerly
called Oh My Fox. FCC remains the name of the existing gateway. The historic Go
prototype keeps its path, binary and `OMF_API_KEY` for compatibility.

## The mark

`selected-reference.png` is the reference image supplied by the user.
`molik-mark.svg` is a clean, editable vector interpretation: an angular M with a
forward-facing cut. Use the SVG for scalable applications; the reference is
preserved for comparison. No separate trademark registration, ownership search
or legal clearance is asserted. Do not imply affiliation with another brand.

- Primary mark: pale mint `#DDFBE2` on forest `#0B3020`.
- Keep generous clear space around the mark; do not stretch or add effects.
- Use an empty `alt` when the adjacent text already names Oh My Molik; use
  `alt="Oh My Molik monogram"` when the mark stands alone.
- Product name: Oh My Molik. Rust executable: `oh-my-molik`. The Go prototype's
  older `ohmyfox` name and `OMF_API_KEY` are compatibility identifiers, not the
  new branding. Do not rename FCC or imply integration that does not yet exist.

## Theme collection

Five named design-token sets live in `../rust/src/theme.rs`:

| Theme | Character | Status |
| --- | --- | --- |
| Forest | Deep green, mint accent | Signature dark palette |
| Catppuccin Mocha | Soft charcoal and mauve | Community palette |
| Ocean | Deep blue and cyan | Dark palette |
| Dusk | Muted violet | Dark palette |
| Dawn | Warm neutral | Light palette |

The Mocha color values are from the [Catppuccin palette](https://github.com/catppuccin/catppuccin),
© 2021 Catppuccin, licensed under MIT; see [CATPPUCCIN-LICENSE](CATPPUCCIN-LICENSE).
Other palettes are project design tokens. Contrast tests check every foreground
against both base and surface at WCAG 2.x AA's 4.5:1 text ratio; they are not a
full interface accessibility audit. The gallery is an illustrative static
preview, and the Rust CLI can list tokens, but agent theme switching and a TUI
are not yet implemented. Separately, FCC Admin offers browser-local appearance
selection inspired by these palettes. Avoid claiming these mockups show a working product.

The repository's [AGPL-3.0-only license](../../LICENSE) applies to project code;
the Catppuccin notice is preserved separately.
