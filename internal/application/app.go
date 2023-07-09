package application

import (
	"log"
	"os"

	"github.com/jinzhu/copier"

	"github.com/MarcosViniciusPinho/versionup/internal/application/dto"
	"github.com/MarcosViniciusPinho/versionup/internal/domain"
	"github.com/MarcosViniciusPinho/versionup/internal/domain/service"
)

func Start() {
	destinationFile := CopyFile("versionup_old.yml", "versionup.yml")
	defer destinationFile.Close()

	entryDataDTO, err := dto.NewEntryDataDTO()
	if err != nil {
		log.Fatal(err)
	}

	entryDataDomain := domain.EntryData{}

	copier.Copy(&entryDataDomain, &entryDataDTO)

	service.NewVersionServicePort().Update(destinationFile, entryDataDomain)
}

func CopyFile(originFileName, destinationFileName string) *os.File {
	originalFile, err := os.ReadFile(originFileName)
	if err != nil {
		log.Fatalf("Error reading source file: %v", err)
	}

	err = os.WriteFile(destinationFileName, originalFile, 0644)
	if err != nil {
		log.Fatal(err)
		log.Fatalf("Error writing to destination file: %v", err)
	}

	destinationFile, err := os.Open(destinationFileName)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	return destinationFile
}
