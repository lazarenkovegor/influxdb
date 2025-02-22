//go:build !assets
// +build !assets

package static

import (
	"errors"
	"os"
)

// The functions defined in this file are placeholders for when the binary is
// built without assets.

var errNoAssets = errors.New("no assets included in binary")

// Asset returns an error stating no assets were included in the binary.
func Asset(string) ([]byte, error) {
	return nil, errNoAssets
}

// AssetInfo returns an error stating no assets were included in the binary.
func AssetInfo(name string) (os.FileInfo, error) {
	return nil, errNoAssets
}

// AssetDir returns nil because there are no assets included in the binary.
func AssetDir(name string) ([]string, error) {
	return nil, errNoAssets
}
