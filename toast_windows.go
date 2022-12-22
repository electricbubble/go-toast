//go:build windows

package toast

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"text/template"
	"time"
	"unsafe"
)

// WithAppID
//
// The name of your app. This value shows up in Windows 10's Action Centre, so make it
// something readable for your users. It can contain spaces, however special characters
// (eg. é) are not supported.
func WithAppID(appID string) NotificationOption {
	return func(n *notification) {
		n.AppID = appID
	}
}

// WithIcon
//
// An optional path to an image on the OS to display to the left of the title & message.
func WithIcon(pathIcon string) NotificationOption {
	return func(n *notification) {
		n.Icon = pathIcon
	}
}

func WithIconRaw(raw []byte) NotificationOption {
	return func(n *notification) {
		randBytes := make([]byte, 4)
		_r.Read(randBytes)
		n._tmpIconFilename = filepath.Join(os.TempDir(), fmt.Sprintf("go-toast-logo-%x.png", randBytes))
		if err := os.WriteFile(n._tmpIconFilename, raw, 0600); err != nil {
			return
		}
		n.Icon = n._tmpIconFilename
	}
}

// WithActivationType
//
// The type of notification level action (like Action)
func WithActivationType(activationType string) NotificationOption {
	return func(n *notification) {
		n.ActivationType = activationType
	}
}

// WithActivationArguments
//
// // The activation/action arguments (invoked when the user clicks the notification)
func WithActivationArguments(activationArguments string) NotificationOption {
	return func(n *notification) {
		n.ActivationArguments = activationArguments
	}
}

// WithProtocolAction
//
// Defines an actionable button.
// See https://msdn.microsoft.com/en-us/windows/uwp/controls-and-patterns/tiles-and-notifications-adaptive-interactive-toasts for more info.
//
// Only protocol type action buttons are actually useful, as there's no way of receiving feedback from the
// user's choice. Examples of protocol type action buttons include: "bingmaps:?q=sushi" to open up Windows 10's
// maps app with a pre-populated search field set to "sushi".
//
//     Action{"protocol", "Open Maps", "bingmaps:?q=sushi"}
func WithProtocolAction(label string, arguments ...string) NotificationOption {
	return func(n *notification) {
		if len(n.Actions) == 0 {
			n.Actions = make([]Action, 0, 5)
		}
		if len(n.Actions) == 5 {
			return
		}
		if len(arguments) == 0 {
			arguments = []string{""}
		}
		n.Actions = append(n.Actions, Action{
			Type:      "protocol",
			Label:     escapeNotificationString(label),
			Arguments: arguments[0],
		})
	}
}

// WithAudioLoop
//
// Whether to loop the audio (default false)
func WithAudioLoop(b bool) NotificationOption {
	return func(n *notification) {
		n.Loop = b
	}
}

// WithDuration
//
// How long the notification should show up for (short/long)
func WithDuration(nd NotificationDuration) NotificationOption {
	return func(n *notification) {
		n.Duration = nd
	}
}

func WithLongDuration() NotificationOption {
	return func(n *notification) {
		n.Duration = Long
	}
}

func WithShortDuration() NotificationOption {
	return func(n *notification) {
		n.Duration = Short
	}
}

type NotificationDuration string

const (
	Short NotificationDuration = "short"
	Long  NotificationDuration = "long"
)

const (
	Silent         Audio = "silent"
	Default        Audio = "ms-winsoundevent:Notification.Default"
	IM             Audio = "ms-winsoundevent:Notification.IM"
	Mail           Audio = "ms-winsoundevent:Notification.Mail"
	Reminder       Audio = "ms-winsoundevent:Notification.Reminder"
	SMS            Audio = "ms-winsoundevent:Notification.SMS"
	LoopingAlarm   Audio = "ms-winsoundevent:Notification.Looping.Alarm"
	LoopingAlarm2  Audio = "ms-winsoundevent:Notification.Looping.Alarm2"
	LoopingAlarm3  Audio = "ms-winsoundevent:Notification.Looping.Alarm3"
	LoopingAlarm4  Audio = "ms-winsoundevent:Notification.Looping.Alarm4"
	LoopingAlarm5  Audio = "ms-winsoundevent:Notification.Looping.Alarm5"
	LoopingAlarm6  Audio = "ms-winsoundevent:Notification.Looping.Alarm6"
	LoopingAlarm7  Audio = "ms-winsoundevent:Notification.Looping.Alarm7"
	LoopingAlarm8  Audio = "ms-winsoundevent:Notification.Looping.Alarm8"
	LoopingAlarm9  Audio = "ms-winsoundevent:Notification.Looping.Alarm9"
	LoopingAlarm10 Audio = "ms-winsoundevent:Notification.Looping.Alarm10"
	LoopingCall    Audio = "ms-winsoundevent:Notification.Looping.Call"
	LoopingCall2   Audio = "ms-winsoundevent:Notification.Looping.Call2"
	LoopingCall3   Audio = "ms-winsoundevent:Notification.Looping.Call3"
	LoopingCall4   Audio = "ms-winsoundevent:Notification.Looping.Call4"
	LoopingCall5   Audio = "ms-winsoundevent:Notification.Looping.Call5"
	LoopingCall6   Audio = "ms-winsoundevent:Notification.Looping.Call6"
	LoopingCall7   Audio = "ms-winsoundevent:Notification.Looping.Call7"
	LoopingCall8   Audio = "ms-winsoundevent:Notification.Looping.Call8"
	LoopingCall9   Audio = "ms-winsoundevent:Notification.Looping.Call9"
	LoopingCall10  Audio = "ms-winsoundevent:Notification.Looping.Call10"
)

// Action
//
// Defines an actionable button.
// See https://msdn.microsoft.com/en-us/windows/uwp/controls-and-patterns/tiles-and-notifications-adaptive-interactive-toasts for more info.
//
// Only protocol type action buttons are actually useful, as there's no way of receiving feedback from the
// user's choice. Examples of protocol type action buttons include: "bingmaps:?q=sushi" to open up Windows 10's
// maps app with a pre-populated search field set to "sushi".
//
//     Action{"protocol", "Open Maps", "bingmaps:?q=sushi"}
type Action struct {
	Type      string
	Label     string
	Arguments string
}

var _ notifier = (*notification)(nil)

func newNotification(message string, opts ...NotificationOption) *notification {
	n := &notification{
		AppID:          "GO APP",
		Title:          "GO APP",
		Message:        message,
		ActivationType: "protocol",
		Duration:       Short,
		Audio:          Silent,
	}
	for _, fn := range opts {
		fn(n)
	}
	n.AppID = escapeNotificationString(n.AppID)
	n.Title = escapeNotificationString(n.Title)
	n.Message = escapeNotificationString(n.Message)
	return n
}

func (n *notification) push() error {
	content, err := n.template()
	if err != nil {
		return err
	}

	randBytes := make([]byte, 4)
	_r.Read(randBytes)
	tmpFilename := filepath.Join(os.TempDir(), fmt.Sprintf("go-toast-%x.ps1", randBytes))

	if err = os.WriteFile(tmpFilename, content, 0600); err != nil {
		return err
	}

	defer func() {
		_ = os.Remove(tmpFilename)
	}()

	launch := "(Get-Content -Encoding UTF8 -Path " + tmpFilename + " -Raw) | Invoke-Expression"
	if len(n._tmpIconFilename) != 0 {
		launch += "; Start-Sleep -m 50 ; Remove-Item " + n._tmpIconFilename
	}
	cmd := exec.Command("PowerShell", "-ExecutionPolicy", "Bypass", launch)
	fixCmd("PowerShell", cmd)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

var (
	_r    = rand.New(rand.NewSource(time.Now().Unix()))
	_tpl  *template.Template
	_once sync.Once
)

func (n *notification) template() (content []byte, err error) {
	_once.Do(func() {
		var tplNotification = `
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.UI.Notifications.ToastNotification, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null

$APP_ID = '{{if .AppID}}{{.AppID}}{{else}}Windows App{{end}}'

$template = @"
<toast activationType="{{.ActivationType}}" launch="{{.ActivationArguments}}" duration="{{.Duration}}">
    <visual>
        <binding template="ToastGeneric">
            {{if .Icon}}
            <image placement="appLogoOverride" src="{{.Icon}}" />
            {{end}}
            {{if .Title}}
            <text><![CDATA[{{.Title}}]]></text>
            {{end}}
            {{if .Message}}
            <text><![CDATA[{{.Message}}]]></text>
            {{end}}
        </binding>
    </visual>
    {{if ne .Audio "silent"}}
	<audio src="{{.Audio}}" loop="{{.Loop}}" />
	{{else}}
	<audio silent="true" />
	{{end}}
    {{if .Actions}}
    <actions>
        {{range .Actions}}
        <action activationType="{{.Type}}" content="{{.Label}}" arguments="{{.Arguments}}" />
        {{end}}
    </actions>
    {{end}}
</toast>
"@

$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($template)
$go_toast = New-Object Windows.UI.Notifications.ToastNotification $xml
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier($APP_ID).Show($go_toast)
`

		_tpl, err = template.New("_tpl").Parse(tplNotification)
	})
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	err = _tpl.Execute(buf, n)
	return buf.Bytes(), err
}

type notification struct {
	// The name of your app. This value shows up in Windows 10's Action Centre, so make it
	// something readable for your users. It can contain spaces, however special characters
	// (eg. é) are not supported.
	AppID string

	// The main title/heading for the notification.
	Title string

	// The single/multi line message to display for the notification.
	Message string

	// An optional path to an image on the OS to display to the left of the title & message.
	Icon             string
	_tmpIconFilename string

	// The type of notification level action (like Action)
	ActivationType string

	// The activation/action arguments (invoked when the user clicks the notification)
	ActivationArguments string

	// Optional action buttons to display below the notification title & message.
	Actions []Action

	// The audio to play when displaying the notification
	Audio Audio

	// Whether to loop the audio (default false)
	Loop bool

	// How long the notification should show up for (short/long)
	Duration NotificationDuration
}

func escapeNotificationString(in string) string {
	noSlash := strings.ReplaceAll(in, "`", "``")
	return strings.ReplaceAll(noSlash, "\"", "`\"")
}

// https://pkg.go.dev/golang.org/x/sys/execabs#Command
func fixCmd(name string, cmd *exec.Cmd) {
	if filepath.Base(name) == name && !filepath.IsAbs(cmd.Path) {
		// exec.Command was called with a bare binary name and
		// exec.LookPath returned a path which is not absolute.
		// Set cmd.lookPathErr and clear cmd.Path so that it
		// cannot be run.
		lookPathErr := (*error)(unsafe.Pointer(reflect.ValueOf(cmd).Elem().FieldByName("lookPathErr").Addr().Pointer()))
		if *lookPathErr == nil {
			*lookPathErr = relError(name, cmd.Path)
		}
		cmd.Path = ""
	}
}

func relError(file, path string) error {
	return fmt.Errorf("%s resolves to executable in current directory (.%c%s)", file, filepath.Separator, path)
}
