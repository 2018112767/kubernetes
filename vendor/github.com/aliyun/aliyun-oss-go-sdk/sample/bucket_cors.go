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

// BucketCORSSample shows how to get or set the bucket CORS.
func BucketCORSSample() {
	// New client
	client, err := oss.New(endpoint, accessID, accessKey)
	if err != nil {
		HandleError(err)
	}

	// Create the bucket with default parameters
	err = client.CreateBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	rule1 := oss.CORSRule{
		AllowedOrigin: []string{"*"},
		AllowedMethod: []string{"PUT", "GET", "POST"},
		AllowedHeader: []string{},
		ExposeHeader:  []string{},
		MaxAgeSeconds: 100,
	}

	rule2 := oss.CORSRule{
		AllowedOrigin: []string{"http://www.a.com", "http://www.b.com"},
		AllowedMethod: []string{"GET"},
		AllowedHeader: []string{"Authorization"},
		ExposeHeader:  []string{"x-oss-test", "x-oss-test1"},
		MaxAgeSeconds: 100,
	}

	// Case 1: Set the bucket CORS rules
	err = client.SetBucketCORS(bucketName, []oss.CORSRule{rule1})
	if err != nil {
		HandleError(err)
	}

	// Case 2: Set the bucket CORS rules. if CORS rules exist, they will be overwritten.
	err = client.SetBucketCORS(bucketName, []oss.CORSRule{rule1, rule2})
	if err != nil {
		HandleError(err)
	}

	// Get the bucket's CORS
	gbl, err := client.GetBucketCORS(bucketName)
	if err != nil {
		HandleError(err)
	}
	fmt.Println("Bucket CORS:", gbl.CORSRules)

	// Delete bucket's CORS
	err = client.DeleteBucketCORS(bucketName)
	if err != nil {
		HandleError(err)
	}

	// Delete bucket
	err = client.DeleteBucket(bucketName)
	if err != nil {
		HandleError(err)
	}

	fmt.Println("BucketCORSSample completed")
}
