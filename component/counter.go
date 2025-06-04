package component

import (
	"strconv"

	"github.com/friedelschoen/st8/notify"
)

func Counter(block *Block, args map[string]string, not *notify.Notification, cacheptr *any) error {
	var count *int
	if *cacheptr != nil {
		count = (*cacheptr).(*int)
	} else {
		initial := 0
		*cacheptr = &initial
		count = &initial
	}

	block.OnClick = func(ClickEvent) {
		(*count)++
	}

	block.Text = strconv.Itoa(*count)
	return nil
}
