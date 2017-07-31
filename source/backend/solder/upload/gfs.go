package upload

import (
	"github.com/mitchellh/mapstructure"
	"github.com/zlepper/gfs"
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	"log"
	"path"
	"strings"
)

type GfsTestResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func UploadFilesToGfs(modpack *types.Modpack, infos []*types.OutputInfo, conn types.WebsocketConnection) {
	conn.Write("started-uploading-all", "")
	var err error

	gfsDetails := modpack.Technic.Upload.GFS

	client, err := gfs.NewClient(gfsDetails.Url, gfsDetails.Username, gfsDetails.Password)
	if err != nil {
		panic(err)
	}

	outdir := path.Join(modpack.OutputDirectory, modpack.Name)
	log.Println("Uploading files to GFS")
	for _, info := range infos {
		err := func(info *types.OutputInfo, client *gfs.Client) error {
			if info.File == "" {
				conn.Write("starting-upload", info.Name)
				conn.Write("finished-uploading", info.Name)
				return nil
			}

			log.Println("Uploading file", info.File)

			remotePath := strings.Replace(info.File, outdir, "", 1)
			conn.Write("starting-upload", info.Name)

			uploadFile, err := gfs.NewUploadFileFromDisk(info.File, path.Dir(remotePath))
			if err != nil {
				return err
			}

			log.Println(uploadFile)

			defer uploadFile.Reader.Close()

			err = client.UploadFile(uploadFile)
			if err != nil {
				return err
			}
			conn.Write("finished-uploading", info.Name)
			return nil

		}(info, client)
		if err != nil {
			log.Panic(err)
		}
	}
}

func testGfs(data interface{}) error {
	var loginInfo types.GfsConfig
	err := mapstructure.Decode(data, &loginInfo)
	if err != nil {
		return err
	}

	_, err = gfs.NewClient(loginInfo.Url, loginInfo.Username, loginInfo.Password)
	return err
}

func TestGfs(conn types.WebsocketConnection, data interface{}) {
	err := testGfs(data)
	testResult := GfsTestResult{
		Success: err == nil,
	}
	if testResult.Success {
		testResult.Message = "GFS connection test was successful"
	} else {
		testResult.Message = err.Error()
	}
	conn.Write("gfs-test", testResult)
}
