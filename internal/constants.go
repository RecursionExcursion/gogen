package internal

/* Supported os and arch */

const Win64 = "win64"
const Linux64 = "linux64"
const Mac64 = "mac64"

type compliationPair struct {
	Os   string
	Arch string
}

var SupportedArchitecture = map[string]compliationPair{
	Win64:   {Os: "windows", Arch: "amd64"},
	Linux64: {Os: "linux", Arch: "arm64"},
	Mac64:   {Os: "darwin", Arch: "arm64"},
}
