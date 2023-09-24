package command

import (
	"github.com/google/wire"
	"github.com/spf13/cobra"
	"go-web-wire-starter/internal/dao"
)

var ProviderSet = wire.NewSet(dao.NewDB, NewDataBaseCommand, NewCommand)

type Command struct {
	dbCmd *DataBaseCommand
}

// NewCommand .
func NewCommand(
	dbCmd *DataBaseCommand,
) *Command {
	return &Command{
		dbCmd: dbCmd,
	}
}

// Register 注册子命令
func Register(rootCmd *cobra.Command, newCmd func() (*Command, func(), error)) {
	rootCmd.AddCommand(
		&cobra.Command{
			Use:   "migrate",
			Short: "数据库迁移",
			Run: func(cmd *cobra.Command, args []string) {
				command, cleanup, err := newCmd()
				if err != nil {
					panic(err)
				}
				defer cleanup()

				command.dbCmd.Migrate(cmd, args)
			},
		},
	)
}
