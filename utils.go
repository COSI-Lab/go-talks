package main

import (
	"log"
	"net"
	"time"
)

// Networks is a slice of net.IPNet
type Networks struct {
	nets []*net.IPNet
}

// NewNetworks creates a new Networks struct from a slice of subnets
func NewNetworks(subnets []string) Networks {
	n := Networks{}
	for _, subnet := range subnets {
		_, net, err := net.ParseCIDR(subnet)
		if err != nil {
			log.Fatalf("[FATAL] Invalid subnet %s: %s", subnet, err)
		}
		n.nets = append(n.nets, net)
	}
	return n
}

// Contains returns true if the given IP is in one of the trusted subnets
// Runs in O(the number of trusted subnets)
func (n *Networks) Contains(ip net.IP) bool {
	for _, net := range n.nets {
		if net.Contains(ip) {
			return true
		}
	}
	return false
}

// Returns the next Wednesday in YYYYMMDD format
func nextWednesday() string {
	const format = "20060102"
	now := time.Now().In(tz)

	daysUntilWednesday := time.Wednesday - now.Weekday()
	if daysUntilWednesday == 0 {
		return now.Format(format)
	} else if daysUntilWednesday > 0 {
		return now.AddDate(0, 0, int(daysUntilWednesday)).Format(format)
	} else {
		return now.AddDate(0, 0, int(daysUntilWednesday)+7).Format(format)
	}
}

func addWeek(week string) string {
	const format = "20060102"

	t, err := time.ParseInLocation(format, week, tz)
	if err != nil {
		log.Fatalln("Failed to parse week:", err)
	}

	return t.AddDate(0, 0, 7).Format(format)
}

func subtractWeek(week string) string {
	const format = "20060102"

	t, err := time.ParseInLocation(format, week, tz)
	if err != nil {
		log.Fatalln("Failed to parse week:", err)
	}

	return t.AddDate(0, 0, -7).Format(format)
}

func weekForHumans(week string) (string, error) {
	const format = "20060102"

	t, err := time.ParseInLocation(format, week, tz)
	if err != nil {
		return "", err
	}

	return t.Format("January _2, 2006"), nil
}

func isPast(current, week string) bool {
	const format = "20060102"

	currentTime, err := time.ParseInLocation(format, current, tz)
	if err != nil {
		log.Fatalln("Failed to parse week:", err)
	}

	weekTime, err := time.ParseInLocation(format, week, tz)
	if err != nil {
		log.Fatalln("Failed to parse week:", err)
	}

	return currentTime.After(weekTime)
}

func duringMeeting() bool {
	return time.Now().In(tz).Weekday() == time.Wednesday && time.Now().In(tz).Hour() >= 19
}
