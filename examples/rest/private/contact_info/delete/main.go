// Deletes one contact-info record by id. Required env:
// DERIVE_CONTACT_ID.
package main

import (
	"log"
	"os"
	"strconv"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	idStr := os.Getenv("DERIVE_CONTACT_ID")
	if idStr == "" {
		log.Fatal("DERIVE_CONTACT_ID required")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Fatalf("DERIVE_CONTACT_ID=%q: %v", idStr, err)
	}
	c := example.MustRESTPrivate()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	deletedID, deleted, err := c.DeleteContactInfo(ctx, id)
	example.Fatal(err)
	example.Print("contact_id", deletedID)
	example.Print("deleted", deleted)
}
