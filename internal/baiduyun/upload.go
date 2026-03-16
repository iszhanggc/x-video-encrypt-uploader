package baiduyun

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	// 默认分块大小 4MB
	defaultChunkSize = 4 * 1024 * 1024
	// 最小分块大小 1MB
	minChunkSize = 1 * 1024 * 1024
	// 最大分块大小 4GB
	maxChunkSize = 4 * 1024 * 1024 * 1024
)

// PreuploadResponse 预上传响应
type PreuploadResponse struct {
	Errno     int    `json:"errno"`
	Path      string `json:"path"`
	UploadID  string `json:"uploadid"`
	BlockList []int  `json:"block_list"`
	RequestID int64  `json:"request_id"`
}

// UploadChunkResponse 分块上传响应
type UploadChunkResponse struct {
	Errno   int    `json:"errno"`
	Md5     string `json:"md5"`
	RequestID int64 `json:"request_id"`
}

// CreateFileResponse 创建文件响应
type CreateFileResponse struct {
	Errno     int    `json:"errno"`
	Path      string `json:"path"`
	Size      int64  `json:"size"`
	CTime     int64  `json:"ctime"`
	MTime     int64  `json:"mtime"`
	MD5       string `json:"md5"`
	RequestID int64  `json:"request_id"`
}

// UploadOption 上传选项
type UploadOption struct {
	ChunkSize int64  // 分块大小，默认4MB
	Overwrite bool   // 是否覆盖已有文件
	UploadID  string // 续传的UploadID，断点续传时使用
}

// DefaultUploadOption 默认上传选项
var DefaultUploadOption = UploadOption{
	ChunkSize: defaultChunkSize,
	Overwrite: false,
}

// UploadFile 上传文件到百度云盘
// localPath: 本地文件路径
// remotePath: 云盘路径，如 /apps/备份/xxx.mp4
// option: 上传选项
func (c *Client) UploadFile(localPath, remotePath string, option ...UploadOption) (*CreateFileResponse, error) {
	opt := DefaultUploadOption
	if len(option) > 0 {
		opt = option[0]
	}

	// 校验分块大小
	if opt.ChunkSize < minChunkSize {
		opt.ChunkSize = minChunkSize
	}
	if opt.ChunkSize > maxChunkSize {
		opt.ChunkSize = maxChunkSize
	}

	// 打开本地文件
	file, err := os.Open(localPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()

	// 计算分块数量
	chunkCount := int(math.Ceil(float64(fileSize) / float64(opt.ChunkSize)))
	if chunkCount == 0 {
		chunkCount = 1
	}

	// 计算每个分块的MD5
	blockMd5s := make([]string, chunkCount)
	for i := 0; i < chunkCount; i++ {
		offset := int64(i) * opt.ChunkSize
		size := opt.ChunkSize
		if i == chunkCount-1 {
			size = fileSize - offset
		}

		buf := make([]byte, size)
		_, err := file.ReadAt(buf, offset)
		if err != nil && err != io.EOF {
			return nil, err
		}

		hash := md5.Sum(buf)
		blockMd5s[i] = hex.EncodeToString(hash[:])
	}

	// 1. 预上传
	preuploadResp, err := c.Preupload(remotePath, fileSize, blockMd5s, opt.Overwrite)
	if err != nil {
		return nil, err
	}

	// 如果是秒传，直接返回
	if preuploadResp.Errno == 0 && len(preuploadResp.BlockList) == 0 {
		return &CreateFileResponse{
			Errno: 0,
			Path:  remotePath,
			Size:  fileSize,
		}, nil
	}

	uploadID := preuploadResp.UploadID
	if opt.UploadID != "" {
		uploadID = opt.UploadID
	}

	// 2. 分块上传需要上传的块
	blockList := preuploadResp.BlockList
	if len(blockList) == 0 {
		// 如果没有返回需要上传的块，就全部上传
		for i := 0; i < chunkCount; i++ {
			blockList = append(blockList, i)
		}
	}

	for _, part := range blockList {
		if part >= chunkCount {
			continue
		}

		offset := int64(part) * opt.ChunkSize
		size := opt.ChunkSize
		if part == chunkCount-1 {
			size = fileSize - offset
		}

		buf := make([]byte, size)
		_, err := file.ReadAt(buf, offset)
		if err != nil && err != io.EOF {
			return nil, err
		}

		// 上传分块
		_, err = c.UploadChunk(remotePath, uploadID, part, buf)
		if err != nil {
			return nil, fmt.Errorf("上传分块%d失败: %v", part, err)
		}

		fmt.Printf("✅ 已上传分块 %d/%d\n", part+1, chunkCount)
	}

	// 3. 创建文件
	createResp, err := c.CreateFile(remotePath, fileSize, uploadID, blockMd5s, opt.Overwrite)
	if err != nil {
		return nil, err
	}

	return createResp, nil
}

// Preupload 预上传接口
func (c *Client) Preupload(remotePath string, fileSize int64, blockMd5s []string, overwrite bool) (*PreuploadResponse, error) {
	params := url.Values{}
	params.Set("method", "precreate")
	params.Set("access_token", c.accessToken)
	params.Set("path", remotePath)
	params.Set("size", strconv.FormatInt(fileSize, 10))
	params.Set("isdir", "0")
	params.Set("autoinit", "1")
	params.Set("block_list", toJSON(blockMd5s))
	if overwrite {
		params.Set("rtype", "1") // 覆盖
	} else {
		params.Set("rtype", "2") // 不覆盖，返回错误
	}

	reqURL := apiURL + "?" + params.Encode()
	resp, err := c.httpClient.Post(reqURL, "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var preuploadResp PreuploadResponse
	if err := json.Unmarshal(body, &preuploadResp); err != nil {
		return nil, err
	}

	if preuploadResp.Errno != 0 {
		return nil, fmt.Errorf("预上传失败，错误码: %d", preuploadResp.Errno)
	}

	return &preuploadResp, nil
}

// UploadChunk 上传分块
func (c *Client) UploadChunk(remotePath, uploadID string, partseq int, data []byte) (*UploadChunkResponse, error) {
	params := url.Values{}
	params.Set("method", "upload")
	params.Set("access_token", c.accessToken)
	params.Set("type", "tmpfile")
	params.Set("path", remotePath)
	params.Set("uploadid", uploadID)
	params.Set("partseq", strconv.Itoa(partseq))

	reqURL := "https://d.pcs.baidu.com/rest/2.0/pcs/superfile2?" + params.Encode()

	body := &bytes.Buffer{}
	boundary := "----WebKitFormBoundary7MA4YWxkTrZu0gW"
	body.WriteString("--" + boundary + "\r\n")
	body.WriteString("Content-Disposition: form-data; name=\"file\"; filename=\"blob\"\r\n")
	body.WriteString("Content-Type: application/octet-stream\r\n\r\n")
	body.Write(data)
	body.WriteString("\r\n--" + boundary + "--\r\n")

	req, err := http.NewRequest("POST", reqURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var chunkResp UploadChunkResponse
	if err := json.Unmarshal(respBody, &chunkResp); err != nil {
		return nil, err
	}

	if chunkResp.Errno != 0 {
		return nil, fmt.Errorf("上传分块失败，错误码: %d", chunkResp.Errno)
	}

	return &chunkResp, nil
}

// CreateFile 创建文件
func (c *Client) CreateFile(remotePath string, fileSize int64, uploadID string, blockMd5s []string, overwrite bool) (*CreateFileResponse, error) {
	params := url.Values{}
	params.Set("method", "create")
	params.Set("access_token", c.accessToken)
	params.Set("path", remotePath)
	params.Set("size", strconv.FormatInt(fileSize, 10))
	params.Set("isdir", "0")
	params.Set("uploadid", uploadID)
	params.Set("block_list", toJSON(blockMd5s))
	if overwrite {
		params.Set("rtype", "1")
	} else {
		params.Set("rtype", "2")
	}

	reqURL := apiURL + "?" + params.Encode()
	resp, err := c.httpClient.Post(reqURL, "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var createResp CreateFileResponse
	if err := json.Unmarshal(body, &createResp); err != nil {
		return nil, err
	}

	if createResp.Errno != 0 {
		return nil, fmt.Errorf("创建文件失败，错误码: %d", createResp.Errno)
	}

	return &createResp, nil
}

// toJSON 转换为JSON字符串
func toJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}
