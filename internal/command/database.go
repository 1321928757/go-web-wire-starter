package command

import (
	"github.com/spf13/cobra"
	"go-web-wire-starter/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 数据库操作命令
type DataBaseCommand struct {
	logger *zap.Logger
	db     *gorm.DB
}

func NewDataBaseCommand(logger *zap.Logger, db *gorm.DB) *DataBaseCommand {
	return &DataBaseCommand{
		logger: logger,
		db:     db,
	}
}

// 数据库迁移，将项目实体与数据库结构同步(代码 -> 数据库)
func (h *DataBaseCommand) Migrate(cmd *cobra.Command, args []string) {
	err := h.db.AutoMigrate(
		&model.User{},
		&model.Media{},
	)

	if err != nil {
		cmd.Println("database migrate error:", err)
	}

	cmd.Println("database migrate success")
}
