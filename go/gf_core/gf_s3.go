/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

package gf_core

import (
	"os"
	"bytes"
	"fmt"
	"net/http"
	"mime"
	"time"
	"path/filepath"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

//---------------------------------------------------
type GFs3Info struct {
	Client   *s3.S3
	Uploader *s3manager.Uploader
	Session  *session.Session
}

//---------------------------------------------------
func S3getFile(pTargetFileS3pathStr string,
	pTargetFileLocalPathStr string,
	pS3bucketNameStr        string,
	pS3info                 *GFs3Info,
	pRuntimeSys             *RuntimeSys) *GFerror {
	
	fmt.Printf("target_file_s3_path - %s\n", pTargetFileS3pathStr)
	fmt.Printf("s3_bucket_name      - %s\n", pS3bucketNameStr)
	
	downloader := s3manager.NewDownloader(pS3info.Session)

	// create a local host FS file to store the downloaded image into
	file, err := os.Create(pTargetFileLocalPathStr)
	if err != nil {
		gfErr := ErrorCreate("failed to create local file on host FS, to save a downloaded S3 file to.",
			"file_create_error", 
			map[string]interface{}{
				"target_file__s3_path_str":    pTargetFileS3pathStr,
				"target_file__local_path_str": pTargetFileLocalPathStr,
				"s3_bucket_name_str":          pS3bucketNameStr,
			}, err, "gf_core", pRuntimeSys)
		return gfErr
	}

	// write downloaded S3 file contents to the local FS file
	bytesDownloadedInt, err := downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(pS3bucketNameStr),
		Key:    aws.String(pTargetFileS3pathStr),
	})

	if err != nil {
		gfErr := ErrorCreate("failed to download an image from S3 bucket",
			"s3_file_download_error", nil, err, "gf_core", pRuntimeSys)
		return gfErr
	}
	fmt.Printf("file downloaded, %d bytes\n", bytesDownloadedInt)


	return nil
}

//---------------------------------------------------
// S3_INIT
func S3init(p_aws_access_key_id_str string,
	p_aws_secret_access_key_str string,
	p_token_str                 string,
	pRuntimeSys                 *RuntimeSys) (*GFs3Info, *GFerror) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_s3.S3init()")

	
	config := &aws.Config{
		Region:           aws.String("us-east-1"),
		Endpoint:         aws.String("s3.amazonaws.com"),
		S3ForcePathStyle: aws.Bool(true),      // <-- without these lines. All will fail! fork you aws!
		// Credentials:      creds,
		// LogLevel:         0, // <-- feel free to crank it up 
	}

	//--------------
	// STATIC_CREDENTIALS - they're non-empty and should be constructed. otherwise AWS creds are acquired
	//                      by the AWS client from the environment.
	if p_aws_access_key_id_str != "" {

		creds  := credentials.NewStaticCredentials(p_aws_access_key_id_str, p_aws_secret_access_key_str, p_token_str)
		_, err := creds.Get()

		if err != nil {
			gf_err := ErrorCreate("failed to acquire S3 static credentials - (credentials.NewStaticCredentials().Get())",
				"s3_credentials_error", nil, err, "gf_core", pRuntimeSys)
			return nil, gf_err
		}

		config.Credentials = creds
	}

	//--------------

	sess := session.New(config)

	s3_uploader := s3manager.NewUploader(sess)
	s3_client   := s3.New(sess)

	s3_info := &GFs3Info{
		Client:   s3_client,
		Uploader: s3_uploader,
		Session:  sess,
	}

	return s3_info, nil
}

//---------------------------------------------------
// S3__GENERATE_PRESIGNED_URL
func S3generatePresignedUploadURL(pTargetFileS3pathStr string,
	pS3bucketNameStr string,
	pS3info          *GFs3Info,
	pRuntimeSys      *RuntimeSys) (string, *GFerror) {

	// INPUT
	fileEXTstr     := filepath.Ext(pTargetFileS3pathStr)
	contentTypeStr := mime.TypeByExtension(fileEXTstr)

	s3_put_object_params := &s3.PutObjectInput{
		ACL:         aws.String("public-read"),
		Bucket:      aws.String(pS3bucketNameStr),
		Key:         aws.String(pTargetFileS3pathStr),
		ContentType: aws.String(contentTypeStr),
	}

	req, _ := pS3info.Client.PutObjectRequest(s3_put_object_params)

	// PRESIGN
	presignedURLstr, err := req.Presign(time.Minute * 1)
	if err != nil { // resp is now filled
		gf_err := ErrorCreate("failed to generate pre-signed S3 putObject URL",
			"s3_file_upload_url_presign_error", nil, err, "gf_core", pRuntimeSys)
		return "", gf_err
	}

	return presignedURLstr, nil
}

//---------------------------------------------------
// S3__UPLOAD_FILE
func S3uploadFile(p_target_file__local_path_str string,
	p_target_file__s3_path_str string,
	p_s3_bucket_name_str       string,
	p_s3_info                  *GFs3Info,
	pRuntimeSys                *RuntimeSys) (string, *GFerror) {
	pRuntimeSys.LogFun("FUN_ENTER", "gf_s3.S3uploadFile()")
	pRuntimeSys.LogFun("INFO",      "p_s3_bucket_name_str       - "+p_s3_bucket_name_str)
	pRuntimeSys.LogFun("INFO",      "p_target_file__s3_path_str - "+p_target_file__s3_path_str)

	//-----------------
	file, fs_err := os.Open(p_target_file__local_path_str)
	if fs_err != nil {
		gf_err := ErrorCreate("failed to open a local file to upload it to S3",
			"file_open_error",
			map[string]interface{}{
				"bucket_name_str":             p_s3_bucket_name_str,
				"target_file__local_path_str": p_target_file__local_path_str,
				"target_file__s3_path_str":    p_target_file__s3_path_str,
			},
			fs_err, "gf_core", pRuntimeSys)
		return "", gf_err
	}
	defer file.Close()
	
	//-----------------

	file_info,_   := file.Stat()
	var size int64 = file_info.Size()

	buffer := make([]byte, size)

	// read file content to buffer
	file.Read(buffer)

	file_bytes := bytes.NewReader(buffer) // convert to io.ReadSeeker type
	file_type  := http.DetectContentType(buffer)

	// Upload uploads an object to S3, intelligently buffering large files 
	// into smaller chunks and sending them in parallel across multiple goroutines.
	result, s3_err := p_s3_info.Uploader.Upload(&s3manager.UploadInput{
		ACL:         aws.String("public-read"),
		Bucket:      aws.String(p_s3_bucket_name_str),
		Key:         aws.String(p_target_file__s3_path_str),
		ContentType: aws.String(file_type),
		Body:        file_bytes,
	})

	if s3_err != nil {
		gf_err := ErrorCreate("failed to upload a file to an S3 bucket",
			"s3_file_upload_error",
			map[string]interface{}{
				"bucket_name_str":             p_s3_bucket_name_str,
				"target_file__local_path_str": p_target_file__local_path_str,
				"target_file__s3_path_str":    p_target_file__s3_path_str,
			},
			s3_err, "gf_core", pRuntimeSys)
		return "", gf_err
	}

	rStr := fmt.Sprint(result)
	return rStr, nil
}

//---------------------------------------------------
// S3__COPY_FILE
func S3copyFile(p_source_bucket_str string,
	p_source_file__s3_path_str string,
	p_target_bucket_name_str   string,
	p_target_file__s3_path_str string,
	pS3info                    *GFs3Info,
	pRuntimeSys                *RuntimeSys) *GFerror {

	fmt.Printf("source_bucket        - %s\n", p_source_bucket_str)
	fmt.Printf("source_file__s3_path - %s\n", p_source_file__s3_path_str)
	fmt.Printf("target_bucket_name   - %s\n", p_target_bucket_name_str)
	fmt.Printf("target_file__s3_path - %s\n", p_target_file__s3_path_str)

	source_bucket_and_file__s3_path_str := filepath.Clean(fmt.Sprintf("/%s/%s", p_source_bucket_str, p_source_file__s3_path_str))

	svc   := s3.New(pS3info.Session)
	input := &s3.CopyObjectInput{
		CopySource: aws.String(source_bucket_and_file__s3_path_str),
	    Bucket:     aws.String(p_target_bucket_name_str),
	    Key:        aws.String(p_target_file__s3_path_str),
	}

	result, err := svc.CopyObject(input)
	if err != nil {
		gf_err := ErrorCreate("failed to copy a file within S3",
			"s3_file_copy_error",
			map[string]interface{}{
				"source_bucket_and_file__s3_path_str": source_bucket_and_file__s3_path_str,
				"target_bucket_name_str":              p_target_bucket_name_str,
				"target_file__s3_path_str":            p_target_file__s3_path_str,
			},
			err, "gf_core", pRuntimeSys)
		return gf_err
	}

	fmt.Println(result)

	return nil
}