package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

const oops = "Oops! Something went wrong, please try again"

// ContainerLogsCmd show container logs
func ContainerLogsCmd(id string) {
	logs, err := GetContainerLogs(id)
	if err != nil {
		logger.Info(err)
		fmt.Println(oops)
		return
	}

	fmt.Println(logs)
}

// ContainersMetricsCmd show containers metrics
func ContainersMetricsCmd(id []string) {
	metrics, err := GetContainersMetrics(strings.Join(id, " "))
	if err != nil {
		logger.Info(err)
		fmt.Println(oops)
		return
	}

	// if metrics length == 0 -> no running containers
	if len(metrics) == 0 {
		fmt.Println("No running containers")
		return
	}

	// print metrics table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "CPU, %", "MEM, %"})
	for _, m := range metrics {
		table.Append([]string{
			m.Name,
			strconv.FormatFloat(m.CPUPercentage, 'f', 2, 64),
			strconv.FormatFloat(m.MemoryPercentage, 'f', 2, 64),
		})
	}
	table.Render()
}

// StoppedContainersCmd show stopped containers
func StoppedContainersCmd() {
	stopped, err := GetStoppedContainers()
	if err != nil {
		logger.Info(err)
		fmt.Println(oops)
		return
	}

	// if first stopped array element == "no stopped containers"
	// -> no stopped containers
	if stopped == nil {
		fmt.Println("No stopped containers")
		return
	}

	for _, s := range stopped {
		fmt.Println(s)
	}
	fmt.Println("Total stopped:", len(stopped))
}

// LaunchedContainersCmd show launched containers
func LaunchedContainersCmd() {
	launched, err := GetLaunchedContainers()
	if err != nil {
		logger.Info(err)
		fmt.Println(oops)
		return
	}

	// if first launched array element == "no launched containers"
	// -> no launched containers
	if launched == nil {
		fmt.Println("No launched containers")
		return
	}

	for _, l := range launched {
		fmt.Println(l)
	}
	fmt.Println("Total launched:", len(launched))
}

// APIStatusCmd show API status
func APIStatusCmd() {
	const APIStatus = "API status:"

	err := GetAPIStatus()
	if err != nil {
		logger.Info(err)
		fmt.Println(APIStatus, 500)
		return
	}
	fmt.Println(APIStatus, 200)
}
