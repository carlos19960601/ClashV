package constant

var Path = func() *path {
	return &path{}
}()

type path struct {
	homeDir string
}

func SetHomeDir(root string) {
	Path.homeDir = root
}
