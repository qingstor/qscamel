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

package unpacker

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/yunify/qingstor-sdk-go/logger"
	"github.com/yunify/qingstor-sdk-go/request/data"
	"github.com/yunify/qingstor-sdk-go/utils"
)

// BaseUnpacker is the base unpacker for all services.
type BaseUnpacker struct {
	operation    *data.Operation
	httpResponse *http.Response
	output       *reflect.Value
}

// UnpackHTTPRequest unpacks http response with an operation and an output.
func (b *BaseUnpacker) UnpackHTTPRequest(o *data.Operation, r *http.Response, x *reflect.Value) error {
	b.operation = o
	b.httpResponse = r
	b.output = x

	err := b.exposeStatusCode()
	if err != nil {
		return err
	}
	err = b.parseResponseHeaders()
	if err != nil {
		return err
	}
	err = b.parseResponseBody()
	if err != nil {
		return err
	}
	err = b.parseResponseElements()
	if err != nil {
		return err
	}

	return nil
}

func (b *BaseUnpacker) exposeStatusCode() error {
	value := b.output.Elem().FieldByName("StatusCode")
	if value.IsValid() {
		switch value.Interface().(type) {
		case int:
			logger.Info(fmt.Sprintf(
				"QingStor response status code: [%d] %d",
				utils.StringToUnixInt(b.httpResponse.Header.Get("Date"), "RFC 822"),
				b.httpResponse.StatusCode))
			value.SetInt(int64(b.httpResponse.StatusCode))
		}
	}

	return nil
}

func (b *BaseUnpacker) parseResponseHeaders() error {
	logger.Info(fmt.Sprintf(
		"QingStor response headers: [%d] %s",
		utils.StringToUnixInt(b.httpResponse.Header.Get("Date"), "RFC 822"),
		fmt.Sprint(b.httpResponse.Header)))

	if b.isResponseRight() {
		fields := b.output.Elem()
		for i := 0; i < fields.NumField(); i++ {
			field := fields.Field(i)
			fieldTagName := fields.Type().Field(i).Tag.Get("name")
			fieldTagLocation := fields.Type().Field(i).Tag.Get("location")
			fieldStringValue := b.httpResponse.Header.Get(fieldTagName)

			if fieldTagName != "" && fieldTagLocation == "headers" {
				switch field.Interface().(type) {
				case string:
					field.Set(reflect.ValueOf(fieldStringValue))
				case int:
					intValue, err := strconv.Atoi(fieldStringValue)
					if err != nil {
						return err
					}
					field.Set(reflect.ValueOf(intValue))
				case bool:
				case time.Time:
					format := fields.Type().Field(i).Tag.Get("format")
					timeValue, err := utils.StringToTime(fieldStringValue, format)
					if err != nil {
						return err
					}
					field.Set(reflect.ValueOf(timeValue))
				}
			}
		}
	}

	return nil
}

func (b *BaseUnpacker) parseResponseBody() error {
	if b.isResponseRight() {
		value := b.output.Elem().FieldByName("Body")
		if value.IsValid() {
			switch value.Type().String() {
			case "string":
				buffer := &bytes.Buffer{}
				buffer.ReadFrom(b.httpResponse.Body)
				b.httpResponse.Body.Close()

				logger.Info(fmt.Sprintf(
					"QingStor response body string: [%d] %s",
					utils.StringToUnixInt(b.httpResponse.Header.Get("Date"), "RFC 822"),
					string(buffer.Bytes())))

				value.SetString(string(buffer.Bytes()))
			case "io.ReadCloser":
				value.Set(reflect.ValueOf(b.httpResponse.Body))
			}
		}
	}

	return nil
}

func (b *BaseUnpacker) parseResponseElements() error {
	if b.isResponseRight() {
		if b.httpResponse.Header.Get("Content-Type") == "application/json" {
			buffer := &bytes.Buffer{}
			buffer.ReadFrom(b.httpResponse.Body)
			b.httpResponse.Body.Close()

			logger.Info(fmt.Sprintf(
				"QingStor response body string: [%d] %s",
				utils.StringToUnixInt(b.httpResponse.Header.Get("Date"), "RFC 822"),
				string(buffer.Bytes())))

			_, err := utils.JSONDecode(buffer.Bytes(), b.output.Interface())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *BaseUnpacker) isResponseRight() bool {
	rightStatusCodes := b.operation.StatusCodes
	if len(rightStatusCodes) == 0 {
		rightStatusCodes = append(rightStatusCodes, 200)
	}

	flag := false
	for _, statusCode := range rightStatusCodes {
		if statusCode == b.httpResponse.StatusCode {
			flag = true
		}
	}

	return flag
}
