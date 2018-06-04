package filelist

import (
	"context"
	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
	"gopkg.in/yaml.v2"
	"path/filepath"
)

// Client is the struct for local file list endpoint.
type Client struct {
	ListPath string `yaml:"list_path"`

	Path string

	AbsPath string
}

// New will create a new file list client.
func New(ctx context.Context, et uint8) (c *Client, err error) {
	t, err := model.GetTask(ctx)
	if err != nil {
		return
	}

	c = &Client{}

	e := t.Src
	if et == constants.DestinationEndpoint {
		e = t.Dst
	}

	content, err := yaml.Marshal(e.Options)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(content, c)
	if err != nil {
		return
	}

	// Set prefix
	c.Path = e.Path
	c.AbsPath, err = filepath.Abs(e.Path)
	if err != nil {
		return
	}

	return
}
