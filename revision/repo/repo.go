package repo

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/wallix/awless/config"
	"github.com/wallix/awless/rdf"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type Rev struct {
	Id   string
	Date time.Time

	Infra  *rdf.Graph
	Access *rdf.Graph
}

type Repo interface {
	Commit(files ...string) error
	LoadRev(version string) (*Rev, error)
}

type NoRevisionRepo struct{}

func (*NoRevisionRepo) Commit(files ...string) error         { return nil }
func (*NoRevisionRepo) LoadRev(version string) (*Rev, error) { return &Rev{}, nil }

type GitRepo struct {
	repo  *git.Repository
	files []string
	path  string
}

func NewRepo() (Repo, error) {
	if IsGitInstalled() {
		return newGitRepo(config.RepoDir)
	} else {
		return &NoRevisionRepo{}, nil
	}
}

func IsGitInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func newGitRepo(path string) (Repo, error) {
	if _, err := os.Stat(filepath.Join(path, ".git")); os.IsNotExist(err) {
		if _, err := newGit(path).run("init"); err != nil {
			return nil, err
		}
	}

	repo, err := git.NewFilesystemRepository(filepath.Join(path, ".git"))
	return &GitRepo{repo: repo, path: path}, err
}

func (r *GitRepo) LoadRev(version string) (*Rev, error) {
	rev := &Rev{Id: version}

	commit, err := r.repo.Commit(plumbing.NewHash(version))
	if err != nil {
		return nil, err
	}

	rev.Date = commit.Committer.When

	f, err := commit.File(config.InfraFilename)
	if err != nil {
		return nil, err
	}
	contents, err := f.Contents()
	if err != nil {
		return nil, err
	}
	infraG := rdf.NewGraph()
	infraG.Unmarshal([]byte(contents))
	rev.Infra = infraG

	f, err = commit.File(config.AccessFilename)
	if err != nil {
		return nil, err
	}
	contents, err = f.Contents()
	if err != nil {
		return nil, err
	}
	accessG := rdf.NewGraph()
	accessG.Unmarshal([]byte(contents))
	rev.Access = accessG

	return rev, nil
}

func (r *GitRepo) Commit(files ...string) error {
	for _, path := range files {
		r.files = append(r.files, path)
	}

	for _, path := range r.files {
		if _, err := newGit(r.path).run("add", path); err != nil {
			return err
		}
	}

	if hasChanges, err := r.hasChanges(); err != nil {
		return err
	} else if !hasChanges {
		return nil
	}

	if _, err := newGit(r.path).run("-c", "user.name='awless'", "-c", "user.email='git@awless.io'", "commit", "-m", "new sync"); err != nil {
		return err
	}

	return nil
}

func (r *GitRepo) hasChanges() (bool, error) {
	stdout, err := newGit(r.path).run("status", "--porcelain")
	if err != nil {
		return false, err
	}

	return !(strings.TrimSpace(stdout) == ""), nil
}
