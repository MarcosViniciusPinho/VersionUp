package infrastructure

import (
	"log"
	"os"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"

	"github.com/MarcosViniciusPinho/versionup/internal/domain"
	"github.com/MarcosViniciusPinho/versionup/internal/domain/port/outbound"
)

type GitServicePort struct{}

func (gs GitServicePort) CreateCommitAndTag(newFile *os.File, entryData domain.EntryData) {
	commit, repo := gs.createCommit(entryData, newFile)
	gs.createTag(commit, repo, entryData.DescriptionTag)
}

func (gs GitServicePort) createCommit(entryData domain.EntryData, newFile *os.File) (plumbing.Hash, *git.Repository) {
	_, err := git.PlainClone(entryData.RepositoryUrl, false, &git.CloneOptions{
		URL: entryData.RepositoryUrl,
	})
	if err != nil {
		log.Fatalf("Error cloning project using SSH authentication: %v", err)
	}

	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatalf("Error opening project: %v", err)
	}

	wt, err := repo.Worktree()
	if err != nil {
		log.Fatalf("Error preparing environment for staging: %v", err)
	}

	_, err = wt.Add(newFile.Name())
	if err != nil {
		log.Fatalf("Error adding file to staging area: %v", err)
	}

	commit, err := wt.Commit("(versionUP) version update triggered", &git.CommitOptions{
		Author: gs.createSignature(entryData.UserName, entryData.UserEmail),
	})
	if err != nil {
		log.Fatalf("Error when committing: %v", err)
	}

	return commit, repo
}

func (gs GitServicePort) createTag(commit plumbing.Hash, repo *git.Repository, tagName string) {
	tagNameRef := plumbing.ReferenceName("refs/tags/" + tagName)

	r, err := repo.Storer.Reference(tagNameRef)
	if err == nil && r != nil {
		log.Fatalf("Tag %s already exists in repository\n", tagName)
	}

	if err := repo.Storer.SetReference(plumbing.NewHashReference(tagNameRef, commit)); err != nil {
		log.Fatal(err)
	}

	// Configurar autenticação SSH com a chave privada
	auth, err := ssh.NewPublicKeysFromFile("git", "/ssh/id_rsa", "")
	if err != nil {
		log.Fatalf("Error configuring ssh authentication: %v", err)
	}

	// Configurar o remote com a URL e a autenticação
	remote, err := repo.Remote("origin")
	if err != nil {
		log.Fatalf("Error configuring repository with ssh: %v", err)
	}

	// Criar o remote com a URL e a autenticação
	rconfig := remote.Config()
	rconfig.URLs = append(rconfig.URLs, remote.Config().URLs...)
	rconfig.URLs = append(rconfig.URLs, remote.Config().URLs...)
	rconfig.URLs = rconfig.URLs[:1] // Usar apenas a primeira URL

	// Push para o repositório remoto
	err = repo.Push(&git.PushOptions{
		RemoteName: rconfig.Name,
		Auth:       auth,
		Progress:   os.Stdout,
		RefSpecs:   []config.RefSpec{"refs/heads/*:refs/heads/*", "refs/tags/*:refs/tags/*"},
	})
	if err != nil {
		if err == transport.ErrEmptyRemoteRepository {
			log.Fatalf("The remote repository is empty: %v", err)
		} else if err == git.NoErrAlreadyUpToDate {
			log.Fatalf("The local repository is already up to date with the remote one: %v", err)
		} else {
			//TODO Investigate why the exception git.ErrNonFastForwardUpdate is not being caught in this part of the code.
			errMsg := err.Error()
			if errMsg == "non-fast-forward update: refs/heads/main" {
				log.Fatalf("Version has been updated and you haven't performed the local update yet. Use 'git pull origin' to update: %v", err)
			} else {
				log.Fatal(err)
			}
		}
	}
}

func (gs GitServicePort) createSignature(userName, userEmail string) *object.Signature {
	return &object.Signature{
		Name:  userName,
		Email: userEmail,
		When:  time.Now(),
	}
}

func NewGitServicePort() outbound.IGitServicePort {
	return &GitServicePort{}
}
