// Fetches the headline referral performance for the configured
// referral code over the last 30 days. Set DERIVE_REFERRAL_CODE to
// scope to one code; otherwise the engine returns the caller's
// own referral metrics if any.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/pkg/rest"
	"github.com/amiwrpremium/go-derive/pkg/types"
)

func main() {
	restNetwork := rest.WithTestnet()
	if os.Getenv("DERIVE_NETWORK") == "mainnet" {
		restNetwork = rest.WithMainnet()
	}
	c, err := rest.New(restNetwork)
	if err != nil {
		log.Fatalf("rest.New: %v", err)
	}
	defer c.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	end := time.Now().UnixMilli()
	start := end - int64(30*24*time.Hour/time.Millisecond)
	q := types.ReferralPerformanceQuery{
		StartMs:      start,
		EndMs:        end,
		ReferralCode: os.Getenv("DERIVE_REFERRAL_CODE"),
	}

	res, err := c.GetReferralPerformance(ctx, q)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%-30s %v\n", "referral_code:", res.ReferralCode)
	fmt.Printf("%-30s %v\n", "fee_share_percentage:", res.FeeSharePercentage.String())
	fmt.Printf("%-30s %v\n", "stdrv_balance:", res.StdrvBalance.String())
	fmt.Printf("%-30s %v\n", "total_notional_volume:", res.TotalNotionalVolume.String())
	fmt.Printf("%-30s %v\n", "total_referred_fees:", res.TotalReferredFees.String())
	fmt.Printf("%-30s %v\n", "total_fee_rewards:", res.TotalFeeRewards.String())
}
