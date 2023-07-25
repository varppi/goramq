<img src="https://github.com/R00tendo/goramq/blob/main/images/goramq.png?raw=true"></img>

# GORAMQ (Golang Ram Query)
<b>Like the name suggests, Goramq loads user supplied files into memory and then lets you search strings inside the files.</b>

## How it works
<b>Goramq reads all files provided in chunks (defined by blocksize) and adds the read data to a variable. By default after reading the files Goramq will ask Golang garbage collector to do a cycle, this frees up memory.

When a query comes through, Goramq activates a "lock" that prevents other queries from being taken while the current one is being executed.</b>

## Installation
Note: Make sure you have Golang installed and the go bin directory added to your PATH variable.
```bash
go install github.com/R00tendo/goramq@local_version
```

## Basic usage examples
```
#Loads 2 files and searches all the queries inside "queries" 
goramq --filenames "documents/database.sql::documents/randomtextfile.txt" --search queries --output out.txt

#Loads a file and does garbage collection only after all files have been loaded
goramq --filenames "database.sql" --pcgo --search qries
```


## Help page
```
________________ ________ _______ ______  __________
__  ____/__  __ \___  __ \___    |___   |/  /__  __ \
_  / __  _  / / /__  /_/ /__  /| |__  /|_/ / _  / / /
/ /_/ /  / /_/ / _  _, _/ _  ___ |_  /  / /  / /_/ /
\____/   \____/  /_/ |_|  /_/  |_|/_/  /_/   \___\_\

Usage of goramq:
  -blocksize int
        How big chunks the files will be loaded in (default 1024)
  -caseins
        Case sensitive
  -filenames string
        Files to serve (separated with ::)
  -nogc
        Disables garbage collection (uses twice as much ram)
  -output string
        File to write the output. (default "output.txt")
  -pcgo
        Only does garbage collection after all files are loaded.
  -quiet
        Won't display file loading progress (perfomance boost)
  -resamount int
        Amount of results to return (default 500000)
  -result-limit int
        How many results you can receive (over limit = error message) (default 50000)
  -search string
        File of search queries to execute
```
