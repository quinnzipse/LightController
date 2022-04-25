package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

const appKey = "<API-KEY-HERE>"
const header = "hue-application-key"

const bridgeIpAddress = "192.168.1.26"
const resource = "light"

const baseURL = "https://" + bridgeIpAddress + "/clip/v2/resource/" + resource

var client *http.Client

func main() {

	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	listLights()
	fallAsleep()
}

func wakeUp() {
	lightOnOff("2cf19f73-a49f-4339-94e2-480181a33ab4", true)
	lightOnOff("0abf139f-d08b-403b-8e22-4c8df123d494", true)
	lightOnOff("bc86efa2-2613-41de-bcad-686aed5a5d44", true)
	time.Sleep(time.Second * 3)
	slowTransition(0, 100, []string{"2cf19f73-a49f-4339-94e2-480181a33ab4", "0abf139f-d08b-403b-8e22-4c8df123d494", "bc86efa2-2613-41de-bcad-686aed5a5d44"})
}

func fallAsleep() {
	slowTransition(100, 0, []string{"2cf19f73-a49f-4339-94e2-480181a33ab4", "0abf139f-d08b-403b-8e22-4c8df123d494", "bc86efa2-2613-41de-bcad-686aed5a5d44"})
	time.Sleep(time.Second * 10)
	lightOnOff("2cf19f73-a49f-4339-94e2-480181a33ab4", false)
	lightOnOff("0abf139f-d08b-403b-8e22-4c8df123d494", false)
	lightOnOff("bc86efa2-2613-41de-bcad-686aed5a5d44", false)
}

func listLights() {
	request, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	request.Header.Add(header, appKey)

	response, err := client.Do(request)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	var responseObject Response
	json.Unmarshal(content, &responseObject)

	fmt.Println(len(responseObject.Data))

	for i := 0; i < len(responseObject.Data); i++ {
		currentLight := responseObject.Data[i]
		var on string

		if currentLight.On.Value {
			on = "on"
		} else {
			on = "off"
		}

		fmt.Println(currentLight.Identifier, currentLight.Metadata.Name, currentLight.Metadata.Type, on)
	}
}

func fastTransition(from, to uint8, ids []string) {
	transitionBrightness(from, to, 50, time.Second, ids)
}

func slowTransition(from, to uint8, ids []string) {
	fmt.Println(time.Second * 5)
	transitionBrightness(from, to, 50, time.Second*5, ids)
}

func transitionBrightness(from, to uint8, steps int, interval time.Duration, ids []string) {

	difference := int(to) - int(from)
	absDiff := int(math.Abs(float64(difference)))
	sign := difference / absDiff

	rate := 1
	if absDiff > steps {
		rate = absDiff / steps
	}

	brightnessObject := LightState{
		Dimming: &Dimming{
			Brightness: from,
		},
	}

	for absDiff > 0 {
		if sign < 0 {
			brightnessObject.Dimming.Brightness -= uint8(rate)
		} else {
			brightnessObject.Dimming.Brightness += uint8(rate)
		}

		absDiff -= rate

		for i := 0; i < len(ids); i++ {
			put(ids[i], brightnessObject)
		}

		time.Sleep(interval)
		fmt.Println("Current Brightness", brightnessObject.Dimming.Brightness)
	}
}

func lightOnOff(id string, on bool) {
	onObject := LightState{
		On: &On{Value: on},
	}

	put(id, onObject)
}

func put(id string, state LightState) {
	putData, err := json.Marshal(state)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	requestBody := strings.NewReader(string(putData))

	request, err := http.NewRequest("PUT", baseURL+"/"+id, requestBody)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	request.Header.Add(header, appKey)

	response, err := client.Do(request)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	defer response.Body.Close()

	_, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
}
