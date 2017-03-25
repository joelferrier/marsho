package repository

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

var fixtureDir = "./test-fixtures"
var repomdFile = "repomd.xml"

var repomdData []byte
var repo Repository

func init() {
	fixtureDir, err := filepath.Abs(fixtureDir)
	if err != nil {
		panic(err)
	}
	repomdData, err = ioutil.ReadFile(filepath.Join(fixtureDir, "repodata", repomdFile))
	if err != nil {
		panic(err)
	}

	repo = DefaultRepository()
}

type repomd struct {
	data             *[]byte
	revision         string
	dataType         string
	dataChecksum     string
	dataOpenChecksum string
	dataLocation     string
	dataTimestamp    string
	dataSize         int
	dataOpenSize     int
}

var repomdtest = repomd{
	&repomdData,
	"1487818901",
	"primary",
	"eba1bffc37262c1866cd5c76775181c507fe84c0321e166e2febda745e13545e",
	"0cf6f38bb35f9ea27158575b290ddcad4264ccbd3bdf01ffad9480ac24cb924f",
	"repodata/0cf6f38bb35f9ea27158575b290ddcad4264ccbd3bdf01ffad9480ac24cb924f-primary.xml.gz",
	"1487818901",
	23803,
	167164,
}

func TestRepoMetadata(t *testing.T) {
	metadata := repoMetadata(*repomdtest.data)

	if metadata.Revision != repomdtest.revision {
		t.Error(
			"For\n", string(*repomdtest.data),
			"expected revision", repomdtest.revision,
			"got", metadata.Revision,
		)
	}
	if metadata.Manifest.RepoType != repomdtest.dataType {
		t.Error(
			"For\n", string(*repomdtest.data),
			"expected type", repomdtest.dataType,
			"got", metadata.Manifest.RepoType,
		)
	}
	if metadata.Manifest.Checksum != repomdtest.dataChecksum {
		t.Error(
			"For\n", string(*repomdtest.data),
			"expected checksum", repomdtest.dataChecksum,
			"got", metadata.Manifest.Checksum,
		)
	}
	if metadata.Manifest.OpenChecksum != repomdtest.dataOpenChecksum {
		t.Error(
			"For\n", string(*repomdtest.data),
			"expected open checksum", repomdtest.dataOpenChecksum,
			"got", metadata.Manifest.OpenChecksum,
		)
	}
	if metadata.Manifest.Location.Href != repomdtest.dataLocation {
		t.Error(
			"For\n", string(*repomdtest.data),
			"expected location", repomdtest.dataLocation,
			"got", metadata.Manifest.Location.Href,
		)
	}
	if metadata.Manifest.Timestamp != repomdtest.dataTimestamp {
		t.Error(
			"For\n", string(*repomdtest.data),
			"expected", repomdtest.dataTimestamp,
			"got", metadata.Manifest.Timestamp,
		)
	}
	if metadata.Manifest.Size != repomdtest.dataSize {
		t.Error(
			"For\n", string(*repomdtest.data),
			"expected", repomdtest.dataSize,
			"got", metadata.Manifest.Size,
		)
	}
	if metadata.Manifest.OpenSize != repomdtest.dataOpenSize {
		t.Error(
			"For\n", string(*repomdtest.data),
			"expected", repomdtest.dataOpenSize,
			"got", metadata.Manifest.OpenSize,
		)
	}
}
