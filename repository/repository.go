package repository

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Repository struct {
	BaseUrl       string
	SkipGPGVerify bool
	metaDir       string
	repoMeta      string
	repoMetaSig   string
	signingKey    string
}

type RepoMetadata struct {
	Revision string           `xml:"revision"`
	Manifest ManifestMetadata `xml:"data"`
}

type ManifestMetadata struct {
	RepoType     string   `xml:"type,attr"`
	Checksum     string   `xml:"checksum"`
	OpenChecksum string   `xml:"open_checksum"`
	Location     Location `xml:"location"`
	Timestamp    string   `xml:"timestamp"`
	Size         int      `xml:"size"`
	OpenSize     int      `xml:"open_size"`
}

var netClient *http.Client

func init() {
	netClient = &http.Client{
		Timeout: time.Second * 60,
	}
}

func repoMetadata(data []byte) RepoMetadata {
	var repo RepoMetadata
	xml.Unmarshal(data, &repo)
	return repo
}

func DefaultRepository() Repository {
	return Repository{
		"https://threatresponse-lime-modules.s3.amazonaws.com/",
		false,
		"repodata/",
		"repomd.xml",
		"repomd.xml.sig",
		"REPO_SIGNING_KEY.asc",
	}
}

func (r *Repository) Get(kernVer string) int {

	modules, err := r.Find(kernVer)
	if err != nil {
		log.Critical(err)
		return 0
	}

	//TODO: prompt if there are multiple matches
	//TODO: implement a multi-get function that downloads all matches
	// Exit if there are multiple matches for a kernel version
	if len(modules) != 1 {
		log.Critical(fmt.Sprintf("multiple matches for: %s", kernVer))
		return 1
	}
	log.Debug(fmt.Sprintf("found module matching: %s", kernVer))

	localPath, err := r.download(modules[0])
	if err != nil {
		log.Critical(err)
		return 0
	}
	log.Info(fmt.Sprintf("module downloaded to %s", localPath))

	return 1
}

func (r *Repository) Find(kernVer string) ([]Module, error) {

	var modules []Module
	manifest, err := r.manifest()
	if err != nil {
		return modules, err
	}

	modules, err = manifest.find(kernVer)
	if err != nil {
		return modules, err
	}

	return modules, nil
}

func (r *Repository) List() (Manifest, error) {
	manifest, err := r.manifest()
	return manifest, err
}

func (r *Repository) manifest() (Manifest, error) {
	repo, err := r.metadata()
	if err != nil {
		return Manifest{}, err
	}

	manifest, err := r.fetchManifest(repo)
	if err != nil {
		return Manifest{}, err
	}

	return manifest, nil
}

func (r *Repository) metadata() (RepoMetadata, error) {
	// fetch repository metadata file
	var url string
	url = fmt.Sprintf("%s%s%s", r.BaseUrl, r.metaDir, r.repoMeta)
	log.Debug(fmt.Sprintf("fetching repo metadata: %s", url))

	resp, err := netClient.Get(url)
	if err != nil {
		return RepoMetadata{},
			errors.New(fmt.Sprintf("unable to fetch repository metadata: %s", err))
	}
	defer resp.Body.Close()
	rawMetadata, err := ioutil.ReadAll(resp.Body)

	if r.SkipGPGVerify == false {
		url = fmt.Sprintf("%s%s", r.BaseUrl, r.signingKey)
		log.Debug(fmt.Sprintf("fetching repo signing key: %s", url))
		resp, err = netClient.Get(url)
		if err != nil {
			return RepoMetadata{},
				errors.New(fmt.Sprintf("error fetching repository signing key: %s", err))
		}
		defer resp.Body.Close()

		repoKey, err := readKey(resp.Body)
		if err != nil {
			return RepoMetadata{},
				errors.New(fmt.Sprintf("error reading repository signing key: %s", err))
		}

		// load user's keyring
		keyring, err := getDefaultKeyring()
		if err != nil {
			return RepoMetadata{},
				errors.New(fmt.Sprintf("error loading user keyring: %s", err))
		}

		//check if repo key is imported to user keychain
		//TODO: expand info in error message
		if !keyring.contains(repoKey) {
			return RepoMetadata{},
				errors.New("Repository key not imported in user keychain")
		}

		// fetch detached repository metadata signature
		url = fmt.Sprintf("%s%s%s", r.BaseUrl, r.metaDir, r.repoMetaSig)
		log.Debug(fmt.Sprintf("fetching repo metadata signature: %s", url))
		resp, err = netClient.Get(url)
		if err != nil {
			return RepoMetadata{},
				errors.New(fmt.Sprintf("error fetching repo metadata signature: %s", err))
		}
		defer resp.Body.Close()

		metadataReader := bytes.NewReader(rawMetadata)
		signer, err := keyring.verifyDetachedSig(metadataReader, resp.Body)
		if err != nil {
			return RepoMetadata{},
				errors.New(fmt.Sprintf("error verifying repo metadata signature: %s", err))
		}
		log.Debug(fmt.Sprintf("verified metadata signature against %s", signer.fingerprint()))
	}
	return repoMetadata(rawMetadata), nil
}

func (r *Repository) fetchManifest(repo RepoMetadata) (Manifest, error) {
	//Download manifest from repository
	url := fmt.Sprintf("%s%s", r.BaseUrl, repo.Manifest.Location.Href)
	log.Debug(fmt.Sprintf("fetching manifest: %s", url))
	resp, err := netClient.Get(url)
	if err != nil {
		return Manifest{}, errors.New(
			fmt.Sprintf("unable to fetch repository manifest: %s", err),
		)
	}
	defer resp.Body.Close()

	// verify gzipped file checksum
	gzBody, err := ioutil.ReadAll(resp.Body)
	valid, calcSum := sha256sum(gzBody, repo.Manifest.Checksum)
	if valid == false {
		return Manifest{},
			errors.New(
				fmt.Sprintf(
					"manifest checksum mismatch expected: %s found: %s",
					repo.Manifest.OpenChecksum, calcSum,
				),
			)
	}

	//Unzip manifest
	//TODO: check for errors
	buf := bytes.NewBuffer(gzBody)
	reader, err := gzip.NewReader(buf)
	defer reader.Close()
	data, err := ioutil.ReadAll(reader)

	// verify manifest open checksum
	valid, calcSum = sha256sum(data, repo.Manifest.OpenChecksum)
	if valid == false {
		return Manifest{},
			errors.New(
				fmt.Sprintf(
					"manifest open checksum mismatch expected: %s found: %s",
					repo.Manifest.OpenChecksum, calcSum,
				),
			)
	}

	// create manifest object
	//TODO: add error handling
	return moduleManifest(data), nil
}

func (r *Repository) download(mod Module) (string, error) {
	var url string
	url = fmt.Sprintf("%s%s", r.BaseUrl, mod.Location.Href)
	log.Debug(fmt.Sprintf("downloading module from: %s", url))
	resp, err := netClient.Get(url)
	if err != nil {
		return "", err
	}
	log.Debug(fmt.Sprintf("get module returned %s", resp.Status))
	defer resp.Body.Close()

	modFile, err := os.Create(mod.Name)
	if err != nil {
		return "", err
	}

	defer modFile.Close()
	io.Copy(modFile, resp.Body)

	return mod.Name, nil
}
