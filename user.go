package iam

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/imroc/req/v3"
	"github.com/patrickmn/go-cache"
)

var cacheInstance *cache.Cache
var onceCache sync.Once

type User struct {
	AccountID     string `json:"AccountId"`
	Email         string `json:"Email"`
	FullName      string `json:"FullName"`
	ID            string `json:"ID"`
	PublicID      string `json:"PublicId"`
	RegisterToken string `json:"RegisterToken"`
}

func getCache() *cache.Cache {
	onceCache.Do(func() {
		cacheInstance = cache.New(5*time.Minute, 10*time.Minute)
	})
	return cacheInstance

}
func GetUserId(accountId string, authHeader string, userProfileUrl string) (user *User, err error) {
	store, ok := getCache().Get(accountId)
	if ok {
		log.Println("==== get user from cache")
		return store.(*User), nil
	} else {
		log.Println("==== get user from account server")
		client := req.C().EnableDumpAll()
		resp, err := client.R().
			SetHeader("Authorization", authHeader).
			SetResult(&user).
			Get(userProfileUrl)

		if resp.IsSuccess() {
			if err == nil {
				getCache().Set(accountId, user, cache.DefaultExpiration)
			}
			fmt.Printf("%s's id is %s\n", user.FullName, user.PublicID)
		}
		return user, err
	}
}
