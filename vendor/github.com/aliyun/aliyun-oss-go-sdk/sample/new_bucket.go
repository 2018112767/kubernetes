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
package sample

import (
	"fmt"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// NewBucketSample shows how to initialize client and bucket
func NewBucketSample() {
	// New client
	client, err := oss.New(endpoint, accessID, accessKey)
	if err != nil {
		HandleError(err)
	}

	// Create bucket
	err = client.CreateBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	// New bucket
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	// Put object, uploads an object
	var objectName = "myobject"
	err = bucket.PutObject(objectName, strings.NewReader("MyObjectValue"))
	if err != nil {
		HandleError(err)
	}

	// Delete object, deletes an object
	err = bucket.DeleteObject(objectName)
	if err != nil {
		HandleError(err)
	}

	// Delete bucket
	err = client.DeleteBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	fmt.Println("NewBucketSample completed")
}
