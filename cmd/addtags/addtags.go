package addtags

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sagan/ptool/client"
	"github.com/sagan/ptool/cmd"
)

var command = &cobra.Command{
	Use:   "addtags <client> <tags> [<infoHash>...]",
	Short: "Add tags to torrents in client.",
	Long: `Add tags to torrents in client.
<tags> : comma-seperated tags list.
<infoHash>...: infoHash list of torrents. It's possible to use state filter to target multiple torrents:
_all, _active, _done, _undone, _downloading, _seeding, _paused, _completed, _error.`,
	Args: cobra.MatchAll(cobra.MinimumNArgs(2), cobra.OnlyValidArgs),
	RunE: addtags,
}

var (
	category = ""
	tag      = ""
	filter   = ""
)

func init() {
	command.Flags().StringVarP(&filter, "filter", "f", "", "Filter torrents by name")
	command.Flags().StringVarP(&category, "category", "c", "", "Filter torrents by category")
	command.Flags().StringVarP(&tag, "tag", "t", "", "Filter torrents by tag. Comma-separated string list. Torrent which tags contain any one in the list will match")
	cmd.RootCmd.AddCommand(command)
}

func addtags(cmd *cobra.Command, args []string) error {
	clientName := args[0]
	tags := strings.Split(args[1], ",")
	args = args[2:]
	if category == "" && tag == "" && filter == "" && len(args) == 0 {
		return fmt.Errorf("you must provide at least a condition flag or hashFilter")
	}
	clientInstance, err := client.CreateClient(clientName)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}

	infoHashes, err := client.SelectTorrents(clientInstance, category, tag, filter, args...)
	if err != nil {
		clientInstance.Close()
		return err
	}
	if infoHashes == nil {
		err = clientInstance.AddTagsToAllTorrents(tags)
		if err != nil {
			clientInstance.Close()
			return err
		}
	} else if len(infoHashes) > 0 {
		err = clientInstance.AddTagsToTorrents(infoHashes, tags)
		if err != nil {
			clientInstance.Close()
			return err
		}
	}
	clientInstance.Close()
	return nil
}
