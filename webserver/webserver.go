package webserver

import (
	"fmt"
	"net/http"
	"strconv"
	"bufio"
	"sync"
	"strings"
	"github.com/charmbracelet/log"
)

type Squery_settings struct {
	Case_insensitive bool
	Query            string 
	Result_amount    int
	Output           string
}

var being_read sync.Mutex

func queryHandle(wr http.ResponseWriter, re *http.Request) {
	being_read.Lock()
	defer being_read.Unlock()
	var raw_output       string = ""
	var processed_output []string
	var query_settings   Squery_settings

	if re.URL.Query().Get("pass") != api_pass {
		fmt.Fprintf(wr, http_errors["AccDnd"])
		return
	}
	query_settings.Query             = re.URL.Query().Get("q")
	case_insensitive                := re.URL.Query().Get("caseins")
	query_settings.Result_amount, _  = strconv.Atoi(re.URL.Query().Get("amount"))

	if len(case_insensitive) != 0 {
		query_settings.Case_insensitive = true
	}

	if len(query_settings.Query) == 0 {
		fmt.Fprintf(wr, http_errors["NoParam"])
		
	} else {
		querych      <- query_settings
		raw_output = (<- querych).Output

		qscanner := bufio.NewScanner(strings.NewReader(raw_output))
		for qscanner.Scan() {
			processed_output = append(processed_output, qscanner.Text())
		}
		fmt.Fprintf(wr, strings.Join(processed_output, "\n"))
	}
}

var http_errors map[string]string
var querych     chan Squery_settings 
var api_pass    string

func Start(killch chan bool, _querych chan Squery_settings, ladr *string, _api_pass string) {
	querych  = _querych
	api_pass = _api_pass

	//Setup predefined errorz
	http_errors = map[string]string {
		"NoParam": "<h1>Specify a search term</h1>",
		"AccDnd": "<h1>Access denied</h1>",
	}
	////////
	
	//Web server settings
	httpsvr := http.Server{Addr: *ladr,}

	http.HandleFunc("/search", queryHandle)
	////////

	//Receive signal when main program wants to exit
	go func() {
		<- killch
		log.Info("Stopping web server...")
		httpsvr.Close()
	}()
	////////

	//Start Handler
	if err := httpsvr.ListenAndServe(); err != http.ErrServerClosed {
		log.Error("HTTP server error.", "ERROR", err)
	}
	////////

	//Informs main that the program has cleaned up and is ready to exit
	killch <- false
	////////
}