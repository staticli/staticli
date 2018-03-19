package lib

import (
	"fmt"
	"net/http"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/inconshreveable/go-update"
	pb "gopkg.in/cheggaaa/pb.v2"
)

var (
	Version     string
	BuildTime   string
	BuildCommit string
)

// PrintVersion prints the version information
func PrintVersion() {
	log.Infof("staticli v%s", Version)

	log.WithFields(log.Fields{
		"Version":   Version,
		"BuildTime": BuildTime,
	}).Debugf("Build Details")

	log.WithFields(log.Fields{
		"BuildCommit": BuildCommit,
	}).Debugf("Git Details")

}

// IsLatestVersion returns a bool indicating whether if the current version is
// the latest, in addition to some information from the Github Release API
// about the latest version
func IsLatestVersion() (bool, GithubRelease, error) {
	releaseUrl := "https://api.github.com/repos/staticli/staticli/releases/latest"

	releaseData := GithubRelease{}
	GetJson(releaseUrl, &releaseData)

	if releaseData.Name == "" {
		return false, releaseData, fmt.Errorf("Couldn't get latest version of staticli, do you have an internet connection?")
	}

	log.Debugf("Latest version is: %s", releaseData.Name)

	if Version != releaseData.Name {
		return false, releaseData, nil
	}

	return true, releaseData, nil
}

// GithubRelease is a small subset of the data returned by the Github Releases API
type GithubRelease struct {
	Name   string
	Assets []GithubReleaseAsset
}

// GithubReleaseAsset is a small subset of the data corresponding to an asset
// in a response from the GitHub Releases API
type GithubReleaseAsset struct {
	Size        int
	Name        string
	DownloadUrl string `json:"browser_download_url"`
}

// updateAsset returns the specific GithubReleaseAsset which corresponds to the
// current Operating System / Archictecture
func updateAsset(releaseData GithubRelease) (GithubReleaseAsset, error) {

	updateName := fmt.Sprintf("staticli.%s.%s", runtime.GOOS, runtime.GOARCH)

	for _, asset := range releaseData.Assets {
		if asset.Name == updateName {
			return asset, nil
		}
	}

	return GithubReleaseAsset{}, fmt.Errorf("Unable to find release asset for %s %s", runtime.GOOS, runtime.GOARCH)
}

// Update updates staticli to the specific GithubRelease
func Update(releaseData GithubRelease) error {

	updateAsset, err := updateAsset(releaseData)
	if err != nil {
		return fmt.Errorf("Unable to get update URL: %s", err)
	}

	log.Infof("Updating to v%s", releaseData.Name)

	resp, err := http.Get(updateAsset.DownloadUrl)
	if err != nil {
		return fmt.Errorf("Could not download new version: %s", err)
	}
	defer resp.Body.Close()

	// start new bar
	bar := pb.New(updateAsset.Size)
	bar.Start()
	// create proxy reader
	body := bar.NewProxyReader(resp.Body)

	err = update.Apply(body, update.Options{})
	if err != nil {
		return fmt.Errorf("Could not apply update: %s", err)
	}

	bar.Finish()
	log.Infof("Up to date! 🎉")
	return nil
}