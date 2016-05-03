package main

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

type (
	ip       string
	zone     string
	location string
)

func DefaultTime(req *http.Request) string {
	clientIp := ipFrom(req)
	clientZone := zoneFromGeo(clientIp)
	clientLocation := locationFrom(clientZone)
	clientNow := time.Now().UTC().In(clientLocation)
	return thirtyMinutesAfter(clientNow)
}

func thirtyMinutesAfter(now time.Time) string {
	inThirtyMinutes := now.Add(time.Duration(30) * time.Minute)
	return inThirtyMinutes.Format("15:04")
}

func ipFrom(req *http.Request) ip {
	clientIpWithPort := req.RemoteAddr
	clientIp, _, _ := net.SplitHostPort(clientIpWithPort)
	return ip(clientIp)
}

func zoneFromGeo(clientIp ip) zone {
	resp, err := http.Get("https://freegeoip.net/json/" +
		string(clientIp))
	info("error when geolocating client ip: %v", err)
	defer resp.Body.Close()
	clientZone, err := jsonToZone(resp.Body)
	info("error parsing client ip: %v", err)
	return zone(clientZone)
}

func locationFrom(zn zone) *time.Location {
	loc, err := time.LoadLocation(string(zn))
	info("error processing client time zone: %v", err)
	return loc
}

func jsonToZone(resp io.ReadCloser) (zone, error) {
	type client struct {
		// As specified by freegeoip.net JSON API
		time_zone string
	}
	c := new(client)
	decoder := json.NewDecoder(resp)
	if err := decoder.Decode(&c); err != nil {
		return "", err
	}
	return zone(c.time_zone), nil
}

func info(format string, err error) {
	if err != nil {
		log.Printf(format, err)
	}
}
