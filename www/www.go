package www

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	abbr "github.com/ultram4rine/raspisos/abbreviation"
	"github.com/ultram4rine/raspisos/data"

	"github.com/PuerkitoBio/goquery"
)

//GetFacs gets all faculties
func GetFacs() ([]data.Fac, error) {
	var (
		addr      = "https://www.sgu.ru/schedule"
		fac       data.Fac
		faculties []data.Fac
	)

	res, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	doc.Find(".panes_item__type_group ul li").Each(func(_ int, s *goquery.Selection) {
		fac.Name = s.Find("a").Text()
		href, _ := s.Find("a").Attr("href")
		fac.Link = strings.Split(href, "/")[2]
		fac.Abbr = abbr.Faculty(fac.Name)
		faculties = append(faculties, fac)
	})

	return faculties, nil
}

//GetTypesofEducation gets types of education on faculty
func GetTypesofEducation(facLink string) ([]string, error) {
	var (
		addr  = "https://www.sgu.ru/schedule/" + facLink
		typ   string
		types []string
	)

	res, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	doc.Find(".form_education").Each(func(_ int, s *goquery.Selection) {
		if s.ChildrenFiltered("legend").Text() == "Дневная форма обучения" {
			s.Find(".group-type").Each(func(_ int, t *goquery.Selection) {
				typ = t.ChildrenFiltered("legend").Text()
				types = append(types, typ)
			})
		}
	})

	return types, nil
}

//GetGroups gets groups of the faculty
func GetGroups(facLink, educationType string) ([]string, error) {
	var (
		addr   = "https://www.sgu.ru/schedule/" + facLink
		group  string
		groups []string
	)

	res, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	doc.Find(".form_education").Each(func(_ int, s *goquery.Selection) {
		if s.ChildrenFiltered("legend").Text() == "Дневная форма обучения" {
			s.Find(".group-type").Each(func(_ int, t *goquery.Selection) {
				typ := t.ChildrenFiltered("legend").Text()
				if typ == educationType {
					t.Find("a").Each(func(_ int, g *goquery.Selection) {
						group = g.Text()
						groups = append(groups, group)
					})
				}
			})
		}
	})

	return groups, nil
}

//DownloadFileFromURL for download
func DownloadFileFromURL(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("bad status: " + resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
