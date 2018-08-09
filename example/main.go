package main

import (
	"flag"
	"fmt"

	"github.com/vinely/layers"
	"github.com/vinely/layers/squashfs"
)

var (
	path = flag.String("path", "./test", "The Path to pack")
)

func main() {
	flag.Parse()

	if !squashfs.CheckSquashfsTools() {
		fmt.Println("Din't find squashfs tools!")
		return
	}

	//out := MakeSquashfsPackage(*path, "./xyz.sb")
	// out := squashfs.ExtraSquashfsPackage("./xyz.sb", "./xyz")
	// fmt.Printf("%v\n", out)

	ly := layers.PackPath(*path)
	if ly != nil {
		fmt.Printf("Layer: %+v\n", ly)
	}
	res := layers.VerifyLayers(ly.Location.ChkFile)
	fmt.Println("Verify result is ", res)

}
