package internal

import (
	"bytes"
	"os"
	"text/template"
)

// GenerateScript - creates a go script that opens files/urls.
// f- file where the script is going to be created.
// args- urls/paths to open.
func GenerateScript(f *os.File, args ...string) error {

	var buff bytes.Buffer

	//write exeFn declaration
	buff.WriteString(execFnTemplate + "\n")

	//parse other templates
	tmp, err := template.New("main").Parse(mainFuncTemplate)
	if err != nil {
		return err
	}
	tmp.New("exeFnCall").Parse(execFnCallTemplate)

	//inject args into exeFn calls
	var exeFnCallBuff bytes.Buffer

	for _, arg := range args {
		tmp.ExecuteTemplate(&exeFnCallBuff, "exeFnCall", struct {
			Arg string
		}{
			arg,
		})
		exeFnCallBuff.WriteString("\n")
	}

	tmp.ExecuteTemplate(&buff, "main", struct {
		Arg string
	}{
		exeFnCallBuff.String(),
	})

	_, err = f.Write(buff.Bytes())
	if err != nil {
		return err
	}
	return nil
}
