package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func printJSON(value interface{}) int {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fail(err)
	}
	fmt.Println(string(data))
	return 0
}

func flagSet(name string) *flag.FlagSet {
	return flag.NewFlagSet(name, flag.ContinueOnError)
}

func printText(text string) int {
	fmt.Print(text)
	return 0
}

func readJSONFile(path string, target interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func printSummary(stressPa float64) int {
	fmt.Printf("sigma_z=%.6g Pa (%.6g kPa)\n", stressPa, stressPa/1000)
	return 0
}
