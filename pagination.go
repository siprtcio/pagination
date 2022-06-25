package pagination

import (
	"errors"
	"math"
	"strconv"
)

// Chapter handles pagination results.
type Chapter struct {
	// The base URL for the endpoint.
	// It is only necessary when using links.
	// Will be omitted from JSON when links are set to false.
	BaseURL string `json:"-"`
	// The next URL string.
	CurrentPageURI string `json:"uri,omitempty"`
	NextPageURI string `json:"next_page_uri,omitempty"`
	FirstPageURI string `json:"first_page_uri,omitempty"`
	PreviousPageURI string `json:"previous_page_uri,omitempty"`
	// The inicial offset position.
	Offset int `json:"-"`
	Limit int `json:"page_size"`
	// The new page number captured on the request params.
	// Will be omitted from JSON, since there is no need for it.
	NewPage int `json:"-"`
	// The current page of the tome.
	// If none is provided, the current page will be setted to 1.
	CurrentPage int `json:"page"`
	// The last page of the tome.
	LastPage int `json:"last_page"`
	Start int `json:"start"`
	End int `json:"end"`
	// The total results, this usually comes from
	// a database query.
	TotalResults int `json:"total_results"`
}

// Paginate handles the pagination calculation.
func (c *Chapter) Paginate() error {
	c.setDefaults()

	if err := c.ceilLastPage(); err != nil {
		return err
	}

	if err := c.doPaginate(); err != nil {
		return err
	}

	if err := c.checkLinks(); err != nil {
		return err
	}

	return nil
}

// Calculates the offset and the limit.
func (c *Chapter) doPaginate() error {
	if c.NewPage < 0 {
		c.NewPage = 0
	}

	if c.NewPage > c.CurrentPage {
		c.CurrentPage = c.NewPage
		c.Offset = (c.CurrentPage) * c.Limit
	}
	
	c.Start = c.Offset
	c.End = c.TotalResults - 1
	if c.CurrentPage != c.LastPage {
	   c.End = c.Offset + c.Limit
	}

	return nil
}

// Ceils the last page and generates
// a integer number.
func (c *Chapter) ceilLastPage() error {
	if c.TotalResults == 0 {
		return errors.New("TotalResults value is missing")
	}

	c.LastPage = int(math.Floor(float64(c.TotalResults) / float64(c.Limit)))

	return nil
}

// Handles links validations.
func (c *Chapter) checkLinks() error {
	if err := c.createLinks(); err != nil {
		return err
	}
	return nil
}

// Creates next and previous links using
// the given base URL.
func (c *Chapter) createLinks() error {
	if c.BaseURL == "" {
		return errors.New("BaseURL value is missing")
	}

	c.CurrentPageURI = c.BaseURL + "?Page=" + strconv.Itoa(c.CurrentPage) + "&PageSize=" + strconv.Itoa(c.Limit)
	
	c.FirstPageURI = c.BaseURL + "?Page=0" + "&PageSize=" + strconv.Itoa(c.Limit)
	
	if c.CurrentPage < c.LastPage {
		c.NextPageURI = c.BaseURL + "?Page=" + strconv.Itoa(c.CurrentPage+1) + "&PageSize=" + strconv.Itoa(c.Limit)
	}

	if c.LastPage > c.CurrentPage {
		c.PreviousPageURI = c.BaseURL + "?Page=" + strconv.Itoa(c.CurrentPage-1) + "&PageSize=" + strconv.Itoa(c.Limit)
	}

	return nil
}

// Sets the defaults values for current page
// and limit if none of them were provided.
func (c *Chapter) setDefaults() {
	if c.Limit == 0 {
		c.Limit = 10
	}
}
