package layers

import "path/filepath"

type LayerLocation struct {
	Path    string
	ChkFile string
}

var (
	Root       = "./"
	Sublay     = "./"
	ChksumFile = "sha256.sum"
)

func GetLayerLocation() *LayerLocation {
	return &LayerLocation{
		Path:    filepath.Join(Root, Sublay),
		ChkFile: filepath.Join(Root, Sublay, ChksumFile),
	}
}
