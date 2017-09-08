package graphic

import (
	"io/ioutil"
	"log"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Resource pack for tortoise enemy
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var tortoiseResPackRegistry map[string]*TortoiseResPack = make(map[string]*TortoiseResPack)

type TortoiseResPack struct {
	ResLeft0  Resource
	ResLeft1  Resource
	ResRight0 Resource
	ResRight1 Resource
}

func GetTortoiseResPack(userID string) *TortoiseResPack {
	resPack, ok := tortoiseResPackRegistry[userID]
	if !ok {
		log.Fatalf("failed to get tortoise res pack for %s", userID)
	}
	return resPack
}

func registerAllTortoiseResPack() {
	files, err := ioutil.ReadDir("assets/faces")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		userID := strings.Split(f.Name(), ".")[0]
		registerTortoiseResPack(userID)
	}
}

func registerTortoiseResPack(userID string) {
	resRight0 := getFacedResource("assets/tortoise-red-right-0.png", userID,
		tortoise_walking_width, tortoise_walking_height, 35, 45, 15, 0, 5, false, false)
	resRight1 := getFacedResource("assets/tortoise-red-right-1.png", userID,
		tortoise_walking_width, tortoise_walking_height, 35, 45, 15, 0, -5, false, false)
	resLeft0 := getFacedResource("assets/tortoise-red-right-0.png", userID,
		tortoise_walking_width, tortoise_walking_height, 35, 45, 0, 0, 5, true, false)
	resLeft1 := getFacedResource("assets/tortoise-red-right-1.png", userID,
		tortoise_walking_width, tortoise_walking_height, 35, 45, 0, 0, -5, true, false)
	resPack := &TortoiseResPack{
		ResLeft0:  resLeft0,
		ResLeft1:  resLeft1,
		ResRight0: resRight0,
		ResRight1: resRight1,
	}
	tortoiseResPackRegistry[userID] = resPack
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Resource pack for boss B
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var bossBResPackRegistry map[string]*BossBResPack = make(map[string]*BossBResPack)

var bossBUserIDs map[string]struct{} = map[string]struct{}{
	"chran":   {},
	"fchen5":  {},
	"xhao":    {},
	"qingyli": {},
}

type BossBResPack struct {
	ResLeft0  Resource
	ResLeft1  Resource
	ResRight0 Resource
	ResRight1 Resource
}

func GetBossBResPack(userID string) *BossBResPack {
	resPack, ok := bossBResPackRegistry[userID]
	if !ok {
		log.Fatalf("failed to get boss B res pack for %s", userID)
	}
	return resPack
}

func registerAllBossBResPack() {
	files, err := ioutil.ReadDir("assets/faces")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		userID := strings.Split(f.Name(), ".")[0]
		if _, ok := bossBUserIDs[userID]; ok {
			registerBossBResPack(userID)
		}
	}
}

func registerBossBResPack(userID string) {
	resRight0 := getFacedResource("assets/boss-b-right-0.png", userID,
		100, 120, 50, 65, 20, 0, 5, false, false)
	resRight1 := getFacedResource("assets/boss-b-right-1.png", userID,
		100, 120, 50, 65, 20, 0, -5, false, false)
	resLeft0 := getFacedResource("assets/boss-b-right-0.png", userID,
		100, 120, 50, 65, 15, 0, 5, true, false)
	resLeft1 := getFacedResource("assets/boss-b-right-1.png", userID,
		100, 120, 50, 65, 15, 0, -5, true, false)
	resPack := &BossBResPack{
		ResLeft0:  resLeft0,
		ResLeft1:  resLeft1,
		ResRight0: resRight0,
		ResRight1: resRight1,
	}
	bossBResPackRegistry[userID] = resPack
}
