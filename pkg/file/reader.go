package file

import (
	"bufio"
	"os"
)

func Load(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	list := make([]string, 0)

	s := bufio.NewScanner(f)
	for s.Scan() {
		list = append(list, s.Text())
	}

	return list, nil
}
