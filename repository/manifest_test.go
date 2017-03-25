package repository

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

var unzippedManifestFile = "unzipped-primary.xml"

var unzippedManifestData []byte

func init() {
	fixtureDir, err := filepath.Abs(fixtureDir)
	if err != nil {
		panic(err)
	}
	unzippedManifestData, err = ioutil.ReadFile(
		filepath.Join(fixtureDir, "repodata", unzippedManifestFile),
	)
	if err != nil {
		panic(err)
	}
}

type testManifest struct {
	data         *[]byte
	modType      string
	modName      string
	modArch      string
	modChecksum  string
	modVersion   string
	modPackager  string
	modLocation  string
	modSignature string
	modPlatform  string
}

var manifesttest = testManifest{
	&unzippedManifestData,
	"lime",
	"lime-4.2.0-17-generic.ko",
	"x86_64",
	"8fd9d9c765bac68763d4741d4726e9b120bffe1aaa7df7949aa94b37f7a6b6f8",
	"4.2.0-17-generic",
	"lime-compiler",
	"modules/lime-4.2.0-17-generic.ko",
	"modules/lime-4.2.0-17-generic.ko.sig",
	"linux",
}

func TestModuleManifest(t *testing.T) {
	man := moduleManifest(*manifesttest.data)

	if len(man.Modules) != 1 {
		t.Error(
			"For\n", string(*manifesttest.data),
			"expected len(Manifest.Modules) == 1",
			"got", len(man.Modules),
		)
	}
	mod := man.Modules[0]

	if mod.ModuleType != manifesttest.modType {
		t.Error(
			"For\n", string(*manifesttest.data),
			"expected type", manifesttest.modType,
			"got", mod.ModuleType,
		)
	}

	if mod.Name != manifesttest.modName {
		t.Error(
			"For\n", string(*manifesttest.data),
			"expected name", manifesttest.modName,
			"got", mod.Name,
		)
	}

	if mod.Arch != manifesttest.modArch {
		t.Error(
			"For\n", string(*manifesttest.data),
			"expected arch", manifesttest.modArch,
			"got", mod.Arch,
		)
	}

	if mod.Checksum != manifesttest.modChecksum {
		t.Error(
			"For\n", string(*manifesttest.data),
			"expected checksum", manifesttest.modChecksum,
			"got", mod.Checksum,
		)
	}

	if mod.Version != manifesttest.modVersion {
		t.Error(
			"For\n", string(*manifesttest.data),
			"expected version", manifesttest.modVersion,
			"got", mod.Version,
		)
	}

	if mod.Packager != manifesttest.modPackager {
		t.Error(
			"For\n", string(*manifesttest.data),
			"expected packager", manifesttest.modPackager,
			"got", mod.Packager,
		)
	}

	if mod.Location.Href != manifesttest.modLocation {
		t.Error(
			"For\n", string(*manifesttest.data),
			"expected location", manifesttest.modLocation,
			"got", mod.Location.Href,
		)
	}

	if mod.Signature.Href != manifesttest.modSignature {
		t.Error(
			"For\n", string(*manifesttest.data),
			"expected signature", manifesttest.modSignature,
			"got", mod.Signature.Href,
		)
	}

	if mod.Platform != manifesttest.modPlatform {
		t.Error(
			"For\n", string(*manifesttest.data),
			"expected platform", manifesttest.modPlatform,
			"got", mod.Platform,
		)
	}
}

type findTest struct {
	version string
	results int
}

var findtests = []findTest{
	{
		"4.2.0-17-generic",
		1,
	},
	{
		"*-generic",
		1,
	},
	{
		"4.2.0*",
		1,
	},
	{
		"*.x86_64",
		0,
	},
}

func TestFind(t *testing.T) {
	man := moduleManifest(*manifesttest.data)
	for _, input := range findtests {
		modules, err := man.find(input.version)
		if len(modules) != input.results {
			t.Error(
				"For", input.version,
				"expected results:", input.results,
				"got", len(modules),
				"with err", err,
			)
		}
	}
}
