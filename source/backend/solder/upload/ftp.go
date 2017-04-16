package upload

import (
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	//"crypto/tls"
	"errors"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/jlaffaye/ftp"
	"github.com/mitchellh/mapstructure"
	"github.com/zlepper/go-websocket-connection"
	"log"
	"os"
	"path"
	"strings"
)

func UploadFilesToFtp(modpack *types.Modpack, infos []*types.OutputInfo, conn websocket.WebsocketConnection) {
	conn.Write("started-uploading-all", "")
	var err error
	var f *ftp.ServerConn

	ftpDetails := modpack.Technic.Upload.FTP

	if f, err = ftp.Connect(ftpDetails.Url); err != nil {
		raven.CaptureErrorAndWait(err, nil)
		panic(err)
	}

	defer f.Quit()

	if err = f.Login(ftpDetails.Username, ftpDetails.Password); err != nil {
		raven.CaptureErrorAndWait(err, nil)
		panic(err)
	}
	outDir := path.Join(modpack.OutputDirectory, modpack.Name)

	remotePath := "~/"
	if ftpDetails.Path != "" {
		remotePath += ftpDetails.Path
	}

	for _, info := range infos {
		// If the file variable is empty, it indicates that the mod was not actually repacked
		if info.File == "" {
			conn.Write("starting-upload", info.Name)
			conn.Write("finished-uploading", info.Name)
			continue
		}
		conn.Write("starting-upload", info.Name)
		err = f.ChangeDir(remotePath)
		if err != nil {
			raven.CaptureErrorAndWait(err, nil)
			panic(err)
		}
		file, err := os.Open(info.File)
		if err != nil {
			raven.CaptureErrorAndWait(err, nil)
			panic(err)
		}
		key := strings.Replace(info.File, outDir, "", -1)
		parts := strings.Split(key, "/")

		for i, part := range parts {
			if i != len(parts)-1 && part != "" {
				if !doesDirectoryExist(f, part) {
					err = f.MakeDir(part)
					if err != nil {
						raven.CaptureErrorAndWait(err, nil)
						panic(err)
					}
				}
				err = f.ChangeDir(part)
				if err != nil {
					raven.CaptureErrorAndWait(err, nil)
					panic(err)
				}
			}
		}

		if err = f.Stor(path.Join(remotePath, key), file); err != nil {
			raven.CaptureErrorAndWait(err, nil)
			panic(err)
		}
		conn.Write("finished-uploading", info.Name)
	}
	conn.Write("finished-all-uploading", "")
}

func doesDirectoryExist(f *ftp.ServerConn, dir string) bool {
	entries, err := f.NameList("")
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		panic(err)
	}
	for _, entry := range entries {
		if entry == dir {
			return true
		}
	}
	return false
}

type FtpTestResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func testFtp(conn websocket.WebsocketConnection, data interface{}) error {
	dict := data.(map[string]interface{})
	var loginInfo types.FtpConfig
	err := mapstructure.Decode(dict, &loginInfo)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.Panic(err)
	}

	log.Println("Connecting to ftp", loginInfo.Url)

	var f *ftp.ServerConn
	fmt.Println(loginInfo.Url)
	if f, err = ftp.Connect(loginInfo.Url); err != nil {
		conn.Log(err.Error())
		return errors.New("Could not connect to the provided address")
	}

	defer f.Quit()
	log.Println("Logging in to ftp")

	if err = f.Login(loginInfo.Username, loginInfo.Password); err != nil {
		return errors.New("Could not login with the supplied credentials")
	}

	log.Println("Attempting ftp directory listing")

	remotePath := "~/"
	if loginInfo.Path != "" {
		remotePath = loginInfo.Path
	}

	_, err = f.List(remotePath)
	if err != nil {
		return errors.New("Could not list content. Might be a permission issue.")
	}

	log.Println("Ftp connection succesful")
	return nil
}

func TestFtp(conn websocket.WebsocketConnection, data interface{}) {
	err := testFtp(conn, data)
	testResult := FtpTestResult{
		Success: err == nil,
	}
	if err == nil {
		testResult.Message = "ftp connection test was successful"
	} else {
		testResult.Message = err.Error()
	}
	conn.Write("ftp-test", testResult)
}
