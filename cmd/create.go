package cmd

import (
	"errors"
	"fmt"
	"log"

	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/bitrise-tools/replica/macosinstaller"
	"github.com/bitrise-tools/replica/vagrantbox"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create INSTALL_MACOS_APP_PATH",
	Short: `Create a vagrant box from an "Install macOS / OS X .." app`,
	Long:  `Create a vagrant box from an "Install macOS / OS X .." app`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("No 'Install macOS / OS X .. app' path provided")
		}
		installMacOSAppPath := args[0]
		return createVagrantBoxFromInstallMacOSApp(installMacOSAppPath)
	},
}

func init() {
	RootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func createVagrantBoxFromInstallMacOSApp(installMacOSAppPath string) error {
	if err := printToolVersions(); err != nil {
		return fmt.Errorf("Failed to print tool versions - missing tool - error: %s", err)
	}

	log.Printf("installMacOSAppPath: %s", installMacOSAppPath)

	macOSInstallDMGPath, err := macosinstaller.CreateInstallDMGFromInstallMacOSApp(installMacOSAppPath)
	if err != nil {
		return fmt.Errorf("Failed to create Install DMG, error: %s", err)
	}
	fmt.Println()
	log.Println(colorstring.Green("Done. Built image is located at " + macOSInstallDMGPath + "."))
	fmt.Println()

	fmt.Println()

	if isInstall, err := goinp.AskForBoolWithDefault("Do you want to create a vagrant box using the installer?", true); err != nil {
		return fmt.Errorf("Invalid input, error: %s", err)
	} else if !isInstall {
		return nil
	}
	fmt.Println()

	fmt.Println()
	log.Println(colorstring.Green(" => Creating vagrant box ..."))
	fmt.Println()

	if err := vagrantbox.CreateVirtualboxVagrantBoxFromPreparedMacOSInstallDMG(macOSInstallDMGPath); err != nil {
		return fmt.Errorf("Failed to create vagrant box, error: %s", err)
	}

	return nil
}
