package mockingjay

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"
)

func addNewLinesToConfig(y []byte) []byte {

	var newConfig bytes.Buffer
	resultWriter := bufio.NewWriter(&newConfig)

	scanner := bufio.NewScanner(bytes.NewReader(y))

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "- ") {
			fmt.Fprintln(resultWriter, "")
		}
		fmt.Fprintln(resultWriter, scanner.Text())
	}

	err := resultWriter.Flush()

	if err != nil {
		log.Println("Oopsie, problem adding new lines to config", err)
	}

	return newConfig.Bytes()
}
