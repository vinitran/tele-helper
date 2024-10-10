package log

import (
	"fmt"
	"log"

	"github.com/fatih/color"
)

type LogHelper struct {
	UserID string
	IP     string
	App    string
}

func NewLogHelper(userId, proxy, app string) *LogHelper {
	return &LogHelper{
		UserID: userId,
		IP:     fmt.Sprintf("🖥️%s", proxy),
		App:    app,
	}
}

func (lh *LogHelper) Normal(msg string) {
	if lh.App == "" {
		log.Printf("[ID: %s _ IP: %s ] %s\n", lh.UserID, lh.IP, msg)
		return
	}
	log.Printf("[ID: %s _ IP: %s _ MiniApp: %s ] %s\n", lh.UserID, lh.IP, lh.App, msg)
}

func (lh *LogHelper) Error(msg string) {
	if lh.App == "" {
		log.Printf("[ID: %s _ IP: %s ] %s\n", lh.UserID, lh.IP, color.RedString(msg))
		return
	}
	log.Printf("[ID: %s _ IP: %s _ MiniApp: %s ] %s\n", lh.UserID, lh.IP, lh.App, color.RedString(msg))
}

func (lh *LogHelper) Success(msg string) {
	if lh.App == "" {
		log.Printf("[ID: %s _ IP: %s ] %s\n", lh.UserID, lh.IP, color.RedString(msg))
		return
	}
	log.Printf("[ID: %s _ IP: %s _ MiniApp: %s ] %s\n", lh.UserID, lh.IP, lh.App, color.GreenString(msg))
}

func (lh *LogHelper) UpdateIp(ip string) {
	lh.IP = ip
}

func (lh *LogHelper) ErrorMessage(err error) error {
	if lh.App == "" {
		return fmt.Errorf("[ID: %s _ IP: %s ] %s\n", lh.UserID, lh.IP, color.RedString(fmt.Sprintf("%e", err)))
	}
	return fmt.Errorf("[ID: %s _ IP: %s _ MiniApp: %s ] %s\n", lh.UserID, lh.IP, lh.App, color.RedString(err.Error()))
}

func (lh *LogHelper) ErrorWithMsg(msg string) error {
	if lh.App == "" {
		return fmt.Errorf("[ID: %s _ IP: %s ] %s\n", lh.UserID, lh.IP, color.RedString(msg))
	}
	return fmt.Errorf("[ID: %s _ IP: %s _ MiniApp: %s ] %s\n", lh.UserID, lh.IP, lh.App, color.RedString(msg))
}
