package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
	//"regexp"
)

/**global vars**/
var (
	userMap   = make(map[string]int)
	debug     bool
	debugDate bool
	rawtxt    bool
	rawcsv    bool
	limit     string
)

/**map functions**/
func mapAddUnique(key string) {

	if _, ok := userMap[key]; ok {
		userMap[key] = userMap[key] + 1
	} else {
		userMap[key] = 1
	}

}

/***Date Functions**/
func DateWithinRange(dateString string) bool {
	dateFormat := "01/02/2006"
	isInRange := false
	dateComponents := strings.Split(dateString, "/")
	months := dateComponents[0]
	days := dateComponents[1]
	years := dateComponents[2]
	date:=dateString

	if len(months) == 1 {
		if len(days) == 1 {
			date = "0" + months + "/" + "0" + days + "/" + years
		} else {
			date = "0" + months + "/" + days + "/" + years
		}

	}else{
		if len(days) == 1 {
			date = months + "/" + "0" + days + "/" + years
		} else {
			date = months + "/" + days + "/" + years
		}
	}

	dateStamp, err := time.Parse(dateFormat, date)

	if err != nil {
		if debug {
			fmt.Println(err.Error())
		}
		isInRange = false
	}

	today := time.Now()

	twelveMonthsAgo := today.AddDate(0, -12, 0) // minus 12 months (1 year)

	if dateStamp.After(twelveMonthsAgo) {
		isInRange = true
	}
	if debugDate {
		fmt.Println(date, "isInRange: ", isInRange)
	}
	return isInRange
}

/**TF.Exe functions**/
func tf_History(stopafter string) string {

	cmd := exec.Command("tf", "history", "/collection:http://tfs.molina.mhc:8080/tfs/HSS/", "$/HSS", "/recursive", "/stopafter:"+stopafter, "/noprompt")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if debug {
			fmt.Println("Error")
		}
		return err.Error()
	}
	return string(out[:])
}

/***FCollect data***/
func collectData(data string) int {
	lines := strings.Split(data, "\n")
	totalLines := len(lines)
	count := 0
	changeset := ""
	userName := ""
	checkInDate := ""
	for _, line := range lines {

		if count != 1 && count < totalLines-1 {

			changeset, userName, checkInDate = formatData(line)

		}
		changeset = changeset
		count++
		//add to map and check date range here
		if userName != "" && userName != "User" {
			if debug {
				fmt.Println(userName + ": " + checkInDate)
			}
			isInDateRange := DateWithinRange(checkInDate)
			if isInDateRange {

				mapAddUnique(userName)
			}

		}

	}

	file, err := os.Create("./userlog.csv")
	if err != nil {
		err = err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	userCount := 0

	for k, v := range userMap {

		fmt.Fprintln(w, k, ",", v)

		if v >= 10 {
			k = k
			userCount++
		}
	}
	err = w.Flush()

	//fmt.Println("Number of users: ", userCount)
	return userCount
}

/**format Data**/
func formatData(line string) (string, string, string) {

	userInfo := strings.Fields(line)

	changeset := "ChangeSet"
	username := "User"
	date := "Date"
	changeset = userInfo[0]
	if 2 < len(userInfo) && (2 < len(userInfo[2]) && (userInfo[2][1] == '/' || userInfo[2][2] == '/')) {
		//Service account

		username = userInfo[1]
		username = strings.Replace(username, ",", "", -1)
		date = userInfo[2]
	} else if 3 < len(userInfo) && (2 < len(userInfo[3]) && (userInfo[3][1] == '/' || userInfo[3][2] == '/')) {
		//two names
		username = userInfo[1] + userInfo[2]
		username = strings.TrimSpace(username)
		username = strings.Replace(username, ",", "", -1)
		date = userInfo[3]
	} else if 4 < len(userInfo) && (2 < len(userInfo[4]) && (userInfo[4][1] == '/' || userInfo[4][2] == '/')) {
		//three names
		username = userInfo[1] + userInfo[2] + userInfo[3]
		username = strings.TrimSpace(username)
		username = strings.Replace(username, ",", "", -1)
		date = userInfo[4]
	}

	return changeset, username, date
}

/**write file **/
func writeToTxtRaw(data string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	lines := strings.Split(data, "\n")

	for _, line := range lines {

		fmt.Fprintln(w, line)

	}

	return w.Flush()
}

func writeToCSVRaw(data string, path string) error {
	changeset := ""
	username := ""
	date := ""
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	lines := strings.Split(data, "\n")
	totalLines := len(lines)

	count := 0

	nameCount := 0
	for _, line := range lines {
		if count != 1 && count < totalLines-1 {

			changeset, username, date = formatData(line)
			fileData := changeset + "," + username + "," + date
			fmt.Fprintln(w, fileData)
		}

		count++

	}
	if debug {
		/*fmt.Println(dateCount)*/
		fmt.Println(nameCount)
	}

	return w.Flush()
}

/***Main***/
func main() {

	flag.BoolVar(&debug, "debug", false, "enable debug printing")
	flag.BoolVar(&rawtxt, "txt", false, "create raw text file")
	flag.BoolVar(&rawcsv, "csv", false, "create raw csv file")
	flag.BoolVar(&debugDate, "debugDate", false, "debugDates")
	flag.StringVar(&limit, "limit", "100000", "Limit changeset amount")
	flag.Parse()
	fmt.Println("Starting Active User Count...")
	history := tf_History(limit)
	if rawtxt {
		writeToTxtRaw(history, "./raw.txt")
	}
	if rawcsv {
		writeToCSVRaw(history, "./raw.csv")
	}
	activeUserCount := collectData(history)
	fmt.Println("Active User Count: ", activeUserCount)
}
