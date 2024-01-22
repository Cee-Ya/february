package mid

import (
	"february/pkg/logx"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/goccy/go-json"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"html"
	"io"
	"net/url"
	"strings"
)

// XSSFilter filters XSS in query string.
func XSSFilter(whitelistURLs []string) gin.HandlerFunc {
	// Do this once for each unique policy, and use the policy for the life of the
	// program Policy creation/editing is not safe to use in multiple goroutines.
	p := bluemonday.UGCPolicy()

	return func(c *gin.Context) {
		for _, u := range whitelistURLs {
			if strings.HasPrefix(c.Request.URL.String(), u) {
				c.Next()
				return
			}
		}

		sanitizedQuery, err := xssFilterQuery(p, c.Request.URL.RawQuery)
		if err != nil {
			logx.ErrorF(c, "XSS error:: ", zap.Error(err))
			c.Abort()
			return
		}
		c.Request.URL.RawQuery = sanitizedQuery

		var sanitizedBody string
		bding := binding.Default(c.Request.Method, c.ContentType())
		body, err := c.GetRawData()
		if err != nil {
			logx.ErrorF(c, "XSS error:: ", zap.Error(err))
			c.Abort()
			return
		}

		// xssFilterJSON() will return error when body is empty.
		if len(body) == 0 {
			c.Next()
			return
		}

		switch bding {
		case binding.JSON:
			sanitizedBody, err = xssFilterJSON(p, string(body))
		case binding.FormMultipart:
			sanitizedBody = xssFilterPlain(p, string(body))
		case binding.Form:
			sanitizedBody, err = xssFilterQuery(p, string(body))
		}

		if err != nil {
			logx.ErrorF(c, "XSS error:: ", zap.Error(err))
			c.Abort()
			return
		}
		c.Request.Body = io.NopCloser(strings.NewReader(sanitizedBody))

		c.Next()
	}
}

func xssFilterPlain(p *bluemonday.Policy, s string) string {
	sanitized := p.Sanitize(s)
	return html.UnescapeString(sanitized)
}

func xssFilterJSON(p *bluemonday.Policy, s string) (string, error) {
	var data interface{}
	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		return "", err
	}

	b := strings.Builder{}
	e := json.NewEncoder(&b)
	e.SetEscapeHTML(false)
	err = e.Encode(xssFilterJSONData(p, data))
	if err != nil {
		return "", err
	}
	// use `TrimSpace` to trim newline char add by `Encode`.
	return strings.TrimSpace(b.String()), nil
}

func xssFilterJSONData(p *bluemonday.Policy, data interface{}) interface{} {
	if s, ok := data.([]interface{}); ok {
		for i, v := range s {
			s[i] = xssFilterJSONData(p, v)
		}
		return s
	} else if m, ok := data.(map[string]interface{}); ok {
		for k, v := range m {
			m[k] = xssFilterJSONData(p, v)
		}
		return m
	} else if str, ok := data.(string); ok {
		return xssFilterPlain(p, str)
	}
	return data
}

func xssFilterQuery(p *bluemonday.Policy, s string) (string, error) {
	values, err := url.ParseQuery(s)
	if err != nil {
		return "", err
	}

	for k, v := range values {
		values.Del(k)
		for _, vv := range v {
			values.Add(k, xssFilterPlain(p, vv))
		}
	}

	return values.Encode(), nil
}
