package main

/*

$> echo "hello @bob@example.com pass the mustard to @alice@mustard.com and doug@localhost before sending it @doug@bob.com and max@gmail.com" | ./bin/list-addresses -
@bob@example.com
@alice@mustard.com
@doug@bob.com

*/

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"bufio"
	"os"
	
	"github.com/sfomuseum/go-activitypub"
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
	
	addrs, err := activitypub.ParseAddressesFromString(body)

	if err != nil {
		log.Fatalf("Failed to parse addresses, %v", err)
	}

	for _, a := range addrs {
		fmt.Println(a)
	}
}
