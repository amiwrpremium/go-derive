// Updates the value of an existing contact-info record. Required
// env: DERIVE_CONTACT_ID, DERIVE_CONTACT_VALUE.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/amiwrpremium/go-derive/pkg/auth"
	"github.com/amiwrpremium/go-derive/pkg/rest"
)

func main() {
	subStr := os.Getenv("DERIVE_SUBACCOUNT")
	if subStr == "" {
		log.Fatal("DERIVE_SUBACCOUNT required")
	}
	subaccount, err := strconv.ParseInt(subStr, 10, 64)
	if err != nil {
		log.Fatalf("DERIVE_SUBACCOUNT=%q: %v", subStr, err)
	}
	key := os.Getenv("DERIVE_SESSION_KEY")
	if key == "" {
		log.Fatal("DERIVE_SESSION_KEY required")
	}
	var signer auth.Signer
	if owner := os.Getenv("DERIVE_OWNER"); owner != "" {
		signer, err = auth.NewSessionKeySigner(key, common.HexToAddress(owner))
	} else {
		signer, err = auth.NewLocalSigner(key)
	}
	if err != nil {
		log.Fatalf("signer: %v", err)
	}

	restNetwork := rest.WithTestnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		restNetwork = rest.WithMainnet()
	}
	c, err := rest.New(restNetwork, rest.WithSigner(signer), rest.WithSubaccount(subaccount))
	if err != nil {
		log.Fatalf("rest.New: %v", err)
	}
	defer c.Close()
	idStr := os.Getenv("DERIVE_CONTACT_ID")
	value := os.Getenv("DERIVE_CONTACT_VALUE")
	if idStr == "" || value == "" {
		log.Fatal("DERIVE_CONTACT_ID and DERIVE_CONTACT_VALUE required")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Fatalf("DERIVE_CONTACT_ID=%q: %v", idStr, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	contact, err := c.UpdateContactInfo(ctx, id, value)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "contact_id:", contact.ID)
	fmt.Printf("%-30s %v\n", "new value:", contact.ContactValue)
	fmt.Printf("%-30s %v\n", "updated_at_sec:", contact.UpdatedAtSec)
}
