package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var siteCmd = &cobra.Command{
	Use:   "site",
	Short: "Creates a new site in your home directory",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("site: please provide the site name")
			return
		}

		home, err := homedir.Dir()
		if err != nil {
			fmt.Println("site: could not detect your home directory")
			return
		}

		dirname := fmt.Sprintf("%s/%s", home, args[0])
		err = os.Mkdir(dirname, 0755)
		if err != nil {
			fmt.Println("site: failed to create site directory")
			return
		}

		if err := setupSkeleton(dirname); err != nil {
			fmt.Println("site: failed to create skeleton for new site")
			return
		}
	},
}

func setupSkeleton(path string) error {
	contentsPath := fmt.Sprintf("%s/contents", path)
	configPath := fmt.Sprintf("%s/config.toml", path)

	err := os.Mkdir(contentsPath, 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("title = \"New Site\"")
	if err != nil {
		return err
	}

	return nil
}
