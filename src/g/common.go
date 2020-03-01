package g

import (
	"os"
)

func Mkdir(path string) error {
	return os.Mkdir(path, 0755)
}
