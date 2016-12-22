// +-------------------------------------------------------------------------
// | Copyright (C) 2016 Yunify, Inc.
// +-------------------------------------------------------------------------
// | Licensed under the Apache License, Version 2.0 (the "License");
// | you may not use this work except in compliance with the License.
// | You may obtain a copy of the License in the LICENSE file, or at:
// |
// | http://www.apache.org/licenses/LICENSE-2.0
// |
// | Unless required by applicable law or agreed to in writing, software
// | distributed under the License is distributed on an "AS IS" BASIS,
// | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// | See the License for the specific language governing permissions and
// | limitations under the License.
// +-------------------------------------------------------------------------

package builder

import (
	"errors"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/yunify/qingstor-sdk-go/request/data"
	"github.com/yunify/qingstor-sdk-go/utils"
)

// BaseBuilder is the base builder for all services.
type BaseBuilder struct {
	parsedURL        string
	parsedProperties *map[string]string
	parsedParams     *map[string]string
	parsedHeaders    *map[string]string
	parsedBodyString string
	parsedBody       io.Reader

	operation *data.Operation
	input     *reflect.Value
}

// BuildHTTPRequest builds http request with an operation and an input.
func (b *BaseBuilder) BuildHTTPRequest(o *data.Operation, i *reflect.Value) (*http.Request, error) {
	b.operation = o
	b.input = i

	_, err := b.parse()
	if err != nil {
		return nil, err
	}

	return b.build()
}

func (b *BaseBuilder) build() (*http.Request, error) {
	httpRequest, err := http.NewRequest(b.operation.RequestMethod, b.parsedURL, b.parsedBody)
	if err != nil {
		return nil, err
	}

	err = b.setupHeaders(httpRequest)
	if err != nil {
		return nil, err
	}

	return httpRequest, nil
}

func (b *BaseBuilder) parse() (*BaseBuilder, error) {
	err := b.parseRequestParamsAndHeaders()
	if err != nil {
		return b, err
	}
	err = b.parseRequestBody()
	if err != nil {
		return b, err
	}
	err = b.parseRequestProperties()
	if err != nil {
		return b, err
	}
	err = b.parseRequestURL()
	if err != nil {
		return b, err
	}

	return b, nil
}

func (b *BaseBuilder) parseRequestParamsAndHeaders() error {
	requestParams := map[string]string{}
	requestHeaders := map[string]string{}
	maps := map[string](map[string]string){
		"params":  requestParams,
		"headers": requestHeaders,
	}

	b.parsedParams = &requestParams
	b.parsedHeaders = &requestHeaders

	if !b.input.IsValid() {
		return nil
	}

	fields := b.input.Elem()
	if !fields.IsValid() {
		return nil
	}

	for i := 0; i < fields.NumField(); i++ {
		tagName := fields.Type().Field(i).Tag.Get("name")
		tagLocation := fields.Type().Field(i).Tag.Get("location")
		tagDefault := fields.Type().Field(i).Tag.Get("default")
		if tagName != "" && tagLocation != "" && maps[tagLocation] != nil {
			switch value := fields.Field(i).Interface().(type) {
			case string:
				if value != "" {
					maps[tagLocation][tagName] = value
				}
			case int:
				numberString := strconv.Itoa(int(value))
				if numberString == "0" {
					numberString = ""
					if tagDefault != "" {
						numberString = tagDefault
					}
				}
				if numberString != "" {
					maps[tagLocation][tagName] = numberString
				}
			case bool:
			case time.Time:
				zero := time.Time{}
				if value != zero {
					var timeString string
					format := fields.Type().Field(i).Tag.Get("format")
					timeString = utils.TimeToString(value, format)
					if timeString != "" {
						maps[tagLocation][tagName] = timeString
					}
				}
			}
		}
	}

	return nil
}

func (b *BaseBuilder) parseRequestBody() error {
	requestData := map[string]interface{}{}

	if !b.input.IsValid() {
		return nil
	}

	fields := b.input.Elem()
	if !fields.IsValid() {
		return nil
	}

	for i := 0; i < fields.NumField(); i++ {
		location := fields.Type().Field(i).Tag.Get("location")
		if location == "elements" {
			name := fields.Type().Field(i).Tag.Get("name")
			requestData[name] = fields.Field(i).Interface()
		}
	}

	if len(requestData) != 0 {
		dataValue, err := utils.JSONEncode(requestData, true)
		if err != nil {
			return err
		}

		b.parsedBodyString = string(dataValue)
		b.parsedBody = strings.NewReader(b.parsedBodyString)
		(*b.parsedHeaders)["Content-Type"] = "application/json"
	} else {
		value := fields.FieldByName("Body")
		if value.IsValid() {
			switch value.Interface().(type) {
			case string:
				if value.String() != "" {
					b.parsedBodyString = value.String()
					b.parsedBody = strings.NewReader(value.String())
				}
			case io.Reader:
				if value.Interface().(io.Reader) != nil {
					b.parsedBody = value.Interface().(io.Reader)
				}
			}
		}
	}

	return nil
}

func (b *BaseBuilder) parseRequestProperties() error {
	propertiesMap := map[string]string{}
	b.parsedProperties = &propertiesMap

	if b.operation.Properties != nil {
		fields := reflect.ValueOf(b.operation.Properties).Elem()
		if fields.IsValid() {
			for i := 0; i < fields.NumField(); i++ {
				switch value := fields.Field(i).Interface().(type) {
				case string:
					propertiesMap[fields.Type().Field(i).Tag.Get("name")] = value
				case int:
					numberString := strconv.Itoa(int(value))
					propertiesMap[fields.Type().Field(i).Tag.Get("name")] = numberString
				}
			}
		}
	}

	return nil
}

func (b *BaseBuilder) parseRequestURL() error {
	return nil
}

func (b *BaseBuilder) setupHeaders(httpRequest *http.Request) error {
	if b.parsedHeaders != nil {
		for headerKey, headerValue := range *b.parsedHeaders {
			httpRequest.Header.Set(headerKey, headerValue)
		}
	}

	if httpRequest.Header.Get("Content-Length") == "" {
		var length int64
		switch body := b.parsedBody.(type) {
		case nil:
			length = 0
		case io.Seeker:
			//start, err := body.Seek(0, io.SeekStart)
			start, err := body.Seek(0, 0)
			if err != nil {
				return err
			}
			//end, err := body.Seek(0, io.SeekEnd)
			end, err := body.Seek(0, 2)
			if err != nil {
				return err
			}
			//body.Seek(0, io.SeekStart)
			body.Seek(0, 0)
			length = end - start
		default:
			return errors.New("Can not get Content-Length")
		}
		if length > 0 {
			httpRequest.ContentLength = length
			httpRequest.Header.Set("Content-Length", strconv.Itoa(int(length)))
		} else {
			httpRequest.Header.Set("Content-Length", "0")
		}
	}
	length, err := strconv.Atoi(httpRequest.Header.Get("Content-Length"))
	if err != nil {
		return err
	}
	httpRequest.ContentLength = int64(length)

	if httpRequest.Header.Get("Date") == "" {
		httpRequest.Header.Set("Date", utils.TimeToString(time.Now(), "RFC 822"))
	}

	return nil
}
