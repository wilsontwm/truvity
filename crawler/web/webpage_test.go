package web_test

import (
	"crawler/web"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Process(t *testing.T) {
	url := "https://www.brandeis.edu/student-affairs"

	wp, err := web.NewWebpage(url)
	assert.Nil(t, err)

	path := "../test"
	os.RemoveAll(path)

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
			return
		}
	} else if err != nil {
		log.Println(err)
		return
	}

	err = wp.Process(url, fmt.Sprintf("./%s", path))
	assert.Nil(t, err)
	assert.Equal(t, 7, len(wp.ChildURLs))
	sort.Strings(wp.ChildURLs)
	assert.Equal(t,
		[]string{
			"https://www.brandeis.edu/student-affairs/about/departments.html",
			"https://www.brandeis.edu/student-affairs/about/index.html",
			"https://www.brandeis.edu/student-affairs/about/leadership.html",
			"https://www.brandeis.edu/student-affairs/about/organizational-structure/index.html",
			"https://www.brandeis.edu/student-affairs/staff/index.html",
			"https://www.brandeis.edu/student-affairs/students/index.html",
			"https://www.brandeis.edu/student-affairs/students/leadership.html",
		},
		wp.ChildURLs,
	)

}
