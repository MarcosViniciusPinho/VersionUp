package dto

import (
	"fmt"
	"os"
	"strings"
)

type EntryDataDTO struct {
	Enviroment     string `copier:"Enviroment"`
	TypeVersion    string `copier:"TypeVersion"`
	DescriptionTag string `copier:"DescriptionTag"`
	RepositoryUrl  string `copier:"RepositoryUrl"`
	UserName       string `copier:"UserName"`
	UserEmail      string `copier:"UserEmail"`
}

const (
	enviroment     = "ENVIROMENT"
	typeVersion    = "TYPE_VERSION"
	descriptionTag = "DESCRIPTION_TAG"
	repositoryUrl  = "REPOSITORY_URL"
	userName       = "USER_NAME"
	userEmail      = "USER_EMAIL"
)

func NewEntryDataDTO() (*EntryDataDTO, error) {
	entryDataDTO := &EntryDataDTO{
		Enviroment:     os.Getenv(enviroment),
		TypeVersion:    os.Getenv(typeVersion),
		DescriptionTag: os.Getenv(descriptionTag),
		RepositoryUrl:  os.Getenv(repositoryUrl),
		UserName:       os.Getenv(userName),
		UserEmail:      os.Getenv(userEmail),
	}
	return entryDataDTO.validate()
}

func (edd EntryDataDTO) validate() (*EntryDataDTO, error) {
	for _, f := range edd.getRequiredFields() {
		if strings.TrimSpace(f.field) == "" {
			return nil, fmt.Errorf("environment variable %s is required", f.fieldName)
		}
	}
	return &edd, nil
}

func (edd EntryDataDTO) getRequiredFields() []struct {
	field     string
	fieldName string
} {
	requiredFields := []struct {
		field     string
		fieldName string
	}{
		{edd.Enviroment, enviroment},
		{edd.TypeVersion, typeVersion},
		{edd.DescriptionTag, descriptionTag},
		{edd.RepositoryUrl, repositoryUrl},
		{edd.UserName, userName},
		{edd.UserEmail, userEmail},
	}
	return requiredFields
}
