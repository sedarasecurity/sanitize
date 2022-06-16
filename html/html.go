package html

import (
	"bufio"
	"bytes"
	"html"
	"html/template"
	"io"
	"log"
	"regexp"
	"strings"
)

type Config struct {
	stripSingleBreaks          bool
	stripNbsp                  bool
	replaceMultipleSpaces      bool
	replaceNewlines            bool
	replaceNewlinesWithBreaks  bool
	replaceLineEndWithBreaks   bool
	escape                     bool
	escapeFirst                bool
	unescape                   bool
	replaceNewlinesWithLineend bool
}

func ReplaceNewLinesWithLineend(b bool) Option {
	return func(c *Config) {
		c.replaceNewlinesWithLineend = b
	}
}

func ReplaceLineEndWithBreaks(b bool) Option {
	return func(c *Config) {
		c.replaceLineEndWithBreaks = b
	}
}

func StripSingleBreaks(b bool) Option {
	return func(c *Config) {
		c.stripSingleBreaks = b
	}
}

func StripNbsp(b bool) Option {
	return func(c *Config) {
		c.stripNbsp = b
	}
}

func ReplaceMultipleSpaces(b bool) Option {
	return func(c *Config) {
		c.replaceMultipleSpaces = b
	}
}

func ReplaceNewlines(b bool) Option {
	return func(c *Config) {
		c.replaceNewlines = b
	}
}

func ReplaceNewlinesWithBreaks(b bool) Option {
	return func(c *Config) {
		c.replaceNewlinesWithBreaks = b
	}
}

func Escape(b bool) Option {
	return func(c *Config) {
		c.escape = b
	}
}

func EscapeFirst(b bool) Option {
	return func(c *Config) {
		c.escapeFirst = b
	}
}

func Unescape(b bool) Option {
	return func(c *Config) {
		c.unescape = b
	}
}

var defaultHtmlOptConfig = Config{
	replaceMultipleSpaces:    true,
	stripSingleBreaks:        true,
	stripNbsp:                true,
	escape:                   true,
	escapeFirst:              false,
	replaceNewlines:          true,
	replaceLineEndWithBreaks: true,
}

var (
	spacesRegexp      = regexp.MustCompile(`\s{2,}`)
	singleBreakRegexp = regexp.MustCompile(`(?i:^(\s+)?\<br\/?\>(\s+)?$)`)
	nbspRegexp        = regexp.MustCompile(`(?i:\&nbsp\;)`)
)

type Option func(c *Config)

func Sanitize(in string, opts ...Option) string {
	var out string
	cfg := defaultHtmlOptConfig

	for _, o := range opts {
		o(&cfg)
	}

	if cfg.unescape {
		in = html.UnescapeString(in)
	}

	if cfg.escape && cfg.escapeFirst {
		in = escapeHtml(in)
	}

	r := bufio.NewReader(bytes.NewReader([]byte(in)))
	buf := new(bytes.Buffer)

	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("unexpected error reading line: %v\n", err)
		}

		if cfg.stripSingleBreaks {
			if singleBreakRegexp.Match(line) {
				continue
			}
		}

		if cfg.replaceMultipleSpaces {
			line = spacesRegexp.ReplaceAll(line, []byte(" "))
		}

		line = bytes.Replace(line, []byte("\u00a0"), []byte(""), -1)

		buf.Write(line)

		if !cfg.replaceNewlines && !cfg.replaceNewlinesWithBreaks {
			buf.Write([]byte("\r\n"))
		} else if cfg.replaceNewlinesWithBreaks {
			buf.Write([]byte(" <br> "))
		} else if cfg.replaceNewlinesWithLineend {
			buf.Write([]byte(" %LINEEND% "))
		}
	}

	out = buf.String()

	if cfg.escape && !cfg.escapeFirst {
		out = escapeHtml(out)
	}

	if cfg.stripNbsp {
		tmp := nbspRegexp.ReplaceAll([]byte(out), []byte(""))
		out = string(tmp)
	}

	if cfg.replaceLineEndWithBreaks {
		out = strings.Replace(out, "%LINEEND%", " <br> ", -1)
	}

	return out
}

func escapeHtml(in string) string {
	buf := new(bytes.Buffer)
	tmpl := template.New("sanitize")
	tmpl, _ = tmpl.Parse(`{{define "T"}}{{.}}{{end}}`)
	_ = tmpl.ExecuteTemplate(buf, "T", in)

	return strings.TrimSpace(buf.String())
}
