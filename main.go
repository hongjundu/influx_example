package main

import (
	"encoding/json"
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/satori/go.uuid"
	"log"
	"math/rand"
	"time"
)

func main() {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
	}
	defer c.Close()

	//write(c)
	//queryByTime(c)
	//queryByTimeAndSN(c)
	queryMeanByTimeAndSN(c)
}

func write(c client.Client) {
	now := time.Now()
	tm := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
	count := 0

	snList := []string{"SN0001", "SN0002", "SN003"}

	tm = tm.Add(-time.Hour * 10000)

	for count < 10000 {
		// Create a new point batch
		bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  "mydb",
			Precision: "s",
		})

		for _, sn := range snList {
			uid := uuid.NewV4()
			tags := map[string]string{
				"SN": sn,
			}
			fields := map[string]interface{}{
				"ID": uid.String(),
			}

			i := 1
			for i <= 30 {
				valueInt := rand.Intn(i * 10000)
				value := float32(valueInt) / 100
				fields[fmt.Sprintf("电流%d", i)] = value
				i++
			}

			pt, err := client.NewPoint("example", tags, fields, tm)
			if err == nil {
				bp.AddPoint(pt)
			} else {
				log.Printf("NewPoint error: %v", err)
			}

			err = c.Write(bp)
			if err != nil {
				log.Printf("Write error: %v", err)
			}
		}

		tm = tm.Add(time.Hour)

		count++
	}
}

func queryByTime(c client.Client) {
	now := time.Now()
	start := time.Date(now.Year(),now.Month(),now.Day(),0,0,0,0,now.Location())
	start = start.Add(-10*time.Hour*24)
	end := start.Add(time.Hour*24)

	cmd := fmt.Sprintf("select * from example where time>=%v and time<=%v TZ('Asia/Shanghai')",start.UnixNano(),end.UnixNano())

	q := client.Query{
		Command:  cmd,
		Database: "mydb",
	}
	if response, err := c.Query(q); err == nil && response.Error() == nil {

		b, _ := json.Marshal(response.Results)
		fmt.Printf("%s",string(b))

	}else {
		log.Printf("err: %v response.Error: %v", err, response.Error())
	}
}

func queryByTimeAndSN(c client.Client) {
	now := time.Now()
	start := time.Date(now.Year(),now.Month(),now.Day(),0,0,0,0,now.Location())
	start = start.Add(-10*time.Hour*24)
	end := start.Add(time.Hour*24)

	cmd := fmt.Sprintf("select * from example where SN='SN0001' and time>=%v and time<=%v TZ('Asia/Shanghai')",start.UnixNano(),end.UnixNano())

	q := client.Query{
		Command:  cmd,
		Database: "mydb",
	}
	if response, err := c.Query(q); err == nil && response.Error() == nil {

		b, _ := json.Marshal(response.Results)
		fmt.Printf("%s",string(b))

	}else {
		log.Printf("err: %v response.Error: %v", err, response.Error())
	}
}

func queryMeanByTimeAndSN(c client.Client) {
	now := time.Now()
	start := time.Date(now.Year(),now.Month(),now.Day(),0,0,0,0,now.Location())
	start = start.Add(-10*time.Hour*24)
	end := start.Add(time.Hour*24)

	cmd := fmt.Sprintf("select mean(\"电流1\") from example where SN='SN0001' and time>=%v and time<=%v TZ('Asia/Shanghai')",start.UnixNano(),end.UnixNano())

	q := client.Query{
		Command:  cmd,
		Database: "mydb",
	}
	if response, err := c.Query(q); err == nil && response.Error() == nil {

		b, _ := json.Marshal(response.Results)
		fmt.Printf("%s",string(b))

	}else {
		log.Printf("err: %v response.Error: %v", err, response.Error())
	}
}
