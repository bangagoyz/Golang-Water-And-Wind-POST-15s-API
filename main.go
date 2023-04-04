package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{}
	ticker := time.NewTicker(15 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				randWind := rand.Intn(100)
				randWater := rand.Intn(100)
				data := map[string]interface{}{
					"Water": randWater,
					"Wind":  randWind,
				}

				requestJson, err := json.Marshal(data)
				if err != nil {
					log.Fatal(err)
					fmt.Println("error marshal")
				}

				req, err := http.NewRequest("POST", "https://jsonplaceholder.typicode.com/posts",
					bytes.NewBuffer(requestJson))

				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					log.Fatal(err)
					fmt.Println("error request")
				}

				res, err := client.Do(req)
				if err != nil {
					log.Fatal(err)
					fmt.Println("error response")
				}
				defer res.Body.Close()

				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					log.Fatal(err)
					fmt.Println("error read")
				}
				log.Println(string(body))

				var responseMap map[string]interface{}
				err = json.Unmarshal(body, &responseMap)
				if err != nil {
					log.Fatal(err)
					fmt.Println("error unmarshal")
				}

				delete(responseMap, "id")
				var windStat, waterStat string
				for key, value := range responseMap {
					switch v := value.(type) {
					case float64:
						intValue := int(v)
						if key == "Water" && intValue <= 5 {
							waterStat = ("Water status: aman")
						} else if key == "Water" && intValue <= 8 {
							waterStat = ("Water status: siaga")
						} else if key == "Water" && intValue > 8 {
							waterStat = ("Water status: bahaya!")
						}

						if key == "Wind" && intValue <= 6 {
							windStat = ("Wind status: aman")
						} else if key == "Wind" && intValue <= 15 {
							windStat = ("Wind status: siaga")
						} else if key == "Wind" && intValue > 15 {
							windStat = ("Wind status: bahaya!")
						}
					}

				}
				fmt.Println(waterStat)
				fmt.Println(windStat)
			}
		}
	}()

	time.Sleep(5 * time.Minute)
	ticker.Stop()
	done <- true
	fmt.Println("Loop stopped")
}
