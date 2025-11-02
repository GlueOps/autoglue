package cmd

import (
	"fmt"
	"time"

	"github.com/glueops/autoglue/internal/app"
	"github.com/glueops/autoglue/internal/keys"
	"github.com/spf13/cobra"
)

var (
	alg     string
	rsaBits int
	kidFlag string
	nbfStr  string
	expStr  string
)

var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage JWT signing keys",
}

var keysGenCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate and store a new signing key",
	RunE: func(_ *cobra.Command, _ []string) error {
		rt := app.NewRuntime()

		var nbfPtr, expPtr *time.Time
		if nbfStr != "" {
			t, err := time.Parse(time.RFC3339, nbfStr)
			if err != nil {
				return err
			}
			nbfPtr = &t
		}
		if expStr != "" {
			t, err := time.Parse(time.RFC3339, expStr)
			if err != nil {
				return err
			}
			expPtr = &t
		}

		rec, err := keys.GenerateAndStore(rt.DB, rt.Cfg.JWTPrivateEncKey, keys.GenOpts{
			Alg:  alg,
			Bits: rsaBits,
			KID:  kidFlag,
			NBF:  nbfPtr,
			EXP:  expPtr,
		})
		if err != nil {
			return err
		}

		fmt.Printf("created signing key\n")
		fmt.Printf("  kid: %s\n", rec.Kid)
		fmt.Printf("  alg: %s\n", rec.Alg)
		fmt.Printf("  active: %v\n", rec.IsActive)
		if rec.NotBefore != nil {
			fmt.Printf("  nbf: %s\n", rec.NotBefore.Format(time.RFC3339))
		}
		if rec.ExpiresAt != nil {
			fmt.Printf("  exp: %s\n", rec.ExpiresAt.Format(time.RFC3339))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(keysCmd)
	keysCmd.AddCommand(keysGenCmd)

	keysGenCmd.Flags().StringVarP(&alg, "alg", "a", "EdDSA", "Signing alg: EdDSA|RS256|RS384|RS512")
	keysGenCmd.Flags().IntVarP(&rsaBits, "bits", "b", 3072, "RSA key size (when alg is RS*)")
	keysGenCmd.Flags().StringVarP(&kidFlag, "kid", "k", "", "Key ID (optional; auto if empty)")
	keysGenCmd.Flags().StringVarP(&nbfStr, "nbf", "n", "", "Not Before (RFC3339)")
	keysGenCmd.Flags().StringVarP(&expStr, "exp", "e", "", "Expires At (RFC3339)")
}
