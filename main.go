package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"mta_app/config"
	"mta_app/email"

	"github.com/tevino/abool"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/yaml.v2"
)

const (
	baseURL = "http://api.prod.obanyc.com/api/siri/vehicle-monitoring.json?"
)

type mtaConfig struct {
	configFile string
}

func main() {
	rand.Seed(time.Now().UnixNano())

	flags := parseFlags()

	yamlFile, err := ioutil.ReadFile(flags.configFile)
	if err != nil {
		panic(err)
	}

	var cfg config.Config

	if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
		panic(err)
	}

	emailClient := email.NewEmailUser(cfg.Email)

	reqURL := constructURL(cfg.MTA)
	jsonResp, err := getLocation(reqURL)
	if err != nil {
		return
	}

	var withinTime abool.AtomicBool
	go checkTime(withinTime, cfg.MTA.BeginTime, cfg.MTA.EndTime)

	for {
		if withinTime.IsSet() {
			found := findClosestBus(jsonResp, cfg.MTA.StopCheck)
			if found {
				if err = emailClient.SendEmail(); err != nil {
					panic(err)
				}
				return
			}
		}
		time.Sleep(2 * time.Second)
	}
}

func checkTime(withinTime abool.AtomicBool, begin, end int) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			now := time.Now().In(loc)
			startBound := time.Date(now.Year(), now.Month(), now.Day(), begin, 0, 0, 0, now.Location())
			endBound := time.Date(now.Year(), now.Month(), now.Day(), end, 0, 0, 0, now.Location())

			if now.After(startBound) && now.Before(endBound) && now.Weekday() != 6 && now.Weekday() != 7 {
				withinTime.SetTo(true)
			} else {
				withinTime.UnSet()
			}
		}
	}
}

func constructURL(mta config.MTAInfo) string {
	var buffer bytes.Buffer

	buffer.WriteString(baseURL)
	buffer.WriteString("key=" + mta.Key)
	buffer.WriteString("&LineRef=" + mta.Line)
	buffer.WriteString("&DirectionRef=" + mta.Direction)

	return buffer.String()
}

func getLocation(url string) (*config.MTAResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var target config.MTAResponse
	json.NewDecoder(resp.Body).Decode(&target)
	return &target, nil
}

func findClosestBus(resp *config.MTAResponse, check string) bool {
	for _, v := range resp.Siri.ServiceDelivery.VehicleMonitoringDelivery {
		for _, ind := range v.VehicleActivity {
			fmt.Println(ind.MonitoredVehicleJourney.MonitoredCall.StopPointName)
			if ind.MonitoredVehicleJourney.MonitoredCall.StopPointName == check {
				return true
			}
		}
	}
	return false
}

func parseFlags() *mtaConfig {
	cfg := mtaConfig{}
	a := kingpin.New(filepath.Base(os.Args[0]), "MTA")

	a.Version("1.0")
	a.HelpFlag.Short('h')
	a.Flag("config", "MTA configuration file").StringVar(&cfg.configFile)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing commandline arguments, got error %v\n", err)
		a.Usage(os.Args[1:])
		os.Exit(2)
	}

	return &cfg
}
