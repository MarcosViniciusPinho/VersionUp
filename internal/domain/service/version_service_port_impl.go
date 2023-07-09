package service

import (
	"fmt"
	"log"
	"os"

	"github.com/MarcosViniciusPinho/versionup/internal/domain"
	"github.com/MarcosViniciusPinho/versionup/internal/domain/port/inbound"
	"github.com/MarcosViniciusPinho/versionup/internal/domain/port/outbound"
	"github.com/MarcosViniciusPinho/versionup/internal/infrastructure"
)

type VersionServicePort struct {
	GitService outbound.IGitServicePort
}

func (vs VersionServicePort) Update(file *os.File, entryData domain.EntryData) {
	var version domain.Version

	extractVersion := version.Extract(file)

	newVersion := vs.increment(extractVersion, entryData)
	vs.updateYml(
		file,
		extractVersion,
		newVersion,
		&entryData,
	)
}

func (vs VersionServicePort) updateYml(
	file *os.File,
	extractVersion *domain.Version,
	newVersion string,
	entryData *domain.EntryData,
) {
	newFile, err := os.Create(file.Name())
	if err != nil {
		log.Fatalf("Error creating YML file: %v", err)
	}
	defer newFile.Close()

	_, err = newFile.Write(extractVersion.ToByte())
	if err != nil {
		log.Fatalf("Error writing to YML file: %v", err)
	}

	fmt.Println("\nPlease wait while we are pushing the version change to the repository and creating your tag...")

	entryData.DescriptionTag = fmt.Sprintf("%s%s", entryData.DescriptionTag, newVersion)
	vs.GitService.CreateCommitAndTag(newFile, *entryData)

	fmt.Println("\nUpdated version successfully!")
	fmt.Println("Tag created:", entryData.DescriptionTag)

	extractVersion.Extract(newFile)
}

func (vs VersionServicePort) increment(
	version *domain.Version,
	entryData domain.EntryData,
) string {
	switch entryData.TypeVersion {
	case "major":
		return version.IncrementByEnviroment(0, entryData.Enviroment)
	case "minor":
		return version.IncrementByEnviroment(1, entryData.Enviroment)
	case "patch":
		return version.IncrementByEnviroment(2, entryData.Enviroment)
	default:
		panic("Invalid version type. Options: major, minor, patch")
	}
}

func NewVersionServicePort() inbound.IVersionServicePort {
	return &VersionServicePort{
		GitService: infrastructure.NewGitServicePort(),
	}
}
