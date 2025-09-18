package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	var db = initDB()

	//Uncomment below lines to insert CSV data into MySQL

	// err := InsertCSVToMySQL(db, "cities.csv", "cities")
	// if err != nil {
	// 	fmt.Printf("Error inserting CSV to MySQL: %v\n", err)
	// }
	distributors, err := loadDistributors(db)
	if err != nil {
		fmt.Errorf("Failed to load distributors:", err)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter distributor name (or type 'exit' to quit): ")
		distInput, _ := reader.ReadString('\n')
		distInput = strings.TrimSpace(distInput)
		if distInput == "exit" {
			fmt.Println("Exiting...")
			break
		}

		dist, ok := distributors[distInput]
		if !ok {
			fmt.Println("Unknown distributor. Please try again.")
			continue
		}

		fmt.Print("Enter region code to check permissions (e.g., CITY-PROVINCE-COUNTRY): ")
		regionInput, _ := reader.ReadString('\n')
		regionInput = strings.TrimSpace(regionInput)
		if regionInput == "" {
			fmt.Println("Region cannot be empty. Please try again.")
			continue
		}

		if dist.Permissions.CanDistribute(regionInput) {
			fmt.Printf("YES, %s can distribute in %s\n", distInput, regionInput)
		} else {
			fmt.Printf("NO, %s cannot distribute in %s\n", distInput, regionInput)
		}
		fmt.Println()
	}
	testCases := []struct {
		distributor string
		region      string
	}{
		{"DISTRIBUTOR1", "CHICAGO-ILLINOIS-UNITEDSTATES"},
		{"DISTRIBUTOR1", "CHENNAI-TAMILNADU-INDIA"},
		{"DISTRIBUTOR1", "BANGALORE-KARNATAKA-INDIA"},
		{"DISTRIBUTOR2", "BANGALORE-KARNATAKA-INDIA"},
		{"DISTRIBUTOR2", "MUMBAI-MAHARASHTRA-INDIA"},
		{"DISTRIBUTOR3", "HUBLI-KARNATAKA-INDIA"},
	}

	for _, tc := range testCases {
		dist, ok := distributors[tc.distributor]
		if !ok {
			fmt.Printf("Distributor %s not found\n", tc.distributor)
			continue
		}
		canDist := dist.Permissions.CanDistribute(tc.region)
		fmt.Printf("%s can distribute in %s? %v\n", tc.distributor, tc.region, canDist)
	}
}
