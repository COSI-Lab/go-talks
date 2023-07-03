package main

import (
	"log"
	"net"
	"time"
)

type Networks struct {
	nets []*net.IPNet
}

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
	now := time.Now().In(TZ)

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

	t, err := time.ParseInLocation(format, week, TZ)
	if err != nil {
		log.Fatalln("Failed to parse week:", err)
	}

	return t.AddDate(0, 0, 7).Format(format)
}

func subtractWeek(week string) string {
	const format = "20060102"

	t, err := time.ParseInLocation(format, week, TZ)
	if err != nil {
		log.Fatalln("Failed to parse week:", err)
	}

	return t.AddDate(0, 0, -7).Format(format)
}

func weekForHumans(week string) (string, error) {
	const format = "20060102"

	t, err := time.ParseInLocation(format, week, TZ)
	if err != nil {
		return "", err
	}

	return t.Format("January _2, 2006"), nil
}

func isPast(current, week string) bool {
	const format = "20060102"

	currentTime, err := time.ParseInLocation(format, current, TZ)
	if err != nil {
		log.Fatalln("Failed to parse week:", err)
	}

	weekTime, err := time.ParseInLocation(format, week, TZ)
	if err != nil {
		log.Fatalln("Failed to parse week:", err)
	}

	return currentTime.After(weekTime)
}

func duringMeeting() bool {
	return time.Now().In(TZ).Weekday() == time.Wednesday && time.Now().In(TZ).Hour() >= 19
}
