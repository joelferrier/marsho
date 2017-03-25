package repository

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/ryanuber/go-glob"
)

type Manifest struct {
	Modules []Module `xml:"module"`
}

type Module struct {
	ModuleType string   `xml:"type,attr"`
	Name       string   `xml:"name"`
	Arch       string   `xml:"arch"`
	Checksum   string   `xml:"checksum"`
	Version    string   `xml:"version"`
	Packager   string   `xml:"packager"`
	Location   Location `xml:"location"`
	Signature  Location `xml:"signature"`
	Platform   string   `xml:"platform"`
}

func (m *Manifest) find(version string) ([]Module, error) {
	var modCollection []Module

	for _, mod := range m.Modules {
		if glob.Glob(version, mod.Version) {
			modCollection = append(modCollection, mod)
		}
	}

	if len(modCollection) > 0 {
		return modCollection, nil
	} else {
		return modCollection, errors.New(fmt.Sprintf("repository: module version %s not found", version))
	}
}

func moduleManifest(data []byte) Manifest {
	var manifest Manifest
	xml.Unmarshal(data, &manifest)
	return manifest
}
