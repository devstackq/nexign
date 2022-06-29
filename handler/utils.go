package handler

import (
	"strings"
	"sync"

	"github.com/devstackq/nexign/config"
	"github.com/devstackq/nexign/entity"
)

func (tg *textgears) concurrentRequest(cfg *config.Config) {
	var (
		wg     sync.WaitGroup
		result = make(chan entity.Result, 1)
		errCh  = make(chan error, 1)
	)

	for idxSent, sentence := range tg.Text {
		wg.Add(1)
		go tg.spelling(cfg, sentence, result, errCh)

		tempRes := <-result

		if <-errCh == nil {
			// var wg2 sync.WaitGroup
			// go func(tempRes entity.Result, idxSent int) {
			if len(tempRes.Errors) > 0 {
				// wg2.Add(1)
				f := ""
				for _, word := range tempRes.Errors {
					f = strings.Replace(*tg.Text[idxSent], word.Bad, word.Better[0], -1)
					tg.Text[idxSent] = &f
				}
				// wg2.Done()
			}
			// }(<-result, idxSent)
			// wg2.Wait()
		}
		wg.Done()
	}
	wg.Wait()
}
