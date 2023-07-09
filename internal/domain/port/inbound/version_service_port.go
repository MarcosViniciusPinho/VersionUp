package inbound

import (
	"os"

	"github.com/MarcosViniciusPinho/versionup/internal/domain"
)

type IVersionServicePort interface {
	Update(file *os.File, entryData domain.EntryData)
}
