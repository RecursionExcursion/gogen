package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/RecursionExcursion/wsd-core/internal"
)

type CreateExeParams struct {
	Name     string
	Arch     string
	Commands []string
}

// TODO update logs to use lib.Log()
func CreateGoExe(params CreateExeParams) (string, string, error) {

	log.Println(`Creating temp dir`)

	tempDir, f, err := createTempDirAndFile()
	if err != nil {
		os.RemoveAll(tempDir)
		return "", "", err
	}
	// defer os.RemoveAll(tempDir)

	log.Println(`Generating script`)
	err = internal.GenerateScript(f, params.Commands...)
	if err != nil {
		return "", "", err
	}

	log.Println(`Building binary`)
	binPath, exeName, err := execCmdOnTempProject(tempDir, params)
	if err != nil {
		return "", "", err
	}

	log.Println(fmt.Sprintf("Build successful. Binary at:%v", binPath), 5)

	// log.Println(`Reading bin`)
	// inMemBin, err := os.ReadFile(binPath)
	// if err != nil {
	// 	return nil, "", err
	// }

	return binPath, exeName, nil
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
	log.Println(`Mod initialization`)
	modInit := createExecCmd(tempDir, "go", "mod", "init", "tmp.com/tmp")
	if out, err := modInit.CombinedOutput(); err != nil {
		fmt.Println("go mod initialization failed:", string(out))
		return "", "", err
	}

	arc := internal.SupportedArchitecture[params.Arch]

	exeName := ""
	if params.Name != "" {
		exeName = params.Name
	} else {
		exeName = "app"
	}

	if params.Arch == internal.Win64 {
		exeName += ".exe"
	}

	log.Println(`Bin initialization`)

	binPath := filepath.Join(tempDir, exeName)
	binCmd := createExecCmd(tempDir, "go", "build", "-o", binPath)
	binCmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%v", arc.Os),
		fmt.Sprintf("GOARCH=%v", arc.Arch),
	)

	log.Println("Starting go build...")

	output, err := binCmd.CombinedOutput()
	if err != nil {
		fmt.Println("Build failed:", string(output))
		return "", "", err
	}

	log.Println("Go build finished.")

	return binPath, exeName, nil
}

func createExecCmd(dir string, name string, args ...string) *exec.Cmd {
	command := exec.Command(name, args...)
	command.Dir = dir
	return command
}
