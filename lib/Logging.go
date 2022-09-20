package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"text/template"
	"time"
)

// LogEntryObject should be able to log Http requests from the server
// and also log the errors that occur in the server
// it should have:
// 1. the method of the request
// 2. the uri of the request
// 3. the status code of the request
// 4. the time it took to process the request
// 5. an error message if there is one
type LogEntryObject struct {
	Method     string       // GET, POST, PUT, DELETE
	URI        string       // in from of /path/to/endpoint
	StatusCode int          // 200, 404, 500, etc
	DateTime   string       // in the format of 2006-01-02 15:04:05
	TimeTaken  int64        // in milliseconds
	Size       float32      // in kilobytes
	Reason     string       // a message to explain the status code
	Error      error        // nil if there is no error
	Message    string       // SUCCESS, FAILURE, ERROR
	Code       LogEntryEnum // 1000, 2000, 3000
}

type LogEntries []LogEntryObject

// LogEntryEnum is an enum that defines the different types of log entries
// there are 3 types of log entries:
// 1. a Success log entry
// 2. a Failure log entry
// 3. an Error log entry
type LogEntryEnum int

const (
	LogSuccess LogEntryEnum = 1000
	LogFailure LogEntryEnum = 2000
	LogError   LogEntryEnum = 3000
)

func New(r *http.Request) LogEntryObject {
	return LogEntryObject{
		Method:    r.Method,
		URI:       r.RequestURI,
		DateTime:  time.Now().Format("RFC3339"),
		TimeTaken: time.Now().UnixNano(),
	}
}

func (ls LogEntries) StartLogging() {
	// check the operating system and load the correct file.
	var logpath string
	if runtime.GOOS == "windows" {
		logpath, _ = filepath.Abs(".\\lib\\log\\tavern.log")
	} else {
		logpath, _ = filepath.Abs("./lib/log/tavern.log")
	}

	// delete the file if it exists
	if _, err := os.Stat(logpath); err == nil {
		os.Remove(logpath)
	}
	// create the file
	f, err := os.Create(logpath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
}

func (ls LogEntries) RenderHtml(w http.ResponseWriter) {
	// TODO: IMPLEMENT THIS IN THE /api/admin/log endpoint
	// also TODO: make a groups i.e.
	// 1. User
	// 2. Premium
	// 3. Admin

	// check the operating system and load the correct file.
	var logpath string
	if runtime.GOOS == "windows" {
		logpath, _ = filepath.Abs(".\\lib\\log\\tavern.log")
	} else {
		logpath, _ = filepath.Abs("./lib/log/tavern.log")
	}

	// load the file,
	f, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	// read the file
	// first try to load the file into the LogEntries struct
	// if the file is empty, then panic
	fstring, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	if len(fstring) > 0 {
		json.Unmarshal(fstring, &ls)
	} else {
		panic("The log file is empty")
	}
	defer f.Close()

	// create the html file
	var htmlpath string
	if runtime.GOOS == "windows" {
		htmlpath, _ = filepath.Abs(".\\lib\\html\\AdminLogData.html")
	} else {
		htmlpath, _ = filepath.Abs("./lib/html/AdminLogData.html")
	}
	tmpl, err := template.ParseFiles(htmlpath)
	if err != nil {
		panic(err)
	}

	tmpl.Execute(w, ls)
}

func (logEntry LogEntryObject) Log(
	size float32,
	statusCode int,
	reason string,
	er ...error,
) {
	ls := LogEntries{}

	i, err := strconv.ParseInt(fmt.Sprintf("%d", logEntry.TimeTaken), 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	logEntry.TimeTaken = time.Since(tm).Milliseconds()
	logEntry.Size = size
	logEntry.StatusCode = statusCode

	if er != nil {
		logEntry.Code = LogError
		logEntry.Message = "error"
	} else if logEntry.StatusCode >= 200 && logEntry.StatusCode < 300 {
		logEntry.Code = LogSuccess
		logEntry.Message = "success"
	} else {
		logEntry.Code = LogFailure
		logEntry.Message = "failure"
	}

	// log to file
	var logpath string
	if runtime.GOOS == "windows" {
		logpath, _ = filepath.Abs(".\\lib\\log\\tavern.log")
	} else {
		logpath, _ = filepath.Abs("./lib/log/tavern.log")
	}
	f, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	// read the file
	// first try to load the file into the LogEntries struct
	// if the file is empty, then just append the new log entry
	fstring, err := io.ReadAll(f)
	if err != nil {
		fmt.Println(err)
	}

	if len(fstring) > 0 {
		err = json.Unmarshal(fstring, &ls)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		ls = append(ls, logEntry)
	}

	j, err := json.Marshal(ls)
	if err != nil {
		fmt.Println(err)
	}
	if _, err := f.Write(j); err != nil {
		fmt.Println(err)
	}
	defer f.Close()
}
