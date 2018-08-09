package layers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/vinely/layers/crypto"
	"github.com/vinely/layers/squashfs"
)

type LayerRepo struct {
	// Version string
	// Author  string
	// Commit  string
	Created time.Time
}

type LayerConfig struct {
	ConfigFile string
}
type LayerEnv struct {
	Path string
}

type Layer struct {
	Id       string
	Bolb     string
	Chksum   string
	Location LayerLocation
	Repo     LayerRepo
	Config   LayerConfig
	Env      LayerEnv
}

var (
	TempFileName    = "tmp.sb"
	PackageFileName = "metadata"
)

func MakeLayer(srcPath string, ll *LayerLocation) *Layer {
	file := filepath.Join(os.TempDir(), TempFileName)
	squashfs.MakeSquashfsPackage(srcPath, file)
	sha256, err := crypto.Sha256sum(file)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	target := filepath.Join(ll.Path, "/"+sha256, PackageFileName)
	os.MkdirAll(filepath.Dir(target), 0777)
	if err := os.Rename(file, target); err != nil {
		// may be meet: The system cannot move the file to a different disk drive.
		// try MoveFile func
		err = squashfs.MoveFile(file, target)
		if err != nil {
			log.Println(err.Error())
			return nil
		}
	}

	s := fmt.Sprintf("%s  %s\n", sha256, target)

	squashfs.AppendFile(ll.ChkFile, s)
	return &Layer{
		Id:     "sha256:" + sha256,
		Bolb:   filepath.Base(target),
		Chksum: sha256,
		Repo: LayerRepo{
			Created: time.Now(),
		},
		Location: LayerLocation{
			Path:    ll.Path,
			ChkFile: ll.ChkFile,
		},
	}
}

func MakeLayerSimple(srcPath, chkFile string) *Layer {
	ll := &LayerLocation{
		Path:    filepath.Dir(chkFile),
		ChkFile: chkFile,
	}
	return MakeLayer(srcPath, ll)
}

func PackPath(srcPath string) *Layer {
	ll := GetLayerLocation()
	return MakeLayer(srcPath, ll)
}

func VerifyLayers(sumFile string) bool {
	return crypto.Check256sumFromFile(sumFile)
}
