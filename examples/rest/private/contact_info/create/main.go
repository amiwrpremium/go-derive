// Registers a new contact-info record on the configured signer's
// wallet. Required env: DERIVE_CONTACT_TYPE (e.g. "email") and
// DERIVE_CONTACT_VALUE.
package main

import (
	"log"
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	contactType := os.Getenv("DERIVE_CONTACT_TYPE")
	contactValue := os.Getenv("DERIVE_CONTACT_VALUE")
	if contactType == "" || contactValue == "" {
		log.Fatal("DERIVE_CONTACT_TYPE and DERIVE_CONTACT_VALUE required")
	}
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	contact, err := c.CreateContactInfo(ctx, contactType, contactValue)
	example.Fatal(err)
	example.Print("contact_id", contact.ID)
	example.Print("contact_type", contact.ContactType)
	example.Print("contact_value", contact.ContactValue)
}
