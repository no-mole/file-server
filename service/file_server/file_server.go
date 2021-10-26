package file_server

import (
	"bufio"
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"file-server/model"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"

	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/redis"

	"google.golang.org/grpc"

	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/config"

	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/registry"
	"smart.gitlab.biomind.com.cn/intelligent-system/biogo/utils"
	pb "smart.gitlab.biomind.com.cn/intelligent-system/protos/file_server"
)

var ConnMap map[string]*grpc.ClientConn

type Service struct {
	*registry.Metadata
	pb.UnimplementedFileServerServiceServer
}

func New() *Service {
	return &Service{
		Metadata: pb.Metadata(),
	}
}

func init() {
	ConnMap = make(map[string]*grpc.ClientConn)
}

func (s *Service) SingleUpload(ctx context.Context, in *pb.UploadInfo) (ret *pb.UpLoadResponse, err error) {
	err = s.CheckBucket(in.Bucket)
	if err != nil {
		return
	}
	filePath := path.Join(utils.GetCurrentAbPath(), model.RootDir, in.Bucket, in.FileName)
	file, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer file.Close()
	_, err = file.Write(in.Chunk.Content)
	if err != nil {
		return
	}

	w := md5.New()
	_, err = io.WriteString(w, string(in.Chunk.Content))
	if err != nil {
		return nil, err
	}
	etag := fmt.Sprintf("%x", w.Sum(nil))

	err = s.SaveFileMetadata(ctx, etag, file, in)
	if err != nil {
		return
	}

	return &pb.UpLoadResponse{
		Message: "success",
		Code:    0,
	}, nil
}

func (s *Service) CheckBucket(bucketName string) error {
	if bucketName == "" {
		return errors.New("bucket not setting")
	}
	path := path.Join(utils.GetCurrentAbPath(), model.RootDir, bucketName)

	if exists(path) {
		return nil
	}
	return os.MkdirAll(path, os.ModePerm)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func (s *Service) ChunkUpload(stream pb.FileServerService_ChunkUploadServer) error {
	upload := &pb.UploadInfo{
		Chunk: &pb.Chunk{},
	}
	for {
		fileChunk, err := stream.Recv()
		if err == io.EOF {
			ret, err := s.SingleUpload(context.Background(), upload)
			if err != nil {
				return err
			}
			return stream.SendAndClose(ret)
		}
		upload.Bucket = fileChunk.Bucket
		upload.FileName = fileChunk.FileName
		upload.Header = fileChunk.Header
		upload.AccessKey = fileChunk.AccessKey
		upload.Chunk.Content = append(upload.Chunk.Content, fileChunk.Chunk.Content...)
	}
}
func (s *Service) Download(ctx context.Context, in *pb.DownloadInfo) (resp *pb.DownloadResponse, err error) {
	if in.Exist {
		return s.DownloadNodeSelf(ctx, in)
	}
	storageNode, err := s.GetFileLocation(ctx, in.FileName, in.Bucket)
	if err != nil {
		return
	}
	if storageNode.NodeServerName != config.GlobalConfig.Name {
		return &pb.DownloadResponse{
			Chunk:    nil,
			Exist:    false,
			NodeName: storageNode.NodeServerName,
		}, nil
	}
	return s.DownloadNodeSelf(ctx, in)
}

func (s *Service) DownloadNodeSelf(_ context.Context, in *pb.DownloadInfo) (resp *pb.DownloadResponse, err error) {
	filePath := path.Join(utils.GetCurrentAbPath(), model.RootDir, in.Bucket, in.FileName)
	body, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	return &pb.DownloadResponse{
		Chunk:    &pb.Chunk{Content: body},
		Exist:    true,
		NodeName: "",
	}, nil
}

func (s *Service) BigFileDownload(in *pb.DownloadInfo, stream pb.FileServerService_BigFileDownloadServer) error {
	if in.Exist {
		return s.BigFileDownloadMyself(in, stream)
	}

	storageNode, err := s.GetFileLocation(context.Background(), in.FileName, in.Bucket)
	if err != nil {
		return err
	}

	if storageNode.NodeServerName != config.GlobalConfig.Name {
		stream.Send(&pb.DownloadResponse{
			Chunk:    nil,
			Exist:    false,
			NodeName: storageNode.NodeServerName,
		})
		return errors.New("file exist other node")
	}
	return s.BigFileDownloadMyself(in, stream)
}

func (s *Service) BigFileDownloadMyself(in *pb.DownloadInfo, stream pb.FileServerService_BigFileDownloadServer) error {
	filePath := path.Join(utils.GetCurrentAbPath(), model.RootDir, in.Bucket, in.FileName)
	fd, err := os.Open(filePath)
	if err != nil {
		return err
	}
	r := bufio.NewReader(fd)
	buffer := make([]byte, in.Size)
	for {
		n, err := r.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		info := &pb.DownloadResponse{
			Chunk:    &pb.Chunk{Content: buffer[:n]},
			Exist:    true,
			NodeName: "",
		}
		if err := stream.Send(info); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) SaveFileMetadata(ctx context.Context, eTag string, f *os.File, in *pb.UploadInfo) (err error) {
	fileInfo, err := f.Stat()
	if err != nil {
		return
	}
	metadata := &pb.UploadFileInfo{
		FileSize:       fileInfo.Size(),
		FileName:       fileInfo.Name(),
		BucketName:     in.Bucket,
		NodeServerName: config.GlobalConfig.Name,
		AccessKey:      in.AccessKey,
		Header:         in.Header,
		FileExtension:  path.Ext(fileInfo.Name()),
		ETage:          eTag,
		CreateTime:     time.Now().Format("2006-01-02 15:04:05"),
	}

	client, exist := redis.Client.GetClient(model.RedisEngine)
	if !exist {
		return
	}

	body, err := json.Marshal(metadata)
	if err != nil {
		return
	}
	key := fmt.Sprintf("%s/%s", in.Bucket, in.FileName)
	return client.Set(ctx, key, string(body), 0).Err()
}

func (s *Service) GetFileLocation(ctx context.Context, fileName, bucketName string) (*pb.UploadFileInfo, error) {
	client, exist := redis.Client.GetClient(model.RedisEngine)
	if !exist {
		return nil, errors.New("not match redis")
	}

	key := fmt.Sprintf("%s/%s", bucketName, fileName)
	body, err := client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	node := new(pb.UploadFileInfo)
	err = json.Unmarshal(body, &node)
	if err != nil {
		return nil, err
	}

	return node, nil
}
