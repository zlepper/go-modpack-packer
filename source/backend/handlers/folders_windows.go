package handlers

import (
	"os"
	"sync"
)

// Gets the available disk drives on the machine
func GetDrives() (disks []string, err error) {
	// I know it's not elegant. I don't care.
	var wg sync.WaitGroup
	var lock sync.Mutex
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		wg.Add(1)
		go func(drive rune) {
			_, err := os.Open(string(drive) + ":\\")
			if err == nil {
				lock.Lock()
				disks = append(disks, string(drive)+":")
				lock.Unlock()
			}
			wg.Done()
		}(drive)
	}
	wg.Wait()
	return
}
