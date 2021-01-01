package channelid

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestDefaultId(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	channelId := New()
	t.Log(channelId.AsShortText())
	t.Log(channelId.AsLongText())
}
