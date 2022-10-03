package main

import (
	"flag"
	"fmt"
	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/keychain"
	"log"
	"time"
)

var (
	wf     *aw.Workflow
	action string
)

const (
	cacheDir            = "./cache"
	cacheKey            = "time"
	actionSetCredential = "set-credential"
	actionLogCredential = "log-credential"
)

func init() {
	wf = aw.New()
	flag.StringVar(&action, "action", "", "action to be executed")
}

func fetchTime() (time.Time, error) {
	cache := aw.NewCache(cacheDir)
	if !cache.Expired(cacheKey, time.Second*20) {
		bytes, err := cache.Load(cacheKey)
		if err != nil {
			return time.Time{}, err
		}
		return time.Parse(time.RFC3339, string(bytes))
	}
	now := time.Now()
	cache.Store(cacheKey, []byte(now.Format(time.RFC3339)))
	return now, nil
}

func run() {
	wf.Args()
	flag.Parse()

	kc := keychain.New("test-alfred-workflow")

	if action == "" {
		wf.NewItem("Set credential").Arg(actionSetCredential, flag.Arg(0)).Valid(len(flag.Arg(0)) > 0)
		wf.NewItem("Log credential").Arg(actionLogCredential).Valid(true)

		t, err := fetchTime()
		if err != nil {
			wf.FatalError(err)
		}
		wf.NewItem(fmt.Sprintf("now: %s, cache: %s", time.Now().Format("15:04:05"), t.Format("15:04:05")))
		wf.SendFeedback()
		return
	}

	switch action {
	case actionSetCredential:
		kc.Set("test", flag.Arg(0))
	case actionLogCredential:
		cred, err := kc.Get("test")
		if err != nil {
			wf.FatalError(err)
		}
		log.Println(fmt.Sprintf("credenital: %s", cred))
	}
	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
