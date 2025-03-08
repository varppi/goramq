<img src="https://github.com/SpoofIMEI/goramq/blob/main/images/goramq.png?raw=true"></img>

# GORAMQ (Golang Ram Query)
<b>Like the name suggests, Goramq loads user supplied files into memory and then lets you search strings inside the files via a restful API.</b>

## How it works
<b>Goramq reads all files provided in chunks (defined by blocksize) and adds the read data to a variable. By default after reading the files Goramq will ask Golang garbage collector to do a cycle,  this frees up memory.

When a query comes through, Goramq activates a "lock" that prevents other queries from being taken while the current one is being executed, to the person/bot doing the query that was put on hold, it will look like the website is just loading slowly.</b>

## Installation
Note: Make sure you have Golang installed and the go bin directory added to your PATH variable.
```bash
go install github.com/SpoofIMEI/goramq@latest
```

## API Docs
```
API endpoint: `/search/`

Parameters:
  q: The query
    Example values:
      - badperson@gmail.com --> Searches for a line that contains single string ("badperson@gmail.com")
      - peter::@randomcompany.com --> Searches for a line that contains both strings, "peter" and "@randomcompany.com"
 
  caseins: Make the query case insensitive
 
  amount: How many results you want to receive (faster and doesn't require as much bandwith)
    Example value:
      - 10 --> Returns first 10 results

Example requests:
  - http://127.0.0.1:9112/search?q=randomguy::Estonia::2003&caseins=true
  - http://127.0.0.1:9112/search?q=username@protonmail.com&pass=verysecretAPIpassword
```

## Basic usage examples
```
#Load a file and start the web server in local only mode
goramq --filenames "database.sql"

#Loads 2 files and let the password protected API listen on all interfaces on port 80
goramq --listener :80 --filenames "documents/database.sql::documents/randomtextfile.txt" --password APIsecret

#Loads a file and does garbage collection only after all files have been loaded
goramq --filenames "database.sql" --pcgo 
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
  -filenames string
        Files to serve (separated with ::)
  -listener string
        Web server listener (IP:PORT) (default "127.0.0.1:9112")
  -nogc
        Disables garbage collection (uses twice as much ram)
  -password string
        Password protects the API
  -pcgo
        Only does garbage collection after all files are loaded.
  -quiet
        When used, Goramq will not display the file loading progress (perfomance boost)
  -result-limit int
        How many results you can receive (over limit = error message) (default 50000)
```
