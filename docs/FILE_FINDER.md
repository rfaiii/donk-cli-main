# DONK File Finder

DONK includes an in-app project file finder inspired by the focused panes in
[Superfile](https://github.com/yorukot/superfile). It is implemented as a DONK
dialog rather than embedding Superfile's standalone application, which keeps
keyboard focus, mouse events, theming, and the agent session in one Bubble Tea
program.

Open it from the `/` command palette by choosing **Open File Finder**, or type
`/finder` to filter directly to it. `ctrl+shift+f` is a quick toggle that opens
the finder and closes it when it is already open. The
finder provides:

- a project directory listing with mouse-clickable rows;
- a live text preview in the adjacent pane;
- selected-file metadata and clipboard status below the panes;
- path copying with `y` (and the clipboard pane reflects the copied path).

Keys: `↑`/`↓` or `j`/`k` navigate, `pgup`/`pgdn` scroll by a page, `home`/`g`
go to the first entry, and `end` goes to the last entry. `enter`/`l` opens a
directory, `h`/`left` or backspace goes to the parent directory, `y` copies the
selected path, `.` toggles hidden dot-files, and `esc` closes the finder.

The finder uses a bounded, responsive panel and reflows with terminal resizes.
The block scrollbar beside the file list shows the current position when more
entries are available. It can be opened by clicking the **📁 OPEN FILE FINDER**
button on the landing screen or the compact folder icon in the header. Click
**[X]** inside the upper-right corner to close it. The finder keeps its current
directory and selection when closed and reopened.

The selected file uses a light-pink highlight for visibility. The file list and
preview are separated by a divider, and the footer uses distinct colors for
metadata, clipboard state, and keyboard instructions. Current implementation
work is tracked in [`FINDER_TASKLIST.md`](FINDER_TASKLIST.md).

Rows are width-constrained before styling, so long filenames cannot wrap the
file list or break the scrollbar into fragments. The title close control stays
on the same row even on narrow terminals.

Processes and pinned tabs are intentionally omitted so the finder remains a
compact companion to the conversation UI.