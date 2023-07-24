package main

import (
	"github.com/R00tendo/goramq/webserver"
	"github.com/charmbracelet/log"
	"os"
	"errors"
	"flag"
	"os/signal"
	"syscall"
	"fmt"
	"sync"
	"bufio"
	"strings"
	"io"
	"runtime"
	"time"
	"bytes"
)


var loaded []byte

func load_qfiles(filenames *string) {
	var loadedS  int
	var loadedb  []byte
	var qfilehs  []*os.File
	var qfilesis []int
	str_filenames := strings.Split(*filenames, "::")

	//File handle and info
	for _, filename := range str_filenames {
		qfileh, err := os.Open(filename)
		if err != nil {
			log.Error("Failed to load file.", "ERROR", err)
			os.Exit(0)
		}
		defer qfileh.Close()
		qfilehs = append(qfilehs, qfileh)

		_qfileS, _ := os.Stat(filename)
		qfilesis    = append(qfilesis, int(_qfileS.Size()))
	}
	////////

	//Read to memory (bytes)
	log.Info("Loading files into memory...")

	var fbuffer []byte = make([]byte, block_size)

	for ind, qfileh := range qfilehs {
		log.Info(fmt.Sprintf("Loading file: %s", str_filenames[ind]))

		for  {
			loadedS += block_size
			if !quiet {
					fmt.Printf("Loaded %d/%d\r", loadedS, qfilesis[ind])
			}
			read, err := qfileh.Read(fbuffer)
			if err == io.EOF {
				break
			}

			loadedb = append(loadedb, fbuffer[:read]...)
		}
		fmt.Printf("%s\r", strings.Repeat("  ", 20))
		log.Info(fmt.Sprintf("File \"%s\" successfully loaded!", str_filenames[ind]))

		if !nogc && !pgco {
			runtime.GC()
			log.Info(fmt.Sprintf("Garbage collecting for %ds (cleaning up previous file)",  loadedS/100000000*2))
			ctime := time.Now()
			for {
				if time.Since(ctime) > time.Duration(loadedS/100000000*2) *time.Second {
					break
				}
			}	
		}
	}

	if !nogc && pgco {
		runtime.GC()
		log.Info(fmt.Sprintf("Garbage collecting for %ds",  loadedS/100000000*2))
		ctime := time.Now()
		for {
			if time.Since(ctime) > time.Duration(loadedS/100000000*2) *time.Second {
				break
			}
		}	
	}
	//////////

	fmt.Printf("%s\r", strings.Repeat("  ", 20))
	log.Info("All files loaded, server ready")

	loaded = loadedb
}


func query_backend() {
	var temp_results      []string
	
	//Query the data and send html response
	for {
		var output_settings webserver.Squery_settings
		temp_results     = []string{}
		query_settings  := <- querych
		search_terms    := strings.Split(query_settings.Query, "::")

		log.Info(fmt.Sprintf("New query:%s", strings.Join(search_terms, ",")))

		sscanner := bufio.NewScanner(bytes.NewReader(loaded))
		for sscanner.Scan() {
			matching_trms := 0
			line := sscanner.Bytes()

			for _, search_term := range search_terms {
				if query_settings.Case_insensitive {
					if bytes.Contains(bytes.ToLower(line), bytes.ToLower([]byte(search_term))) {
						matching_trms += 1
					}
				} else if bytes.Contains(line, []byte(search_term)) {
					matching_trms += 1
				} 
			}

			if matching_trms == len(search_terms) {
				temp_results = append(temp_results, string(line))
				if query_settings.Result_amount != 0 && (len(temp_results) >= query_settings.Result_amount) {
					break
				}
				if len(temp_results) > result_limit {
					break
				}
			}
		}

		if len(temp_results) > result_limit {
			log.Warn(fmt.Sprintf("Too many results %d", len(temp_results)))
			output_settings.Output = "<h1>Too many results</h1>"
			querych <- output_settings
			continue
		} 

		log.Info(fmt.Sprintf("Search completed, %d results", len(temp_results)))
		
		output_settings.Output = ""
		if len(temp_results) == 0 {
			output_settings.Output += "No results"
		} else {
			for _, line := range temp_results {
				output_settings.Output += fmt.Sprintf("%s\n", line)
			}
		}

		querych <- output_settings
	}
	///////////

}


var querych chan webserver.Squery_settings
var quiet        bool 
var nogc         bool
var pgco         bool
var result_limit int 
var block_size   int
var logo string = string("\033[36m") + `
________________ ________ _______ ______  __________ 
__  ____/__  __ \___  __ \___    |___   |/  /__  __ \
_  / __  _  / / /__  /_/ /__  /| |__  /|_/ / _  / / /
/ /_/ /  / /_/ / _  _, _/ _  ___ |_  /  / /  / /_/ / 
\____/   \____/  /_/ |_|  /_/  |_|/_/  /_/   \___\_\ 
` + string("\033[0m")

func main() {
	//Variables
	killch  := make(chan bool, 1)
	var killwg sync.WaitGroup

	querych =  make(chan webserver.Squery_settings, 1)
	//////////

	//Setup keyboard interrupt handler
	kbintr := make(chan os.Signal, 1)
	signal.Notify(kbintr, os.Interrupt, os.Kill, syscall.SIGTERM)
	killwg.Add(1)
	go func() {
		<- kbintr
		signal.Stop(kbintr)
		killch <- true
		<- killch
		log.Info("Exiting...")
		killwg.Done()
		os.Exit(0)
	}()
	//////////

	fmt.Println(logo)

	//Flag arguments
	queryfiles := flag.String("filenames", "", "Files to serve (separated with ::)")
	ladr       := flag.String("listener", "127.0.0.1:9112", "Web server listener (IP:PORT)")
	api_pass   := flag.String("password", "", "Password protects the API")
	flag.BoolVar(&quiet, "quiet", false, "When used, Goramq will not display the file loading progress (perfomance boost)")
	flag.IntVar(&block_size, "blocksize", 1024, "How big chunks the files will be loaded in")
	flag.IntVar(&result_limit, "result-limit", 50000, "How many results you can receive (over limit = error message)")
	flag.BoolVar(&pgco, "pcgo", false, "Only does garbage collection after all files are loaded.")
	flag.BoolVar(&nogc, "nogc", false, "Disables garbage collection (uses twice as much ram)")
	flag.Parse()

	for _, queryfile := range strings.Split(*queryfiles, "::") {
		if _, err := os.Stat(queryfile); errors.Is(err, os.ErrNotExist) {
			log.Error("File doesn't exist.", "FILENAME", queryfile)
			os.Exit(0)
		}
	}
	//////////

	//Start processes
	load_qfiles(queryfiles)
	go query_backend()
	webserver.Start(killch, querych, ladr, *api_pass)
	//////////

	killwg.Wait()
}