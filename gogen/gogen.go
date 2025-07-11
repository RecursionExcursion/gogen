package gogen

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/RecursionExcursion/gogen/internal"
)

var SupportedArchitecture = map[string]internal.CompilationPair{
	internal.Win64:   {Os: "windows", Arch: "amd64"},
	internal.Linux64: {Os: "linux", Arch: "arm64"},
	internal.Mac64:   {Os: "darwin", Arch: "arm64"},
}

type CreateExeParams struct {
	Name     string
	Arch     string
	Commands []string
	Logger   *log.Logger
}

/* returns binaryPath, exeName, cleanup fn (should be defered till after bin is served), and error */
func GenerateGoExe(params CreateExeParams) (binPath string, exeName string, cleanup func(), err error) {
	//creates throwaway logger if nil
	if params.Logger == nil {
		params.Logger = log.New(io.Discard, "", 0)
	}

	params.Logger.Println(`Creating temp dir`)

	tempDir, f, err := createTempDirAndFile()
	if err != nil {
		os.RemoveAll(tempDir)
		return "", "", nil, err
	}

	params.Logger.Println(`Generating script`)
	err = internal.GenerateScript(f, params.Commands...)
	if err != nil {
		return "", "", nil, err
	}

	params.Logger.Println(`Building binary`)
	binPath, exeName, err = execCmdOnTempProject(tempDir, params)
	if err != nil {
		return "", "", nil, err
	}

	params.Logger.Printf("Build successful. Binary at:%v\n", binPath)

	cleanup = func() {
		os.RemoveAll(tempDir)
	}

	return binPath, exeName, cleanup, nil
}

func createTempDirAndFile() (string, *os.File, error) {
	tempDir, err := internal.CreateTempDir("", "go-app")
	if err != nil {
		return "", nil, err
	}

	f, err := internal.CreateTempFile(tempDir, "foo")
	if err != nil {
		return tempDir, nil, err
	}
	return tempDir, f, nil
}

func execCmdOnTempProject(tempDir string, params CreateExeParams) (string, string, error) {
	// BUILD go.mod, must be done before the exe bin is created or it will fail
	params.Logger.Println(`Mod initialization`)
	modInit := createExecCmd(tempDir, "go", "mod", "init", "tmp.com/tmp")
	if out, err := modInit.CombinedOutput(); err != nil {
		params.Logger.Println("go mod initialization failed:", string(out))
		return "", "", err
	}

	arc := SupportedArchitecture[params.Arch]

	exeName := ""
	if params.Name != "" {
		exeName = params.Name
	} else {
		exeName = "app"
	}

	if params.Arch == internal.Win64 {
		exeName += ".exe"
	}

	params.Logger.Println(`Bin initialization`)

	binPath := filepath.Join(tempDir, exeName)
	binCmd := createExecCmd(tempDir, "go", "build", "-o", binPath)
	binCmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%v", arc.Os),
		fmt.Sprintf("GOARCH=%v", arc.Arch),
	)

	params.Logger.Println("Starting go build...")

	output, err := binCmd.CombinedOutput()
	if err != nil {
		fmt.Println("Build failed:", string(output))
		return "", "", err
	}

	params.Logger.Println("Go build finished.")

	return binPath, exeName, nil
}

func createExecCmd(dir string, name string, args ...string) *exec.Cmd {
	command := exec.Command(name, args...)
	command.Dir = dir
	return command
}
