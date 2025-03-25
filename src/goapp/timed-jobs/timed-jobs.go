package timedjobs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/config"
	"main/model"
	"main/service"
	"net/http"
	"strconv"
	"time"
)

type timedJobs struct {
	Service       *service.Service
	configManager config.ConfigManager
}

func NewTimedJobs(s *service.Service, configManager config.ConfigManager) TimedJobs {
	return &timedJobs{
		Service:       s,
		configManager: configManager,
	}
}

func (t *timedJobs) ReprocessFailedCallbacks() {
	freq := t.configManager.GetCallbackRetryFreq()
	freqInt, _ := strconv.ParseInt(freq, 0, 64)
	if freqInt > 0 {
		for range time.NewTicker(time.Duration(freqInt) * time.Minute).C {
			f, err := t.Service.Item.GetFailedCallbacks()
			if err != nil {
				fmt.Printf("Failed to get failed callbacks: %v", err.Error())
				return
			}

			for _, id := range f {
				go t.postCallback(id)
			}
		}
	}
}

func (t *timedJobs) postCallback(id string) {
	item, err := t.Service.Item.GetItemById(id)
	if err != nil {
		fmt.Println("Error getting item by id: ", id)
		return
	}

	if item.CallbackUrl == "" {
		fmt.Println("No callback url found")
		return
	} else {
		params := model.ResponseCallback{
			ItemId:       id,
			IsApproved:   item.IsApproved,
			Remarks:      item.ApproverRemarks,
			ResponseDate: item.DateResponded,
			RespondedBy:  item.RespondedBy,
		}

		jsonReq, err := json.Marshal(params)
		if err != nil {
			return
		}

		token, err := t.Service.Authenticator.GenerateToken()
		if err != nil {
			fmt.Println("Error generating token: ", err)
			return
		}

		req, err := http.NewRequest("POST", item.CallbackUrl, bytes.NewBuffer(jsonReq))
		if err != nil {
			return
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{
			Timeout: time.Second * 90,
		}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Error posting callback: ", err)
			return
		}

		isCallbackFailed := res.StatusCode != 200

		err = t.Service.Item.UpdateItemCallback(id, isCallbackFailed)
		if err != nil {
			fmt.Println("Error updating item callback: ", err)
			return
		}
	}
}
