// Lists every contact-info record on the configured signer's wallet.
// Optional env: DERIVE_CONTACT_TYPE to filter by type.
package main

import (
	"os"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	contacts, err := c.GetContactInfo(ctx, os.Getenv("DERIVE_CONTACT_TYPE"))
	example.Fatal(err)
	example.Print("contacts", len(contacts))
	for _, ct := range contacts {
		example.Print("contact", ct.ID)
		example.Print("  type", ct.ContactType)
		example.Print("  value", ct.ContactValue)
	}
}
