package main

import (
	"fmt"
	"strings"
	"testing"
)

/***Map Functions***/
func TestAddMap(t *testing.T) {
	debug = true
	mapAddUnique("user1")
	mapAddUnique("user2")
	mapAddUnique("user2")
	mapAddUnique("user2")
	mapAddUnique("user2")

	fmt.Println("user1 count: ", userMap["user1"])
	fmt.Println("user2 count: ", userMap["user2"])
}

/**format data**/
func TestFormatData(t *testing.T) {
	debug = true
	debugDate=true
	fmt.Println("FormatData Test")

	changeset := ""
	user := ""
	checkin := ""

	fmt.Println("one Name")
	input := "211996    Gnanadhiraviam... 11/9/2018  Updated"
	changeset, user, checkin = formatData(input)
	fmt.Println(changeset, user, checkin)

	fmt.Println("Service Account")
	input = "211995    TFSSERVICE        11/9/2018  Updating file"
	changeset, user, checkin = formatData(input)
	fmt.Println(changeset, user, checkin)

	fmt.Println("two names")
	input = "211995    luis,	albert       11/9/2018  Updating file"
	changeset, user, checkin = formatData(input)
	fmt.Println(changeset, user, checkin)

	fmt.Println("three names")
	input = "211995    luis,	albert One       11/9/2018  Updating file"
	changeset, user, checkin = formatData(input)
	fmt.Println(changeset, user, checkin)

}

/*****Date Functions *****/
func TestGenerateDateRange1(t *testing.T) {
	debug = false
	date := "2/1/2018"
	isInRange := DateWithinRange(date)
	fmt.Println(date, " After funct Is in Range true: ", isInRange)
}
func TestGenerateDateRange2(t *testing.T) {
	debug = true
	date := "11/9/2017"
	isInRange := DateWithinRange(date)
	fmt.Println(date, "After funct Is in Range false: ", isInRange)
}
func TestGenerateDateRange3(t *testing.T) {
	debug = true
	date := "2/12/2018"
	isInRange := DateWithinRange(date)
	fmt.Println(date, "After funct Is in Range true: ", isInRange)
}
func TestGenerateDateRange4(t *testing.T) {
	debug = true
	date := "10/12/2015"
	isInRange := DateWithinRange(date)
	fmt.Println(date, "After funct Is in Range false: ", isInRange)
}
func TestGenerateDateRangeFalse(t *testing.T) {
	debug = true
	isInRange := DateWithinRange("02/01/2017")
	fmt.Println(isInRange)
}

func TestGenerateDateRangeTrue(t *testing.T) {
	debug = true
	isInRange := DateWithinRange("02/01/2018")
	fmt.Println(isInRange)
}

/*func TestWriteTextCSV(t *testing.T) {
	debug = true
	data := tf_History("60000")
	writeToTxtRaw(data, "./raw.txt")

}*/
func TestWriteRawCSV(t *testing.T) {
	debug = true
	data := tf_History("10")
	writeToCSVRaw(data, "./raw.csv")
}

/*func TestCollectData(t *testing.T) {
	debug = true
	userMap = nil
	userMap = make(map[string]int)
	data := tf_History("1")
	collectData(data)
}*/

/**TF EXE functions**/
func TestTFHistory(t *testing.T) {
	debug = true
	stopAfter := "10"

	output := tf_History(stopAfter)
	temp := strings.Split(output, "\n")
	fmt.Println(temp[1])
	fmt.Println(temp[2])
	userinfo := strings.Fields(temp[2])

	fmt.Println(userinfo)

	username := userinfo[1]
	userlastname := userinfo[2]
	checkindate := userinfo[3]

	fmt.Println("Username: ", username+" "+userlastname)
	fmt.Println("checkinDate: ", checkindate)

	//fmt.Println(userinfo[0])
	//fmt.Println(userinfo[1])
	//fmt.Println(userinfo[2])
	//fmt.Println(userinfo[3])
	//fmt.Println(userinfo[4])

}
