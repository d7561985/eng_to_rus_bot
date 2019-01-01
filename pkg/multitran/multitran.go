package multitran

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"strings"
	"unicode/utf8"
)

const UrlRequestTemplate = "https://www.multitran.com/m.exe?l1=1&l2=2&SHL=2&s=:w"

type Res struct {
	Sub, Trans string
}

// Res translate work representation.
type Result struct {
	List      []Res
	Word, Url string
}

// String output results.
// @limit - max trans. limit
func (r Result) String(limit int) string {
	buf := bytes.NewBufferString(r.Url)
	buf.WriteString("\n")

	for _, v := range r.List {
		buf.WriteString(v.Sub)
		buf.WriteString(" -> ")
		if limit > 0 && len(v.Trans) > limit {
			V := 0
			for ; !utf8.ValidString(v.Trans[:limit+V]); V++ {

			}
			tr := v.Trans[:limit+V] + "..."
			buf.WriteString(tr)
		} else {
			buf.WriteString(v.Trans)
		}

		buf.WriteString("\n")
	}
	return buf.String()
}

// HasSub check if exist in list current sub
func (r Result) HasSub(sub string) bool {
	for _, v := range r.List {
		if v.Sub == sub {
			return true
		}
	}
	return false
}

// GetWord translation from url.
func GetWord(word string) (*Result, error) {
	w := strings.Replace(word, " ", "+", -1)
	url := strings.Replace(UrlRequestTemplate, ":w", w, -1)
	//fmt.Println(url)

	res := &Result{List: make([]Res, 0), Url: url, Word: word}

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, nil
	}

	sub, skip := "", false
	doc.Find(".subj,.trans").Each(func(i int, selection *goquery.Selection) {
		if selection.Get(0).Attr[0].Val == "subj" {
			sub = strings.Trim(selection.Text(), " , - >")
			if strings.Contains(sub, ` `) {
				skip = true
			}

		} else if selection.Get(0).Attr[0].Val == "trans" {
			if skip || res.HasSub(sub) {
				skip = false
				return
			}
			c := (i + 1) / 2
			if len(res.List) > c {
				log.Error().Msgf("search idx %d, trans. count %d", i, c)
				return
			}
			res.List = append(res.List, Res{Sub: sub, Trans: selection.Text()})
		}
	})

	return res, nil
}
