module github.com/yunify/qscamel

go 1.13

require (
	cloud.google.com/go v0.53.0 // indirect
	cloud.google.com/go/storage v1.6.0
	github.com/aliyun/aliyun-oss-go-sdk v0.0.0-20180615125516-36bf7aa2f916
	github.com/aws/aws-sdk-go v1.30.7
	github.com/baiyubin/aliyun-sts-go-sdk v0.0.0-20180326062324-cfa1a18b161f // indirect
	github.com/cenkalti/backoff v1.1.0
	github.com/colinmarc/hdfs/v2 v2.1.1
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/pengsrc/go-shared v0.2.1-0.20190131101655-1999055a4a14
	github.com/qiniu/api.v7 v0.0.0-20190307065957-039fdba59f73
	github.com/qiniu/x v7.0.8+incompatible
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/sirupsen/logrus v1.5.0
	github.com/spf13/cobra v0.0.7
	github.com/stretchr/testify v1.5.1
	github.com/syndtr/goleveldb v0.0.0-20180521045021-5d6fca44a948
	github.com/tencentyun/cos-go-sdk-v5 v0.0.0-20191022140025-978604eaa44d
	github.com/upyun/go-sdk v2.1.0+incompatible
	github.com/vmihailenco/msgpack v3.3.3+incompatible
	github.com/yunify/qingstor-sdk-go/v3 v3.2.0
	google.golang.org/api v0.20.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0-20170531160350-a96e63847dc3
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/qiniu/x => github.com/Xuanwo/qiniu_x v0.0.0-20190416044656-4dd63e731f37
