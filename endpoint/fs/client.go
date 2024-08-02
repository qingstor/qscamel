package fs

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"path/filepath"
	"strings"

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
	EnableLinkFollow bool   `yaml:"enable_link_follow"`
	Encoding         string `yaml:"encoding"`
}

func (o *Options) Check() error {
	switch o.Encoding {
	case "":
	case constants.GBK:
	case constants.HZGB2312:
	case constants.Big5:
	case constants.Windows1252:
	default:
		logrus.Errorf("%s is not a valid value for task encoding", o.Encoding)
		return constants.ErrTaskInvalid
	}

	return nil
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

	err = opt.Check()
	if err != nil {
		return
	}

	c.Options = opt
	return
}

func (c *Client) Encode(key string) (string, error) {
	if c.Options.Encoding != "" {
		utf8, err := encode(key, c.Options.Encoding)
		if err != nil {
			return "", err
		}
		return utf8, nil
	}

	return key, nil
}

func encode(input, encodingName string) (string, error) {
	var enc encoding.Encoding
	switch strings.ToLower(encodingName) {
	case constants.GBK:
		enc = simplifiedchinese.GBK
	case constants.HZGB2312:
		enc = simplifiedchinese.HZGB2312
	case constants.Big5:
		enc = traditionalchinese.Big5
	case constants.Windows1252:
		enc = charmap.Windows1252
	default:
		return "", fmt.Errorf("unsupported encoding: %s", encodingName)
	}

	reader := transform.NewReader(strings.NewReader(input), enc.NewEncoder())
	resBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(resBytes), nil
}
