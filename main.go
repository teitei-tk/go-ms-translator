package malwiya

import "github.com/teitei-tk/malwiya/translator"

type (
	Trasnlate func(text, from, to string) (string, error)

	Malwiya struct {
		SubscriptionKey string
	}
)

func New(subscriptionKey string) *Malwiya {
	m := &Malwiya{
		SubscriptionKey: subscriptionKey,
	}

	return m
}

func (m *Malwiya) Translate(text, from, to string) (string, error) {
	return translator.Translate(m.SubscriptionKey, text, from, to)
}
