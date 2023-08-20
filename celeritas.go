package celeritas

import "fmt"

const version = "1.0.0"

type Celeritas struct {
	AppName string
	Debug   bool
	Version string
}

func (c *Celeritas) New(rootPath string) error {
	pathConfing := initPaths{
		rootPath:    rootPath,
		folderNames: []string{"handlers", "migrations", "views", "data", "public", "tmp", "logs", "middleware"},
	}

	err := c.Init(pathConfing)
	if err != nil {
		return err
	}

	return err
}

func (c *Celeritas) Init(p initPaths) error {
	for _, path := range p.folderNames {
		err := c.CreateDirIfNotExist(fmt.Sprintf("%s/%s", p.rootPath, path))

		if err != nil {
			return err
		}
	}

	return nil
}
