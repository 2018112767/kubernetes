/*
Copyright (c) 2014-2020 CGCL Labs
Container_Migrate is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
         http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.
*/
// main of samples

package main

import (
	"fmt"
    
	"github.com/aliyun/aliyun-oss-go-sdk/sample"
)

func main() {
	sample.CreateBucketSample()
	sample.NewBucketSample()
	sample.ListBucketsSample()
	sample.BucketACLSample()
	sample.BucketLifecycleSample()
	sample.BucketRefererSample()
	sample.BucketLoggingSample()
	sample.BucketCORSSample()

	sample.ObjectACLSample()
	sample.ObjectMetaSample()
	sample.ListObjectsSample()
	sample.DeleteObjectSample()
	sample.AppendObjectSample()
	sample.CopyObjectSample()
	sample.PutObjectSample()
	sample.GetObjectSample()

	sample.CnameSample()
	sample.SignURLSample()

	sample.ArchiveSample()

	fmt.Println("All samples completed")
}
