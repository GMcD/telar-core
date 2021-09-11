package utils

import (
	"log"
	"runtime/debug"
	"strings"

	"github.com/icza/bitio"
)

func pkgversion(pkg string) string {
	_ = bitio.NewReader
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		log.Printf("Failed to read build info for package version.")
		return ""
	}

	for _, dep := range bi.Deps {
		if strings.Contains(dep.Path, pkg) {
			return dep.Version
		}
	}
	return ""
}
