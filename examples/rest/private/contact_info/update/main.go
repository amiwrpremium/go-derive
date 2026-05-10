// Updates the value of an existing contact-info record. Required
// env: DERIVE_CONTACT_ID, DERIVE_CONTACT_VALUE.
package main

import (
	"log"
	"os"
	"strconv"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	idStr := os.Getenv("DERIVE_CONTACT_ID")
	value := os.Getenv("DERIVE_CONTACT_VALUE")
	if idStr == "" || value == "" {
		log.Fatal("DERIVE_CONTACT_ID and DERIVE_CONTACT_VALUE required")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Fatalf("DERIVE_CONTACT_ID=%q: %v", idStr, err)
	}
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	contact, err := c.UpdateContactInfo(ctx, id, value)
	example.Fatal(err)
	example.Print("contact_id", contact.ID)
	example.Print("new value", contact.ContactValue)
	example.Print("updated_at_sec", contact.UpdatedAtSec)
}
