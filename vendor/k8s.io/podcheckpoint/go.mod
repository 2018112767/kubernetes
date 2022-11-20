module k8s.io/podcheckpoint

go 1.13

require k8s.io/apimachinery v0.25.4

require (
	github.com/aliyun/aliyun-oss-go-sdk v2.2.5+incompatible
	github.com/go-criu v0.0.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.5.2
	k8s.io/api v0.25.4
	k8s.io/client-go v0.25.4
	k8s.io/klog/v2 v2.70.1
)

replace github.com/go-criu => ../../github.com/go-criu
