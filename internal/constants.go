package internal

/* Supported os and arch */

const (
	Win64   = "win64"
	Linux64 = "linux64"
	Mac64   = "mac64"
)

type CompilationPair struct {
	Os   string
	Arch string
}
