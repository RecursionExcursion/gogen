package internal

/* Supported os and arch */

const Win64 = "win64"
const Linux64 = "linux64"
const Mac64 = "mac64"

type CompliationPair struct {
	Os   string
	Arch string
}
