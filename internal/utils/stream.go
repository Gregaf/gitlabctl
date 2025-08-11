package utils

import (
	"io"
)

func GetOutputStream(isEnabled bool, enabledWriter, disabledWriter io.Writer) io.Writer {
	if !isEnabled {
		return disabledWriter
	}

	return enabledWriter
}
