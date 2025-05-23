package main

/*

$> echo "hello @bob@example.com pass the mustard to @alice@mustard.com and doug@localhost before sending it @doug@bob.com and max@gmail.com" | ./bin/list-addresses -
@bob@example.com
@alice@mustard.com
@doug@bob.com

*/

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sfomuseum/go-activitypub/ap"
)

func main() {

	flag.Parse()

	body := strings.Join(flag.Args(), " ")

	if body == "-" {

		body = ""
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			line := scanner.Text()
			body = fmt.Sprintf("%s %s", body, line)
		}

		err := scanner.Err()

		if err != nil {
			log.Fatalf("Failed to read data, %v", err)
		}
	}

	addrs, err := ap.ParseAddressesFromString(body)

	if err != nil {
		log.Fatalf("Failed to parse addresses, %v", err)
	}

	for _, a := range addrs {
		fmt.Println(a)
	}
}
