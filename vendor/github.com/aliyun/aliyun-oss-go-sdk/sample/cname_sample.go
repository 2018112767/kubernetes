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
	"io/ioutil"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// CnameSample shows the cname usage
func CnameSample() {
	// New client
	client, err := oss.New(endpoint4Cname, accessID4Cname, accessKey4Cname,
		oss.UseCname(true))
	if err != nil {
		HandleError(err)
	}

	// Create bucket
	err = client.CreateBucket(bucketName4Cname)
	if err != nil {
		HandleError(err)
	}

	// Set bucket ACL
	err = client.SetBucketACL(bucketName4Cname, oss.ACLPrivate)
	if err != nil {
		HandleError(err)
	}

	// Look up bucket ACL
	gbar, err := client.GetBucketACL(bucketName4Cname)
	if err != nil {
		HandleError(err)
	}
	fmt.Println("Bucket ACL:", gbar.ACL)

	// List buckets, the list operation could not be done by cname's endpoint
	_, err = client.ListBuckets()
	if err == nil {
		HandleError(err)
	}

	bucket, err := client.Bucket(bucketName4Cname)
	if err != nil {
		HandleError(err)
	}

	objectValue := "长忆观潮, 满郭人争江上望。来疑沧海尽成空, 万面鼓声中。弄潮儿向涛头立, 手把红旗旗不湿。别来几向梦中看, 梦觉尚心寒。"

	// Put object
	err = bucket.PutObject(objectKey, strings.NewReader(objectValue))
	if err != nil {
		HandleError(err)
	}

	// Get object
	body, err := bucket.GetObject(objectKey)
	if err != nil {
		HandleError(err)
	}
	data, err := ioutil.ReadAll(body)
	body.Close()
	if err != nil {
		HandleError(err)
	}
	fmt.Println(objectKey, ":", string(data))

	// Put object from file
	err = bucket.PutObjectFromFile(objectKey, localFile)
	if err != nil {
		HandleError(err)
	}

	// Get object to file
	err = bucket.GetObjectToFile(objectKey, newPicName)
	if err != nil {
		HandleError(err)
	}

	// List objects
	lor, err := bucket.ListObjects()
	if err != nil {
		HandleError(err)
	}
	fmt.Println("objects:", lor.Objects)

	// Delete object
	err = bucket.DeleteObject(objectKey)
	if err != nil {
		HandleError(err)
	}

	fmt.Println("CnameSample completed")
}
