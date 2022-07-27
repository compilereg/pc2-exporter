package main
 
import (
    "flag"
    "log"
    "fmt"
    "net/http"
 
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
   	fmt.Println("Starting..")
	//Get parameters for EFWebServerIP, EFWebServerPort, EFUSername, EFPassword
	EFWebServerIP := flag.String("EVENTFEEDER_WEBSERVERIP","localhost","Event feeder Web server IP")
	EFWebServerPort := flag.String("EVENTFEEDER_WEBSERVERPORT","50443", "Event feeder Web server port")
	EFWebServerUser := flag.String("EVENTFEEDER_WEBSERVERUSERNAME","admin", "Event feeder Web server username")
	EFWebServerPassword := flag.String("EVENTFEEDER_WEBSERVERPASSWORD","admin", "Event feeder Web server password")
	listenAddress := flag.String("EXPORTER_PORT",":9102","PC2 Exporter TCP Port number")
	metricPath := flag.String("EXPORTER_METRICPATH","/metrics","Path under which to expose metrics")
	flag.Parse()
	fmt.Println("IP " , *EFWebServerIP)
	fmt.Println("Port ", *EFWebServerPort)
	fmt.Println("Username ", *EFWebServerUser)
	fmt.Println("Password ", *EFWebServerPassword)
	fmt.Println("Listen address", *listenAddress)
	fmt.Println("Metric path", *metricPath)

	//Call the Prometheus metric collector, and pass to it required values to scrape the pc2 webserver values
   	PC2 := newPC2Collector(*EFWebServerIP,*EFWebServerPort)
   	prometheus.MustRegister(PC2)
   	log.Fatal(serverMetrics(*listenAddress, *metricPath))
}
 
func serverMetrics(listenAddress, metricsPath string) error {
    http.Handle(metricsPath, promhttp.Handler())
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`
            <html>
            <head><title>Volume Exporter Metrics</title></head>
            <body>
            <h1>ConfigMap Reload</h1>
            <p><a href='` + metricsPath + `'>Metrics</a></p>
            </body>
            </html>
        `))
    })
    return http.ListenAndServe(listenAddress, nil)
}
