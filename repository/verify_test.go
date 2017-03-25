package repository

import "testing"

type checksuminput struct {
	data     []byte
	checksum string
	valid    bool
}

var checksumtests = []checksuminput{
	{
		[]byte("passing checksum test value"),
		"a4fe0d179401df2c292d6a013f9d30521f486185923981e5accaaa20e4a44b7e",
		true,
	},
	{
		[]byte("failing checksum test value"),
		"1886e7b285e9d1a29f2a208dce7d7586aefec980b07cb68026f689bea4b52133",
		false,
	},
}

func TestSha256sum(t *testing.T) {
	for _, input := range checksumtests {
		valid, calcSum := sha256sum(input.data, input.checksum)
		if valid != input.valid {
			t.Error(
				"For data:", string(input.data), "(cast from []byte),",
				"and checksum", input.checksum,
				"expected valid?", input.valid,
				"got valid?", valid,
				"with checksum", calcSum,
			)
		}
	}
}
