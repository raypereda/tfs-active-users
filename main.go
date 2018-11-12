// tfs-active-users is a utiltity that counts the number of active users in the last year.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

var collection string
var minCheckins int

func main() {
	flag.StringVar(&collection, "collection", "http://tfs.molina.mhc:8080/tfs/HSS/", "TFS collection URL")
	flag.IntVar(&minCheckins, "mincheckins", 10, "minimum number of check-ins for active users")
	flag.Parse()

	h := history()
	c := countByUser(h)

	print(c)
}

func path() string {
	parts := strings.Split(collection, "/")
	collection := parts[len(parts)-2]
	rootPath := "$/" + collection
	return rootPath
}

func dateRange() string {
	dateLayout := "2006-01-02"
	now := time.Now()
	nowDate := now.Format(dateLayout)
	yearAgo := now.AddDate(-1, 0, 0)
	yearAgoDate := yearAgo.Format(dateLayout)
	versionSpec := "/version:D%s~D%s"
	spec := fmt.Sprintf(versionSpec, yearAgoDate, nowDate)
	return spec
}

func history() []string {
	cmdline := "tf history /collection:" + collection +
		" " + path() + " /recursive " + dateRange() + " /noprompt"
	fmt.Println("executing the command line below")
	fmt.Println(cmdline)
	fmt.Println()

	// example output of tfs history command:
	// Changeset User              Date       Comment
	// --------- ----------------- ---------- ----------------------------------------
	// 30        Raisa Pokrovskaya 4/23/2012
	// 29        Jamal Hartnett    4/23/2012  Fix bug in new method
	// 20        Raisa Pokrovskaya 4/12/2012  Add new method, add program2.cs to

	cmd := exec.Command("tf", "history",
		"/collection:"+collection,
		path(),
		"/recursive",
		dateRange(),
		"/noprompt")
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "tf.exe error: %s \n", err)
		os.Exit(1)
	}
	lines := strings.Split(string(out), "\n")
	return lines
}

type stat struct {
	name  string
	count int
}

func countByUser(history []string) []stat {
	c := make(map[string]int)
	history = history[2:]
	for _, l := range history {
		if len(l) < 28 {
			continue
		}
		user := strings.TrimSpace(l[10:28])
		c[user]++
	}

	var list []stat
	for user, count := range c {
		if count >= minCheckins {
			s := stat{user, count}
			list = append(list, s)
		}
	}
	isGreater := func(i, j int) bool {
		return list[i].count > list[j].count ||
			list[i].count == list[j].count && list[i].name < list[j].name
	}
	sort.Slice(list, isGreater)
	return list
}

func print(s []stat) {
	fmt.Printf("%d active users with at least %d check-ins\n\n", len(s), minCheckins)
	fmt.Printf("Check-ins      User\n")
	fmt.Printf("--------- -----------------\n")

	for _, line := range s {
		fmt.Printf("   %3d    %s\n", line.count, line.name)
	}
}
