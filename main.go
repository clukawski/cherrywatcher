package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// CherryRIPData is parsed WikiData json response from a defined sparql query
// See queries.go for predifined query + a test query of another Don Cherry
// who is currently dead
type CherryRIPData struct {
	Head struct {
		Vars []string `json:"vars"`
	} `json:"head"`
	Results struct {
		Bindings []struct {
			Item struct {
				Type  string `json:"type"`
				Value string `json:"value"`
			} `json:"item"`
			ItemLabel struct {
				XMLLang string `json:"xml:lang"`
				Type    string `json:"type"`
				Value   string `json:"value"`
			} `json:"itemLabel"`
			Rip struct {
				Datatype string    `json:"datatype,omitempty"`
				Type     string    `json:"type,omitempty"`
				Value    time.Time `json:"value,omitempty"`
			} `json:"RIP,omitempty"`
		} `json:"bindings"`
	} `json:"results"`
}

// Flags
var (
	test        bool
	pushToken   string
	queryString string
	pushURL     string
)

func main() {
	flag.BoolVar(&test, "t", false, "Test service with a different, but dead Don Cherry")
	flag.StringVar(&pushToken, "p", "", "Token used for gotify server push POST request")
	flag.StringVar(&pushURL, "u", "", "URL for gotify server")
	flag.Parse()

	// Don't run without the pushToken
	if pushToken == "" {
		log.Fatalln("missing push token for gotify")
	}

	// Don't run without the pushURL
	if pushURL == "" {
		log.Fatalln("missing push token for gotify")
	}

	// If we're testing use the Don Cherry that isn't the one we care about, but is dead
	if test {
		queryString = otherCherryRIPQuery
	} else {
		queryString = cherryRIPQuery
	}

	// Parse our dumbass query
	q, err := url.Parse(queryString)
	if err != nil {
		log.Fatalf("fatal: %s", err)
	}

	// Make our client with timeout defaults
	c := &http.Client{
		Timeout: 15 * time.Second,
	}

	// Include this header or you'll get XML back!
	// -H 'Accept: application/sparql-results+json'
	req := &http.Request{
		Method: "GET",
		URL:    q,
		Header: http.Header{`Accept`: []string{`application/sparql-results+json`}},
	}

	// Send it a few times just so we know for sure
	var count int
	// Loop making request every minute until we find out don cherry is dead
	for {
		if dead, err := checkCherryRIP(c, req); err != nil {
			log.Println("client: %s", err)
		} else if dead {
			if err = pushCherryRIP(c); err != nil {
				log.Printf("unable to send push request about don cherry being dead, dang: %s", err)
			} else if count > 5 {
				// We have nothing left to do
				os.Exit(0)
			}
			count++
		}
		time.Sleep(time.Minute)
	}
}

// checkCherryRIP checks if Don Cherry is dead yet or not
func checkCherryRIP(c *http.Client, req *http.Request) (dead bool, err error) {
	// Make the request and fetch the body
	res, err := c.Do(req)
	if err != nil {
		return false, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	// Unmarshal the WikiData response json
	crd := CherryRIPData{}
	err = json.Unmarshal(body, &crd)
	if err != nil {
		return false, err
	}

	// Check if ?RIP is set, if so, he's dead
	if len(crd.Results.Bindings) > 0 {
		dead = crd.Results.Bindings[0].Rip.Type != ""
	}

	return dead, nil
}

// pushCherryRIP sends a POST request to local Gotify server to send a push
// notification to registered devices when Don Cherry is dead
func pushCherryRIP(c *http.Client) error {
	// -F "title=my title" -F "message=my message" -F "priority=5"
	formValues := map[string]io.Reader{
		`title`:    strings.NewReader(`🏒 🥳 🎉 DON CHERRY IS DEAD NOW!! REJOICE!! 🎉 🥳 🏒`),
		`message`:  strings.NewReader(`🎉 🥳 🏒 DON CHERRY'S DEAD LETS NOT TALK ABOUT ANY GOOD GUYS!!! 🏒 🥳 🎉`),
		`priority`: strings.NewReader(`1`),
	}

	// Prepare a multipart form for Gotify
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range formValues {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add file
		var err error
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return err
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return err
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}

	}
	w.Close()

	// Craft our request including our multipart io.Writer
	url := fmt.Sprintf("https://%s/message?token=%s", pushURL, pushToken)
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return err
	}
	// Set Form Data Content-Type header
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Make the request to Gotify
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	// Check the response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	log.Printf("gotify response body: %s", string(body))
	return nil
}
