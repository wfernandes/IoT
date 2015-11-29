package notifiers

type Notifier interface {
	Notify(string) error
}
