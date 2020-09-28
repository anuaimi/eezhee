package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(teardownCmd)
}

var teardownCmd = &cobra.Command{
	Use:   "teardown",
	Short: "Will teardown a running app",
	Long:  `All software has versions. This is Eezhee's`,
	Run: func(cmd *cobra.Command, args []string) {

		teardownVM()
	},
}

func teardownVM() {
	// see if there is a state file so we know what we're supposed to teardown

	// read deploy.state
	viper.SetConfigName("deploy-state")
	viper.SetConfigFile("./deploy-state.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("nothing to teardown as no state file found")
		} else {
			fmt.Println("error reading state file")
		}
		return
	}

	cloud := viper.GetString("cloud")
	if strings.Compare(cloud, "digitalocean") != 0 {
		fmt.Println("state file reference cloud we don't support: ", cloud)
		return
	}

	// get details of VM
	ID := viper.GetInt("id")
	if ID == 0 {
		fmt.Println("invalid VM ID:", ID, ".  Can no teardown VM")
		return
	}

	// doctl compute droplet delete id
	cmd := exec.Command("doctl", "compute", "droplet", "delete", strconv.Itoa(ID), "--force", "-o", "json")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(stdoutStderr))
		fmt.Println(err)
		return
	}

	fmt.Println(string(stdoutStderr))

	// remove deploy.yaml
	err = os.Remove("./deploy-state.yaml")
	if err != nil {
		fmt.Println("could not remove deploy-state.yaml file")
		fmt.Println(err)
	}

}