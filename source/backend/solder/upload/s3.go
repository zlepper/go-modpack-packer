package upload

import (
	//"github.com/aws/aws-sdk-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/getsentry/raven-go"
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	"github.com/zlepper/go-websocket-connection"
	"log"
	"os"
	"path"
	"strings"
)

func getConnection(modpack *types.Modpack) *s3.S3 {
	a := modpack.Technic.Upload.AWS
	svc := s3.New(session.New(&aws.Config{
		Region:      aws.String(a.Region),
		Credentials: credentials.NewStaticCredentials(a.AccessKey, a.SecretKey, ""),
	}))
	return svc
}

func GetAwsBuckets(conn websocket.WebsocketConnection, d interface{}) {
	modpack := types.CreateSingleModpackData(d)

	svc := getConnection(&modpack)

	result, err := svc.ListBuckets(&s3.ListBucketsInput{})

	if err != nil {
		conn.Error("TECHNIC.UPLOAD.AWS.UNABLE_TO_LIST_BUCKETS")
		return
	}

	buckets := make([]string, 0)
	for _, b := range result.Buckets {
		buckets = append(buckets, *b.Name)
	}
	conn.Write("found-aws-buckets", buckets)
}

func UploadFilesToS3(modpack *types.Modpack, infos []*types.OutputInfo, conn websocket.WebsocketConnection) {
	conn.Write("started-uploading-all", "")
	svc := getConnection(modpack)
	uploader := s3manager.NewUploaderWithClient(svc)
	outDir := path.Join(modpack.OutputDirectory, modpack.Name)
	log.Println(modpack.Technic.Upload.AWS.Bucket)
	bucket := aws.String(modpack.Technic.Upload.AWS.Bucket)
	for _, info := range infos {
		// If the filename is empty it indicates that the file was not actually repacked, but
		// is assumed to already be on solder
		if info.File == "" {
			continue
		}
		conn.Write("starting-upload", info.ProgressKey)
		key := strings.Replace(info.File, outDir, "", -1)
		log.Println("Key: " + key)
		keyString := aws.String(key)

		file, err := os.Open(info.File)
		if err != nil {
			raven.CaptureError(err, nil)
			log.Panic(err)
		}
		_, err = uploader.Upload(&s3manager.UploadInput{
			Body:   file,
			Bucket: bucket,
			Key:    keyString,
		})
		if err != nil {
			raven.CaptureError(err, nil)
			log.Panic(err)
		}
		conn.Write("finished-uploading", info.ProgressKey)
	}
	conn.Write("finished-all-uploading", "")
}
