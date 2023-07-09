package domain

type EntryData struct {
	Enviroment     string `copier:"Enviroment"`
	TypeVersion    string `copier:"TypeVersion"`
	DescriptionTag string `copier:"DescriptionTag"`
	RepositoryUrl  string `copier:"RepositoryUrl"`
	UserName       string `copier:"UserName"`
	UserEmail      string `copier:"UserEmail"`
}
