package internal

var mainFuncTemplate = `
func main(){
	{{.Arg}}
}`

var execFnCallTemplate = `execCommand({{.Arg}})`

var execFnTemplate = `
package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"log"
)
	
	func execCommand(arg string) error {

	cmdPrefix := "cmd:"
	urlPrefix := "url:"

	var cmd *exec.Cmd
	if strings.HasPrefix(arg, cmdPrefix) {
		parts := strings.SplitN(aeg, ":", 2)
		if len(parts) != 2 {
			log.Fatal("invalid code command")
		}
		cmdParts := strings.SplitN(parts[1], " ", -1)
		cmd = exec.Command(cmdParts[0], cmdParts[1:]...)
	} else if strings.HasPrefix(arg, urlPrefix) {
		parts := strings.SplitN(arg, ":", 2)
		if len(parts) != 2 {
			log.Fatal("invalid code command")
		}
		url := parts[1]

		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
		case "darwin":
			cmd = exec.Command("open", url)
		case "linux":
			cmd = exec.Command("xdg-open", url)
		default:
			return fmt.Errorf("unsupported platform")
		}
	} else {
		log.Fatal("Missing command prefix")
	}

	return cmd.Start()
}`
