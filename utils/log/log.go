package log

import (
	"fmt"

	"github.com/fatih/color"
)

type LogHelper struct {
	Index  int
	UserID string
	IP     string
}

func NewLogHelper(index int, userId string) *LogHelper {
	return &LogHelper{
		Index:  index,
		UserID: userId,
		IP:     "🖥️",
	}
}

func (lh *LogHelper) Log(msg string) {
	fmt.Printf("[ No %d _ ID: %s _ IP: %s ] %s\n", lh.Index, lh.UserID, lh.IP, msg)
}

func (lh *LogHelper) LogError(msg string) {
	fmt.Printf("[ No %d _ ID: %s _ IP: %s ] %s\n", lh.Index, lh.UserID, lh.IP, color.RedString(msg))
}

func (lh *LogHelper) LogSuccess(msg string) {
	fmt.Printf("[ No %d _ ID: %s _ IP: %s ] %s\n", lh.Index, lh.UserID, lh.IP, color.GreenString(msg))
}

func (lh *LogHelper) UpdateIp(ip string) {
	lh.IP = ip
}
