package iwallet

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/iost-official/go-iost/iwallet/contract"
	"github.com/spf13/cobra"
)

// Generate ABI file.
func generateABI(codePath string) (string, error) {
	contractToRun := fmt.Sprintf(`
	const fs = require('fs');
	%s;
	const processContract = module.exports;
	processContract("%s");
	`, contract.CompiledContract, codePath)

	cmd := exec.Command("node", "-")
	cmd.Stdin = strings.NewReader(contractToRun)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to compile, error: %v\n", err)
		fmt.Println(string(output))
		fmt.Printf("Please make sure node.js has been installed\n")
		return "", err
	}

	return codePath + ".abi", nil
}

// compileCmd represents the compile command.
var compileCmd = &cobra.Command{
	Use:     "compile codePath",
	Short:   "Generate contract abi",
	Long:    `Generate abi from contract javascript code`,
	Example: `  iwallet compile ./example.js`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := checkArgsNumber(cmd, args, "codePath"); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		codePath := args[0]
		abiPath, err := generateABI(codePath)
		if err != nil {
			return fmt.Errorf("failed to generate abi: %v", err)
		}
		fmt.Printf("Successfully generated abi file as: %v\n", abiPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
}
