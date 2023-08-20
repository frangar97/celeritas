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

func (c *Celeritas) CreateFileIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)

		if err != nil {
			return err
		}

		defer func(file *os.File) {
			file.Close()
		}(file)

	}

	return nil
}
