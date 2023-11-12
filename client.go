package iam

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/patrickmn/go-cache"

	"github.com/imroc/req/v3"
	http_mw "github.com/zitadel/zitadel-go/v2/pkg/api/middleware/http"
)

var (
	UserId            string
	IsAdmin           bool
	introspection     *http_mw.IntrospectionInterceptor
	once              sync.Once
	userCacheInstance *cache.Cache
	userOnceCache     sync.Once
)

func NewZitadelClient(baseUrl string, keyPath string, user *ZitadelUser, userId *string) *ZitadelClient {
	once.Do(func() {
		var err error
		introspection, err = http_mw.NewIntrospectionInterceptor(baseUrl, keyPath)
		if err != nil {
			log.Fatal(err)
		}
	})

	return &ZitadelClient{
		IntrospectionInterceptor: introspection,
		ZitadelUser:              user,
		UserId:                   userId,
		Client: req.C().
			SetBaseURL(baseUrl).
			//SetCommonErrorResult(&ErrorMessage{}).
			EnableDumpEachRequest().
			OnAfterResponse(func(client *req.Client, resp *req.Response) error {
				if resp.Err != nil { // There is an underlying error, e.g. network error or unmarshal error.
					return nil
				}
				/*
					if errMsg, ok := resp.ErrorResult().(*ErrorMessage); ok {
						resp.Err = errMsg // Convert api error into go error
						return nil
					} */
				if !resp.IsSuccessState() {
					// Neither a success response nor a error response, record details to help troubleshooting
					resp.Err = fmt.Errorf("bad status: %s\nraw content:\n%s", resp.Status, resp.Dump())
				}
				return nil
			}),
	}
}

func (client *ZitadelClient) ZitadelAuth(next http.Handler) http.Handler {
	return client.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		cacheKey := md5hash((token))
		store, ok := getCache().Get(cacheKey)
		if ok {
			log.Println("==== get user from cache")
			client = store.(*ZitadelClient)
		} else {
			log.Println("==== get user from auth server")
			resp, _ := client.R().
				SetContext(r.Context()).
				SetHeader("Authorization", token).
				SetSuccessResult(&client.ZitadelUser).
				Get("/oidc/v1/userinfo")
			if resp.IsSuccessState() {
				fmt.Println("Get ZitadelUser")
				fmt.Println(client.ZitadelUser)
				*client.UserId = md5hash(client.ZitadelUser.Email)
				getCache().Set(cacheKey, client, cache.DefaultExpiration)
			}
		}
		next.ServeHTTP(w, r)
	})
}
func md5hash(text string) string {
	algorithm := md5.New()
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}
