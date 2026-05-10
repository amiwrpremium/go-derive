// Fetches the headline referral performance for the configured
// referral code over the last 30 days. Set DERIVE_REFERRAL_CODE to
// scope to one code; otherwise the engine returns the caller's
// own referral metrics if any.
package main

import (
	"os"
	"time"

	"github.com/amiwrpremium/go-derive/examples/example"
)

func main() {
	c := example.MustRESTPublic()
	defer c.Close()
	ctx, cancel := example.Timeout()
	defer cancel()

	end := time.Now().UnixMilli()
	start := end - int64(30*24*time.Hour/time.Millisecond)
	params := map[string]any{
		"start_ms": start,
		"end_ms":   end,
	}
	if code := os.Getenv("DERIVE_REFERRAL_CODE"); code != "" {
		params["referral_code"] = code
	}

	res, err := c.GetReferralPerformance(ctx, params)
	example.Fatal(err)
	example.Print("referral_code", res.ReferralCode)
	example.Print("fee_share_percentage", res.FeeSharePercentage.String())
	example.Print("stdrv_balance", res.StdrvBalance.String())
	example.Print("total_notional_volume", res.TotalNotionalVolume.String())
	example.Print("total_referred_fees", res.TotalReferredFees.String())
	example.Print("total_fee_rewards", res.TotalFeeRewards.String())
}
