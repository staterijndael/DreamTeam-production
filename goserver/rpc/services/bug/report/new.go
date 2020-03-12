package report

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

//сохраняет баг репорт. не требует четкой структуры. но чем больше инфы вы засунете в json - тем лучше.
//zenrpc:report JSON любой структуры
func (s *Service) New(report interface{}) {
	msg, err1 := json.Marshal(report)
	dir, err2 := filepath.Abs(s.conf.BugReportsDir)
	if err1 != nil || err2 != nil {
		return
	}

	filename := fmt.Sprintf("%d.json", time.Now().Unix())
	fp := filepath.Join(dir, filename)
	ioutil.WriteFile(fp, msg, os.ModePerm)
}
