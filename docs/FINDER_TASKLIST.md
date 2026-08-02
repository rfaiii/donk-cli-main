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
- [x] Rename the command menu entry to **Toggle Beast Mode**.

## Follow-up

- [ ] Move directory reads and file previews off the UI update path.
- [ ] Add a refresh/reload action for changed directories.
- [ ] Add render-bound tests using a fake terminal screen at multiple sizes.
- [ ] Consider theme-owned finder styles instead of local color constants.