package db

import (
	"encoding/json"
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	"log"
	"net/http"
	"strings"
	"sync"
)

type ftbPermission struct {
	ModName             string                 `json:"modName"`
	ModAuthors          string                 `json:"modAuthors"`
	LicenseLink         string                 `json:"licenseLink"`
	ModLink             string                 `json:"modLink"`
	PrivateLicenceLink  string                 `json:"privateLicenceLink"`
	PrivateStringPolicy types.PermissionPolicy `json:"privateStringPolicy"`
	PublicStringPolicy  types.PermissionPolicy `json:"publicStringPolicy"`
	Modids              string                 `json:"modids"`
	CustomData          string                 `json:"customData"`
	ShortName           string                 `json:"shortName"`
}

type PermissionData struct {
	ModName            string                 `json:"modName"`
	ModAuthors         string                 `json:"modAuthors"`
	LicenseLink        string                 `json:"licenseLink"`
	ModLink            string                 `json:"modLink"`
	PrivateLicenceLink string                 `json:"privateLicenceLink"`
	PrivatePolicy      types.PermissionPolicy `json:"privateStringPolicy"`
	PublicPolicy       types.PermissionPolicy `json:"publicStringPolicy"`
	Modids             []string               `json:"modids"`
}

type PermissionsDB struct {
	Permissions []*PermissionData
}

func init() {
	GetPermissionsDb()
}

var permissionDBInstance *PermissionsDB
var ready sync.WaitGroup

func GetPermissionsDb() *PermissionsDB {
	if permissionDBInstance != nil {
		ready.Wait()
		return permissionDBInstance
	}

	permissionDBInstance = &PermissionsDB{}

	ready.Add(1)
	go UpdatePermissionStore()
	return permissionDBInstance
}

func UpdatePermissionStore() {
	res, err := http.Get("http://legacy.feed-the-beast.com/mods/json")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	permissions := make([]ftbPermission, 0)
	err = json.NewDecoder(res.Body).Decode(&permissions)
	if err != nil {
		// We shouldn't panic here, as this is not fatal for the application to be able to work.
		// The application just won't be able to check permissions against the ftb list
		log.Println(err)
		return
	}

	for _, permission := range permissions {
		data := PermissionData{
			ModName:            permission.ModName,
			ModAuthors:         permission.ModAuthors,
			LicenseLink:        permission.LicenseLink,
			ModLink:            permission.ModLink,
			PrivateLicenceLink: permission.PrivateLicenceLink,
			PrivatePolicy:      permission.PrivateStringPolicy,
			PublicPolicy:       permission.PublicStringPolicy,
			Modids:             strings.Split(strings.ToLower(permission.Modids), " "),
		}
		data.Modids = append(data.Modids, permission.ShortName)

		permissionDBInstance.Permissions = append(permissionDBInstance.Permissions, &data)
	}
	ready.Done()
	log.Println("Done updating permission store")
}

var permissionCache = make(map[string]*PermissionData)

func (db *PermissionsDB) GetPermission(modId string) *PermissionData {
	if permission, ok := permissionCache[modId]; ok {
		return permission
	}
	for _, p := range db.Permissions {
		for _, id := range p.Modids {
			if id == modId {
				permissionCache[modId] = p
				return p
			}
		}
	}
	return nil
}

func (db *PermissionsDB) GetPermissionPolicy(modid string, isPublic bool) types.PermissionPolicy {
	if modid == "" {
		return types.Unknown
	}
	permission := db.GetPermission(modid)
	if permission == nil {
		return types.Unknown
	}
	if isPublic {
		return permission.PublicPolicy
	} else {
		return permission.PrivatePolicy
	}
}
