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

// ObjectMetaSample shows how to get and set the object metadata
func ObjectMetaSample() {
	// Create bucket
	bucket, err := GetTestBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	// Delete object
	err = bucket.PutObject(objectKey, strings.NewReader("YoursObjectValue"))
	if err != nil {
		HandleError(err)
	}

	// Case 0: Set bucket meta. one or more properties could be set
	// Note: Meta is case insensitive
	options := []oss.Option{
		oss.Expires(futureDate),
		oss.Meta("myprop", "mypropval")}
	err = bucket.SetObjectMeta(objectKey, options...)
	if err != nil {
		HandleError(err)
	}

	// Case 1: Get the object metadata. Only return basic meta information includes ETag, size and last modified.
	props, err := bucket.GetObjectMeta(objectKey)
	if err != nil {
		HandleError(err)
	}
	fmt.Println("Object Meta:", props)

	// Case 2: Get all the detailed object meta including custom meta
	props, err = bucket.GetObjectDetailedMeta(objectKey)
	if err != nil {
		HandleError(err)
	}
	fmt.Println("Expires:", props.Get("Expires"))

	// Case 3: Get the object's all metadata with contraints. When constraints are met, return the metadata.
	props, err = bucket.GetObjectDetailedMeta(objectKey, oss.IfUnmodifiedSince(futureDate))
	if err != nil {
		HandleError(err)
	}
	fmt.Println("MyProp:", props.Get("X-Oss-Meta-Myprop"))

	_, err = bucket.GetObjectDetailedMeta(objectKey, oss.IfModifiedSince(futureDate))
	if err == nil {
		HandleError(err)
	}

	goar, err := bucket.GetObjectACL(objectKey)
	if err != nil {
		HandleError(err)
	}
	fmt.Println("Object ACL:", goar.ACL)

	// Delete object and bucket
	err = DeleteTestBucketAndObject(bucketName)
	if err != nil {
		HandleError(err)
	}

	fmt.Println("ObjectMetaSample completed")
}
