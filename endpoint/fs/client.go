package fs

import (
	"context"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/yunify/qscamel/constants"
	"github.com/yunify/qscamel/model"
)

// Client is the struct for POSIX file system endpoint.
type Client struct {
	Path    string
	AbsPath string
	Options Options
}

// Options is the struct for fs options
type Options struct {
	EnableLinkFollow bool `yaml:"enable_link_follow"`
}

// New will create a Fs.
func New(ctx context.Context, et uint8) (c *Client, err error) {
	t, err := model.GetTask(ctx)
	if err != nil {
		return
	}

	e := t.Src
	if et == constants.DestinationEndpoint {
		e = t.Dst
	}

	c = &Client{}

	// Set prefix.
	c.Path = e.Path
	c.AbsPath, err = filepath.Abs(e.Path)
	if err != nil {
		return
	}

	opt := Options{}
	content, err := yaml.Marshal(e.Options)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(content, &opt)
	if err != nil {
		return
	}

	c.Options = opt
	return
}
