package main

import (
	"log"
	"net"
	"time"
)

var ipv4net *net.IPNet
var ipv6net *net.IPNet

// Helper function to check if an ip address is within our subnets
func isInSubnet(ip net.IP) bool {
	if ipv4net == nil {
		_, ipv4net, _ = net.ParseCIDR("128.153.144.0/23")
		_, ipv6net, _ = net.ParseCIDR("2605:6480:c051::1/48")
	}

	return ipv4net.Contains(ip) || ipv6net.Contains(ip)
}

// Returns the next Wednesday in YYYYMMDD format
// If today is a wednesday today's date is returned
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
