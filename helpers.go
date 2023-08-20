package celeritas

import "os"

func (c *Celeritas) CreateDirIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)

		if err != nil {
			return err
		}
	}

	return nil
}
