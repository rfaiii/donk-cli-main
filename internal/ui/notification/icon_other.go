//go:build !darwin

package notification

import (
	_ "embed"
)

//go:embed donk-icon-solo.png
var Icon []byte
