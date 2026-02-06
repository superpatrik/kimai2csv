package kimai

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Project struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Activity struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type TimeSheet struct {
	Id           int       `json:"id"`
	ActivityId   int       `json:"activity"`
	Begin        KimaiTime `json:"begin"`
	End          KimaiTime `json:"end"`
	Duration     int       `json:"duration"`
	Description  string    `json:"description"`
	ProjectId    int       `json:"project"`
	Rate         float64   `json:"rate"`
	InternalRate float64   `json:"internalRate"`
}

type Client struct {
	httpClient *http.Client
	url        string
	apiToken   string
}

type KimaiTime struct {
	time.Time
}

const kimaiLayout = "2006-01-02T15:04:05-0700"

func (t *KimaiTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	parsed, err := time.Parse(kimaiLayout, s)
	if err != nil {
		return err
	}
	t.Time = parsed
	return nil
}

func NewClient(url string, apiToken string) *Client {
	c := &Client{
		httpClient: &http.Client{},
		url:        url,
		apiToken:   apiToken,
	}

	return c
}

func (c *Client) Projects() []Project {
	req, _ := http.NewRequest("GET", c.url+"/projects", nil)
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	res, _ := c.httpClient.Do(req)

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var projects []Project

	json.Unmarshal(body, &projects)

	return projects
}

func (c *Client) Activities() []Activity {
	req, _ := http.NewRequest("GET", c.url+"/activities", nil)
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	res, _ := c.httpClient.Do(req)

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var activities []Activity

	json.Unmarshal(body, &activities)

	return activities
}

func (c *Client) TimeSheets(userId string, begin string, end string, projects []string, size int) []TimeSheet {
	req, _ := http.NewRequest("GET", c.url+"/timesheets", nil)
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	q := req.URL.Query()
	q.Add("user", userId)
	q.Add("begin", begin)
	q.Add("end", end)

	for _, project := range projects {
		q.Add("projects[]", project)
	}

	q.Add("size", strconv.Itoa(size))

	req.URL.RawQuery = q.Encode()

	res, _ := c.httpClient.Do(req)

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var timesheets []TimeSheet

	err := json.Unmarshal(body, &timesheets)

	if err != nil {
		panic(err)
	}

	return timesheets
}
