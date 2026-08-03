# Finder Task List

## Completed

- [x] Add a visible landing-page **Open File Finder** button.
- [x] Add a clickable `[x]` close control.
- [x] Keep the finder responsive with bounded dimensions.
- [x] Prevent long paths, previews, metadata, and clipboard text from wrapping the panel.
- [x] Add keyboard, mouse-wheel, page, home, and end navigation.
- [x] Add a block scrollbar to the file list.
- [x] Add internal pane separators and a footer divider.
- [x] Use bold neon-green finder title text.
- [x] Use a high-contrast light-pink selected-file row.
- [x] Use an uppercase neon-pink `[X]` close button aligned to the right.
- [x] Keep the scrollbar track intact when filenames are wider than the list.
- [x] Prevent styled pane rows from wrapping and fragmenting the scrollbar.
- [x] Keep the `[X]` button on the same right-aligned title row at narrow sizes.
- [x] Hide dot-files by default with `.` to toggle them back on.
- [x] Derive the close hitbox from the final fixed panel rectangle.
- [x] Keep the close control inside the panel's upper-right corner.
- [x] Normalize carriage returns before clipping preview rows.
- [x] Draw list and preview panes independently inside hard width boundaries.
- [x] Prevent RTF/control-heavy previews from reflowing across the center divider.
- [x] Prevent long first preview lines from wrapping into adjacent rows.
- [x] Add an inner right safety margin so preview text cannot touch the outer border.
- [x] Add a top pane header rule to make the list and information bounds explicit.
- [x] Use the available terminal height while the Finder is open to cover the large banner and maximize the preview viewport.
- [x] Reserve a deliberate three-cell gutter before the outer right frame.
- [x] Keep the scrollbar origin aligned after adding the pane-header row.
- [x] Add pane-width tests for normal and very narrow terminal content areas.
- [x] Document the hard-rectangle rendering and clipping contract for future maintenance.
- [x] Move directory and preview reads off the UI update path with stale-result protection.
- [x] Add Finder refresh (`r`) plus loading and error states.
- [x] Rename the command menu entry to **Toggle Beast Mode**.

## Follow-up

- [x] Add richer file actions and explicit context handoff.
- [x] Add render-bound tests using a fake terminal screen at multiple sizes.
- [x] Move finder styles into the active theme.

## Maintenance / rollback notes

- Keep `finderPaneWidths` as the single source of truth for left pane, divider,
  preview pane, and list-track widths. If the right border regresses, inspect
  the `rightMargin` before changing preview text handling.
- Keep `fixedLines` as the last normalization boundary before styling. New file
  formats should not bypass it or be rendered as an unconstrained styled row.
- If the Finder becomes too tall on a future layout, adjust the outer `height`
  inset in `FileBrowser.Draw`; do not restore a small fixed height without
  checking short terminals and banner overlap.
- If mouse selection or the scrollbar shifts vertically, compare the body origin
  (`titleH + paneHeaderH`) in `Draw` with `contentRect` and scrollbar drawing.
- Validate changes with representative `.rtf`, `.txt`, `.md`, `.sh`, URL-heavy,
  and long-first-line files at both wide and narrow terminal sizes.