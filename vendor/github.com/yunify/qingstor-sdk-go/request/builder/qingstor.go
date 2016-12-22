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
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/yunify/qingstor-sdk-go"
	"github.com/yunify/qingstor-sdk-go/logger"
	"github.com/yunify/qingstor-sdk-go/request/data"
	"github.com/yunify/qingstor-sdk-go/utils"
)

// QingStorBuilder is the request builder for QingStor service.
type QingStorBuilder struct {
	baseBuilder *BaseBuilder
}

// BuildHTTPRequest builds http request with an operation and an input.
func (qb *QingStorBuilder) BuildHTTPRequest(o *data.Operation, i *reflect.Value) (*http.Request, error) {
	qb.baseBuilder = &BaseBuilder{}
	qb.baseBuilder.operation = o
	qb.baseBuilder.input = i

	_, err := qb.baseBuilder.parse()
	if err != nil {
		return nil, err
	}
	err = qb.parseURL()
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequest(qb.baseBuilder.operation.RequestMethod,
		qb.baseBuilder.parsedURL, qb.baseBuilder.parsedBody)
	if err != nil {
		return nil, err
	}

	err = qb.baseBuilder.setupHeaders(httpRequest)
	if err != nil {
		return nil, err
	}
	err = qb.setupHeaders(httpRequest)
	if err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf(
		"Built QingStor request: [%d] %s",
		utils.StringToUnixInt(httpRequest.Header.Get("Date"), "RFC 822"),
		httpRequest.URL.String()))

	logger.Info(fmt.Sprintf(
		"QingStor request headers: [%d] %s",
		utils.StringToUnixInt(httpRequest.Header.Get("Date"), "RFC 822"),
		fmt.Sprint(httpRequest.Header)))

	if qb.baseBuilder.parsedBodyString != "" {
		logger.Info(fmt.Sprintf(
			"QingStor request body string: [%d] %s",
			utils.StringToUnixInt(httpRequest.Header.Get("Date"), "RFC 822"),
			qb.baseBuilder.parsedBodyString))
	}

	return httpRequest, nil
}

func (qb *QingStorBuilder) parseURL() error {
	config := qb.baseBuilder.operation.Config

	zone := (*qb.baseBuilder.parsedProperties)["zone"]
	port := strconv.Itoa(config.Port)
	endpoint := config.Protocol + "://" + config.Host + ":" + port
	if zone != "" {
		endpoint = config.Protocol + "://" + zone + "." + config.Host + ":" + port
	}

	requestURI := qb.baseBuilder.operation.RequestURI
	for key, value := range *qb.baseBuilder.parsedProperties {
		endpoint = strings.Replace(endpoint, "<"+key+">", value, -1)
		requestURI = strings.Replace(requestURI, "<"+key+">", value, -1)
	}
	requestURI = regexp.MustCompile(`/+`).ReplaceAllString(requestURI, "/")

	qb.baseBuilder.parsedURL = endpoint + requestURI

	if qb.baseBuilder.parsedParams != nil {
		paramsParts := []string{}
		for key, value := range *qb.baseBuilder.parsedParams {
			paramsParts = append(paramsParts, key+"="+value)

		}

		joined := strings.Join(paramsParts, "&")
		if joined != "" {
			qb.baseBuilder.parsedURL += "?" + joined
		}
	}

	return nil
}

func (qb *QingStorBuilder) setupHeaders(httpRequest *http.Request) error {
	method := httpRequest.Method
	if method == "POST" || method == "PUT" || method == "DELETE" {
		if httpRequest.Header.Get("Content-Type") == "" {
			mimeType := mime.TypeByExtension(path.Ext(httpRequest.URL.Path))
			if mimeType != "" {
				httpRequest.Header.Set("Content-Type", mimeType)
			}
		}
	}

	if httpRequest.Header.Get("User-Agent") == "" {
		version := "Go v" + strings.Replace(runtime.Version(), "go", "", -1) + ""
		system := runtime.GOOS + "_" + runtime.GOARCH + "_" + runtime.Compiler
		ua := "qingstor-sdk-go/" + sdk.Version + " (" + version + "; " + system + ")"
		httpRequest.Header.Set("User-Agent", ua)
	}

	if qb.baseBuilder.operation.APIName == "Delete Multiple Objects" {
		buffer := &bytes.Buffer{}
		buffer.ReadFrom(httpRequest.Body)
		httpRequest.Body = ioutil.NopCloser(bytes.NewReader(buffer.Bytes()))

		md5Value := md5.Sum(buffer.Bytes())
		httpRequest.Header.Set("Content-MD5", base64.StdEncoding.EncodeToString(md5Value[:]))
	}

	return nil
}
