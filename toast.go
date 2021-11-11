package toast

type Audio string

type NotificationOption func(*notification)

// WithTitle
//
// The main title/heading for the notification.
func WithTitle(title string) NotificationOption {
	return func(n *notification) {
		n.Title = title
	}
}

// WithMessage
//
// The single/multi line message to display for the notification.
func WithMessage(msg string) NotificationOption {
	return func(n *notification) {
		n.Message = msg
	}
}

// WithAudio
//
// The audio to play when displaying the notification
func WithAudio(audio Audio) NotificationOption {
	return func(n *notification) {
		n.Audio = audio
	}
}

type notifier interface {
	push() error
}

func Push(message string, opts ...NotificationOption) error {
	var n notifier = newNotification(message, opts...)
	return n.push()
}
