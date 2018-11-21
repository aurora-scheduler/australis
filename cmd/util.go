package cmd

import (
    "bytes"
    "encoding/json"
    "fmt"
	"github.com/paypal/gorealis/gen-go/apache/aurora"

	log "github.com/sirupsen/logrus"
)

func toJSON(v interface{}) string {

    output, err := json.Marshal(v)

    if err != nil {
        log.Fatalln("Unable to serialize Aurora response: %+v", v)
    }

    return string(output)
}


func getLoggingLevels() string {

    var buffer bytes.Buffer

    for _, level := range log.AllLevels {
        buffer.WriteString(level.String())
        buffer.WriteString(" ")
    }

    buffer.Truncate(buffer.Len()-1)

    return buffer.String()

}


func maintenanceMonitorPrint(hostResult map[string]bool, desiredStates []aurora.MaintenanceMode) {
	if len(hostResult) > 0 {
		// Create anonymous struct for JSON formatting
		output := struct{
			DesiredStates []string `json:desired_states`
			Transitioned []string `json:transitioned`
			NonTransitioned []string `json:non-transitioned`
		}{
			make([]string, 0),
			make([]string, 0),
			make([]string, 0),
		}

		for _,state := range desiredStates {
			output.DesiredStates = append(output.DesiredStates, state.String())
		}

		for host, ok := range hostResult {
			if ok {
				output.Transitioned = append(output.Transitioned, host)
			} else {
				output.NonTransitioned = append(output.NonTransitioned, host)
			}
		}

		if toJson {
			fmt.Println(toJSON(output))
		} else {
			fmt.Printf("Entered %v status: %v\n", output.DesiredStates, output.Transitioned)
			fmt.Printf("Did not enter %v status: %v\n", output.DesiredStates, output.NonTransitioned)
		}
	}
}
