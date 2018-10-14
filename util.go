package main

import (
	"sort"
	"strconv"
)

func validPort(data string) bool {
	if port, err := strconv.Atoi(data); err != nil || port < 4000 || port > 5000 {
		return false
	}
	return true
}

func getBootstrapNode(port string) Contact {
	return Contact{
		ID:   NodeID{},
		Host: getLocalIP(),
		Port: port,
	}
}

func sortByDistance(list []Contact, id NodeID) {
	sort.SliceStable(list, func(i, j int) bool {
		d1 := Distance(id, list[i].ID)
		d2 := Distance(id, list[j].ID)
		return Less(d1, d2)
	})
}
