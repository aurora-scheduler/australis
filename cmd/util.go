package cmd

import (
    "bytes"
    "encoding/json"

    log "github.com/sirupsen/logrus"
)

func toJSON(v interface{}) string {

    output, err := json.Marshal(v)

    if err != nil {
        log.Fatalln("Unable to serialize statuses")
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