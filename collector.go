package main
import (
	"log"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

//Define variables hold the webserver info
var (
	EFWebServerIP = ""
	EFWebServerPort = ""
)
//Define a struct for you collector that contains pointers
//to prometheus descriptors for each metric you wish to expose.
//Note you can also include fields of other types if they provide utility
//but we just won't be exposing them as metrics.
type PC2Collector struct {
	NumPendingMetric *prometheus.Desc
	NumSolvedMetric *prometheus.Desc
	ProbStats *prometheus.Desc
}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func newPC2Collector(ServerIP string, ServerPort string) *PC2Collector {
	EFWebServerIP = ServerIP
	EFWebServerPort = ServerPort
	return &PC2Collector{
		NumPendingMetric: prometheus.NewDesc("NumPending_metric",
			"Number of pending submissions waitting",
			nil, nil,
		),
		NumSolvedMetric: prometheus.NewDesc("NumSolved_metric",
			"Number of solved problems",
			nil, nil,
		),
		ProbStats: prometheus.NewDesc("problem_statistics",
                        "Problem statistics",
                        []string{"problem_id","metric_type"}, nil,
                ),
	}
}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *PC2Collector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.NumPendingMetric
	ch <- collector.NumSolvedMetric
	ch <- collector.ProbStats
}

//Collect implements required collect function for all promehteus collectors
func (collector *PC2Collector) Collect(ch chan<- prometheus.Metric) {
	var (
		NP int
		NS int
		NumPending float64
		NumSolved float64
		ProbTotal map[string]int
		ProbSolved map[string]int
	)

	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor or call other functions that do so.
	//var metricValue float64
	//if 1 == 1 {
		//metricValue = 1
	//}

	ProbTotal = make(map[string]int)
	ProbSolved = make(map[string]int)

	scoreboardPage := "https://" + EFWebServerIP + ":" + EFWebServerPort + "/contest/scoreboard"
	fmt.Println("Address ", scoreboardPage)
	data,err := readJSONfile("scoreboard.json")
         if err != nil {
                log.Fatal(err)
        }
	NP, NS , ProbTotal , ProbSolved  = parseScoreboard(data)
	NumPending = float64(NP)
	NumSolved = float64(NS)

	//Write latest value for each metric in the prometheus metric channel.
	//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
	ch <- prometheus.MustNewConstMetric(collector.NumPendingMetric, prometheus.GaugeValue, NumPending)
	ch <- prometheus.MustNewConstMetric(collector.NumSolvedMetric, prometheus.CounterValue, NumSolved)
	//Change the loop to loop through map keys 
	for key, value := range ProbTotal {
		ch <- prometheus.MustNewConstMetric(collector.ProbStats, prometheus.CounterValue, float64(value) , key,"submission")
	}
	for key, value := range ProbSolved {
                ch <- prometheus.MustNewConstMetric(collector.ProbStats, prometheus.CounterValue, float64(value) , key,"solved")
        }	
}
//ProbStats *prometheus.Desc
        //ProbSolved *prometheus.Desc

