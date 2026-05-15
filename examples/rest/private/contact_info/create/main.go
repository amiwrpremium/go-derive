// Registers a new contact-info record on the configured signer's
// wallet. Required env: DERIVE_CONTACT_TYPE (e.g. "email") and
// DERIVE_CONTACT_VALUE.
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
	"github.com/amiwrpremium/go-derive/pkg/types"
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
	contactType := os.Getenv("DERIVE_CONTACT_TYPE")
	contactValue := os.Getenv("DERIVE_CONTACT_VALUE")
	if contactType == "" || contactValue == "" {
		log.Fatal("DERIVE_CONTACT_TYPE and DERIVE_CONTACT_VALUE required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	contact, err := c.CreateContactInfo(ctx, types.CreateContactInfoInput{ContactType: contactType, ContactValue: contactValue})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "contact_id:", contact.ID)
	fmt.Printf("%-30s %v\n", "contact_type:", contact.ContactType)
	fmt.Printf("%-30s %v\n", "contact_value:", contact.ContactValue)
}
