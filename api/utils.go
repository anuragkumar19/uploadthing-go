package api

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/anuragkumar19/uploadthing-go/config"
	"github.com/goccy/go-json"
)

func (ut *UploadthingApi) getPresignedUrls(files []UploadFileMeta) ([]UploadthingFile, *UploadthingError) {
	body, err := json.Marshal(map[string]interface{}{
		"files":              files,
		"contentDisposition": config.Inline,
		"metadata":           "{}",
	})

	if err != nil {
		panic(fmt.Errorf("uploadthing: failed to marshal %v to json: %v", files, err.Error()))
	}

	client := http.Client{}
	req, err := http.NewRequest("POST", ut.getRequestUrl("/api/uploadFiles"), bytes.NewBuffer(body))

	if err != nil {
		panic(fmt.Errorf("uploadthing: failed to create request: %v", err.Error()))
	}

	req.Header = ut.getDefaultHeaders()

	resp, err := client.Do(req)

	if err != nil {
		panic(fmt.Errorf("uploadthing: failed to make request: %v", err.Error()))
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		panic(fmt.Errorf("uploadthing: failed to read server response: %v", err.Error()))
	}

	if resp.StatusCode != http.StatusOK {
		errResp := &UploadthingErrorResponse{}

		if err := json.Unmarshal(respBody, errResp); err != nil {
			panic(fmt.Errorf("uploadthing: failed to unmarshal server response: %v", err.Error()))
		}

		fmt.Println(resp.StatusCode)

		//TODO: handle error code better
		return nil, &UploadthingError{
			Code:    resp.Status,
			Message: "Failed to generate presigned URL",
			Data: struct{ Message string }{
				Message: errResp.Error,
			},
		}
	}

	data := &UploadthingPresignedURLResponse{}

	if err := json.Unmarshal(respBody, data); err != nil {
		panic(fmt.Errorf("uploadthing: failed to unmarshal server response: %v", err.Error()))
	}

	return data.Data, nil
}

type UploadPartParams struct {
	Url                    string
	Key                    string
	Chunk                  []byte
	ContentType            string
	ContentDispositionType config.ContentDisposition
	FileName               string
	MaxRetries             int
	PartNumber             int
}

func (ut *UploadthingApi) uploadPart(wg *sync.WaitGroup, ch chan EtagReturn, params *UploadPartParams, retryCount int) {
	defer wg.Done()

	client := http.Client{}
	req, err := http.NewRequest("PUT", params.Url, bytes.NewBuffer(params.Chunk))

	if err != nil {
		panic(fmt.Errorf("uploadthing: failed to create request: %v", err.Error()))
	}

	req.Header = http.Header{
		"Content-Type": {params.ContentType},
		// TODO: try without join
		"Content-Disposition": {strings.Join([]string{
			string(params.ContentDispositionType),
			fmt.Sprintf("filename=\"%v\"", url.QueryEscape(params.FileName)),
			fmt.Sprintf("filename*=UTF-8''%v", url.QueryEscape(params.FileName)),
		}, "; ")},
	}

	resp, err := client.Do(req)

	if err != nil {
		panic(fmt.Errorf("uploadthing: failed to make request - params = %v: %v", params, err.Error()))
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		etag := resp.Header.Get("Etag")

		if etag == "" {
			ch <- EtagReturn{
				Error: &UploadthingError{
					//TODO:
					Code:    "",
					Message: "Upload failed. Etag not found",
					Data: struct{ Message string }{
						Message: "Etag not found",
					},
				},
				PartNumber: params.PartNumber,
			}
			return
		}

		ch <- EtagReturn{
			Etag:       strings.ReplaceAll(etag, "\"", ""),
			PartNumber: params.PartNumber,
		}
		return
	}

	if retryCount < params.MaxRetries {
		delay := int(math.Pow(2, float64(retryCount)) * 1000)
		time.Sleep(time.Duration(delay))
		ut.uploadPart(wg, ch, params, retryCount+1)
	}

	ut.failureCallback(params.Key)

	respBody, err := io.ReadAll(req.Body)

	if err != nil {
		panic(fmt.Errorf("uploadthing: failed to read response: %v", err.Error()))
	}

	ch <- EtagReturn{
		Error: &UploadthingError{
			Code:    "UPLOAD_FAILED",
			Message: "Failed to upload file to storage provider",
			Data: struct{ Message string }{
				Message: string(respBody),
			},
		},
		PartNumber: params.PartNumber,
	}
}

func (ut *UploadthingApi) failureCallback(key string) {
	b := map[string]string{
		"fileKey": key,
	}

	body, err := json.Marshal(b)

	if err != nil {
		panic(fmt.Errorf("uploadthing: failed to marshal %v to json: %v", b, err.Error()))
	}

	client := http.Client{}
	req, err := http.NewRequest("POST", ut.getRequestUrl("/api/failureCallback"), bytes.NewBuffer(body))

	if err != nil {
		panic(fmt.Errorf("uploadthing: failed to create request: %v", err.Error()))
	}

	req.Header = ut.getDefaultHeaders()

	_, err = client.Do(req)

	if err != nil {
		panic(fmt.Errorf("uploadthing: failed to make request: %v", err.Error()))
	}
}

func (ut *UploadthingApi) completeMultipart(key string, uploadId string, etags []string) {
	b := map[string]interface{}{
		"fileKey":  key,
		"uploadId": uploadId,
		"etags":    etags,
	}

	body, err := json.Marshal(b)

	if err != nil {
		panic(fmt.Errorf("uploadthing: failed to marshal %v to json: %v", b, err.Error()))
	}

	client := http.Client{}
	req, err := http.NewRequest("POST", ut.getRequestUrl("/api/completeMultipart"), bytes.NewBuffer(body))

	if err != nil {
		panic(fmt.Errorf("uploadthing: failed to create request: %v", err.Error()))
	}

	req.Header = ut.getDefaultHeaders()

	_, err = client.Do(req)

	if err != nil {
		panic(fmt.Errorf("uploadthing: failed to make request: %v", err.Error()))
	}
}

func (ut *UploadthingApi) confirmUpload(key string, maxRetries int, retryCount int) bool {
	client := http.Client{}
	req, err := http.NewRequest("GET", ut.getRequestUrl("/api/pollUpload/"+key), nil)

	if err != nil {
		panic(fmt.Errorf("uploadthing: failed to create request: %v", err.Error()))
	}

	h := ut.getDefaultHeaders().Clone()
	h.Del("Content-Type")
	req.Header = h

	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		if retryCount < maxRetries {
			delay := int(math.Pow(2, float64(retryCount)) * 1000)
			time.Sleep(time.Duration(delay))
			return ut.confirmUpload(key, maxRetries, retryCount+1)
		} else {
			return false
		}
	}

	return true
}
