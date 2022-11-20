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

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// BucketACLSample shows how to get and set the bucket ACL
func BucketACLSample() {
	// New client
	client, err := oss.New(endpoint, accessID, accessKey)
	if err != nil {
		HandleError(err)
	}

	// Create a bucket with default parameters
	err = client.CreateBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	// Set bucket ACL. The valid ACLs are ACLPrivate、ACLPublicRead、ACLPublicReadWrite
	err = client.SetBucketACL(bucketName, oss.ACLPublicRead)
	if err != nil {
		HandleError(err)
	}

	// Get bucket ACL
	gbar, err := client.GetBucketACL(bucketName)
	if err != nil {
		HandleError(err)
	}
	fmt.Println("Bucket ACL:", gbar.ACL)

	// Delete the bucket
	err = client.DeleteBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	fmt.Println("BucketACLSample completed")
}
