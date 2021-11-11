//go:build darwin

package toast

import (
	"fmt"
	"os/exec"
	"strings"
)

func WithSubtitle(subtitle string) NotificationOption {
	return func(n *notification) {
		n.Subtitle = subtitle
	}
}

const (
	Basso     Audio = "Basso"
	Blow      Audio = "Blow"
	Bottle    Audio = "Bottle"
	Frog      Audio = "Frog"
	Funk      Audio = "Funk"
	Glass     Audio = "Glass"
	Hero      Audio = "Hero"
	Morse     Audio = "Morse"
	Ping      Audio = "Ping"
	Pop       Audio = "Pop"
	Purr      Audio = "Purr"
	Sosumi    Audio = "Sosumi"
	Submarine Audio = "Submarine"
	Tink      Audio = "Tink"
)

var _ notifier = (*notification)(nil)

func newNotification(message string, opts ...NotificationOption) *notification {
	n := &notification{
		Title:   "GO APP",
		Message: message,
	}
	for _, fn := range opts {
		fn(n)
	}
	return n
}

func (n *notification) push() error {
	script := n.template()
	osa, err := exec.LookPath("osascript")
	if err != nil {
		return err
	}
	cmd := exec.Command(osa, "-e", script)
	return cmd.Run()
}

func (n *notification) template() (script string) {
	tpl := `display notification "%s" with title "%s"`
	script = fmt.Sprintf(tpl, escapeNotificationString(n.Message), escapeNotificationString(n.Title))
	if len(n.Subtitle) != 0 {
		script += fmt.Sprintf(` subtitle "%s"`, escapeNotificationString(n.Subtitle))
	}
	if len(n.Audio) != 0 {
		script += fmt.Sprintf(` sound name "%s"`, escapeNotificationString(string(n.Audio)))
	}
	return
}

type notification struct {
	// The main title/heading for the notification.
	Title string

	Subtitle string

	// The single/multi line message to display for the notification.
	Message string

	// The audio to play when displaying the notification
	Audio Audio
}

func escapeNotificationString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	return strings.ReplaceAll(s, `"`, `\"`)
}
