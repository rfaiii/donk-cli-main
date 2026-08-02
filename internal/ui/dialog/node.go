package dialog

import "github.com/rfaiii/donk-cli-main/internal/node"

// NodeSettingsID is the identifier for the node settings dialog.
const NodeSettingsID = "node-settings"

// ActionNodeRename requests renaming the selected node.
type ActionNodeRename struct {
	DeviceID string
	Nickname string
}

// ActionNodeUpdateStatus requests updating a node status.
type ActionNodeUpdateStatus struct {
	DeviceID string
	Status   node.DeviceStatus
}

// NodeSettingsUpdateMsg is emitted when the node registry changes while the
// settings dialog is open.
type NodeSettingsUpdateMsg struct {
	Devices []node.Device
}

// NodeSettingsConfirmMsg is emitted when the settings dialog is closed with
// a durable selection.
type NodeSettingsConfirmMsg struct {
	SelectedDeviceID string
}

// NodeSettingsCloseMsg closes the dialog.
type NodeSettingsCloseMsg struct{}

// NodeRenamePromptMsg requests the rename input for a device.
type NodeRenamePromptMsg struct {
	DeviceID string
	Current  string
}

// ToNodeRenameAction converts a NodeRenamePromptMsg into an action for the
// command palette/dialog dispatcher.
func (msg NodeRenamePromptMsg) ToNodeRenameAction() Action {
	return ActionOpenDialog{DialogID: NodeSettingsID}
}

// ToNodeSettingsOpenAction opens the node settings dialog.
func ToNodeSettingsOpenAction() Action {
	return ActionOpenDialog{DialogID: NodeSettingsID}
}

// NodeSlashCommand returns the slash command string for node settings.
func NodeSlashCommand() string {
	return "/node"
}
