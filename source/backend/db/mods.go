package db

import (
	"encoding/json"
	"fmt"
	"github.com/zlepper/go-modpack-packer/source/backend/types"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

func init() {
	GetModsDb()
}

type modsDb struct {
	Mods []*types.Mod
	m    sync.Mutex
}

var modsDbInstance *modsDb

func GetModsDb() *modsDb {
	if modsDbInstance != nil {
		return modsDbInstance
	}

	modsDbInstance = &modsDb{
		Mods: make([]*types.Mod, 0),
	}

	dataDirectory := os.Args[1]
	modsFile := filepath.Join(dataDirectory, "mods.json")
	f, err := os.Open(modsFile)
	if err != nil {
		return modsDbInstance
	}

	err = json.NewDecoder(f).Decode(&modsDbInstance.Mods)
	if err != nil {
		panic(err)
	}

	return modsDbInstance
}

func (m *modsDb) GetModFromMd5(md5 string) *types.Mod {
	//fmt.Println("Looking for mod " + md5)
	for i, _ := range m.Mods {
		if m.Mods[i].Md5 == md5 {
			//fmt.Println("Found mod with " + md5)
			return m.Mods[i]
		}
	}
	return nil
}

func (m *modsDb) GetModsWithModId(modId string) []*types.Mod {
	modId = strings.ToLower(modId)
	mods := make([]*types.Mod, 0)
	for _, mod := range m.Mods {
		if strings.ToLower(mod.ModId) == modId {
			mods = append(mods, mod)
		}
	}
	return mods
}

func (m *modsDb) Save() {
	// Don't start saving this multiple times, that will mess things up severely
	m.m.Lock()
	dataDirectory := os.Args[1]
	modsFile := filepath.Join(dataDirectory, "mods.json")
	f, err := os.Create(modsFile)
	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(f).Encode(m.Mods)
	if err != nil {
		panic(err)
	}
	fmt.Println("Wrote " + strconv.Itoa(len(m.Mods)) + " mods to the mods database")
	m.m.Unlock()
}

func (m *modsDb) AddMod(mod *types.Mod) {
	fmt.Println(mod.Md5)
	// Check if mods exists
	for i, _ := range m.Mods {
		// Check if the mods are the same
		// If the md5 matches, then it's likely the same mod
		if m.Mods[i].Md5 == mod.Md5 {
			// Update
			m.Mods[i] = mod
			return
		}
	}

	// We couldn't find the mod, so we'll just add it
	m.Mods = append(m.Mods, mod)
	fmt.Println("Adding mod: " + mod.Name)
}

func (m *modsDb) AddMods(mods []*types.Mod) {
	for _, mod := range mods {
		m.AddMod(mod)
	}
}

func (m *modsDb) MarkModAsOnSolder(md5 string) {
	for _, mod := range m.Mods {
		if mod.Md5 == md5 {
			mod.IsOnSolder = true
			return
		}
	}
}
