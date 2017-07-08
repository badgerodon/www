package rbsa

import (
	"bufio"
	"fmt"
	. "github.com/badgerodon/lalg"
	"github.com/badgerodon/statistics"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	cache           = NewCache(100)
	throttle        = time.NewTicker(time.Second * 1)
	DEFAULT_INDICES = map[string]string{
		"IWB": "Large Cap",
		"IWD": "Large Cap Value",
		"IWF": "Large Cap Growth",
		"IWM": "Small Cap",
		"IWN": "Small Cap Value",
		"IWO": "Small Cap Growth",
		"IWR": "Mid Cap",
		"EEM": "Emerging Markets",
		"ICF": "Real Estate",
		"EFA": "International",
		"AGG": "Fixed Income",
	}
)

func getYahoo(symbol string, year, month int) (*http.Response, error) {
	u := "http://ichart.finance.yahoo.com/table.csv?s=" + url.QueryEscape(symbol) +
		fmt.Sprint("&a=", (month-1), "&b=5&c=", (year-4),
			"&d=", (month-1), "&b=5&c=", year,
			"&ignore=.csv")
	log.Println("GET", u)
	res, err := http.Get(u)
	if err != nil {
		log.Println("- ERROR", err)
		return nil, err
	}
	if res.StatusCode/100 != 2 {
		return nil, fmt.Errorf("Error retrieving data for %v", symbol)
	}
	log.Println("-", res.Status)
	return res, err
}

func getData(symbol string) (Vector, error) {
	t := time.Now()

	y := t.Year()
	m := t.Month()

	if m == 1 {
		m = 12
		y--
	} else {
		m--
	}

	vec, err := cache.Get(fmt.Sprint(y, ":", m, ":", symbol), func() (interface{}, error) {
		<-throttle.C

		r, err := getYahoo(symbol, y, int(m))
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()

		csv := bufio.NewReader(r.Body)

		vec := NewVector(37)

		for i := 0; i <= len(vec); i++ {
			line, err := csv.ReadString('\n')
			if err != nil {
				break
			}

			// Skip the headers
			if i == 0 {
				continue
			}

			// Read the data
			parts := strings.Split(line, ",")
			if len(parts) < 6 {
				continue
			}

			v, err := strconv.ParseFloat(strings.Trim(parts[6], "\r\n"), 64)

			if err != nil {
				v = 0
			}

			vec[i-1] = v
		}
		vec = statistics.Relativize(vec)
		return vec, nil
	})

	if err != nil {
		return nil, err
	}

	return vec.(Vector), nil

}

//http://ichart.finance.yahoo.com/table.csv?s=%5EGSPC&a=00&b=3&c=1950&d=05&e=2&f=2011&g=m&ignore=.csv
func Analyze(id string) (map[string]float64, error) {
	alg := New()
	for k, _ := range DEFAULT_INDICES {
		data, err := getData(k)
		if err != nil {
			return nil, err
		}
		alg.AddIndex(k, data)
	}

	data, err := getData(id)
	if err != nil {
		return nil, err
	}

	solution, err := alg.Run(data)
	if err != nil {
		return nil, err
	}

	return solution, nil
}
