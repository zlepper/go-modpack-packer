package upload

import (
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	//"crypto/tls"
	"github.com/jlaffaye/ftp"
	"github.com/zlepper/go-websocket-connection"
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
		panic(err)
	}

	defer f.Quit()

	if err = f.Login(ftpDetails.Username, ftpDetails.Password); err != nil {
		panic(err)
	}
	outDir := path.Join(modpack.OutputDirectory, modpack.Name)

	for _, info := range infos {
		conn.Write("starting-upload", info.Name)
		err = f.ChangeDir("/")
		if err != nil {
			panic(err)
		}
		file, err := os.Open(info.File)
		if err != nil {
			panic(err)
		}
		key := strings.Replace(info.File, outDir, "", -1)
		parts := strings.Split(key, "/")

		for i, part := range parts {
			if i != len(parts)-1 && part != "" {
				if !doesDirectoryExist(f, part) {
					err = f.MakeDir(part)
					if err != nil {
						panic(err)
					}
				}
				err = f.ChangeDir(part)
				if err != nil {
					panic(err)
				}
			}
		}

		if err = f.Stor(key, file); err != nil {
			panic(err)
		}
		conn.Write("finished-uploading", info.Name)
	}
	conn.Write("finished-all-uploading", "")
}

func doesDirectoryExist(f *ftp.ServerConn, dir string) bool {
	entries, err := f.NameList("")
	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		if entry == dir {
			return true
		}
	}
	return false
}
