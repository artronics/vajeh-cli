package internal

import (
	"fmt"
	"github.com/valyala/fasttemplate"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"time"
)

const releaseTemplate = ` #  --- DO NOT EDIT --- Auto-generated at: {{ time }}
version: {{ version }}
releaseId: {{ releaseId }}
commitId: {{ commitId }}
`

type ReleaseData struct {
	Version string
}

func (r ReleaseData) Write(path string) error {
	content := renderTemplate(r)

	f, err := os.Create(path)
	defer f.Close()
	_, err = fmt.Fprintf(f, "%s", content)
	if err != nil {
		return err
	}

	return nil
}

func renderTemplate(rd ReleaseData) string {
	template := fasttemplate.New(releaseTemplate, "{{ ", " }}")

	data := map[string]interface{}{
		"time":    time.Now().UTC().Format("2006-01-02 15:04:05"),
		"version": rd.Version}

	return template.ExecuteString(data)
}

func parseReleaseFile(path string) (map[string]interface{}, error) {
	var data map[string]interface{}

	file, err := os.Open(path)
	if err != nil {
		return data, err
	}
	defer func() {
		err = file.Close()
	}()
	bb, err := ioutil.ReadAll(file)
	if err != nil {
		return data, err
	}

	err = yaml.Unmarshal(bb, &data)
	if err != nil {
		return data, err
	}

	return data, err
}
