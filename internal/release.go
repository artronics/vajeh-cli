package internal

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/valyala/fasttemplate"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"time"
)

const releaseTemplate = ` #  --- DO NOT EDIT --- Auto-generated at: {{ time }}
version: {{ version }}
version_prefix: {{ version_prefix }}
`

type ReleaseData struct {
	Prefix  string `yaml:"version_prefix"`
	Version string `yaml:"version"`
}

func NewReleaseData(version, prefix string) (ReleaseData, error) {
	ver, err := semver.NewVersion(version)
	if err != nil {
		return ReleaseData{}, err
	}

	return ReleaseData{Version: ver.String(), Prefix: prefix}, nil
}

func (r ReleaseData) Write(path string) error {
	content, err := renderTemplate(r)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	defer f.Close()
	_, err = fmt.Fprintf(f, "%s", content)
	if err != nil {
		return err
	}

	return nil
}

func renderTemplate(rd ReleaseData) (string, error) {
	template := fasttemplate.New(releaseTemplate, "{{ ", " }}")
	data := map[string]interface{}{
		"time":           time.Now().UTC().Format("2006-01-02 15:04:05"),
		"version":        rd.Version,
		"version_prefix": rd.Prefix}

	return template.ExecuteString(data), nil
}

func ParseReleaseFile(path string) (ReleaseData, error) {
	var data ReleaseData

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
