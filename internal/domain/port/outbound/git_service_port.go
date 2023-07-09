package outbound

import (
	"os"

	"github.com/MarcosViniciusPinho/versionup/internal/domain"
)

type IGitServicePort interface {
	CreateCommitAndTag(newFile *os.File, entryData domain.EntryData)
}
