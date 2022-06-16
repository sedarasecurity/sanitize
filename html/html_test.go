package html

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCorpus = [][]string{
	{
		"&lt;br&gt;&lt;font style=&#34;font-family: Arial;font-size: 12px;&#34;&gt;&lt;br&gt;&lt;b&gt;Device Name : MULTISORB-EDGE-FW-01&lt;/b&gt;&lt;br&gt;&lt;br&gt;Following alert(s) were generated as per the Device time 2017-05-15 17:10 and time zone - America/New_York configured in Device.&lt;br&gt;&lt;br&gt;     Critical Intrusion attacks with serial number C40057F2MRBJV6A is 2 in the last 5 minute(s).&lt;br&gt;&lt;/font&gt;&lt;br&gt;",
		"&lt;br&gt;&lt;font style=&#34;font-family: Arial;font-size: 12px;&#34;&gt;&lt;br&gt;&lt;b&gt;Device Name : MULTISORB-EDGE-FW-01&lt;/b&gt;&lt;br&gt;&lt;br&gt;Following alert(s) were generated as per the Device time 2017-05-15 17:10 and time zone - America/New_York configured in Device.&lt;br&gt;&lt;br&gt;Critical Intrusion attacks with serial number C40057F2MRBJV6A is 2 in the last 5 minute(s).&lt;br&gt;&lt;/font&gt;&lt;br&gt;",
	},
	{
		"a     b     c                              d     e",
		"a b c d e",
	},
	{
		`<p>this is a test</p><br/>
    <br/>
    <p>this is a test</p><br/>
    <br/>
    <br/>`,
		"&lt;p&gt;this is a test&lt;/p&gt;&lt;br/&gt; &lt;p&gt;this is a test&lt;/p&gt;&lt;br/&gt;",
	},
}

func TestSanitize(t *testing.T) {
	assert := assert.New(t)

	opts := []Option{
		Unescape(true),
		StripNbsp(true),
		ReplaceMultipleSpaces(true),
		ReplaceNewlines(true),
		StripSingleBreaks(true),
		Escape(true),
	}

	for _, corpus := range testCorpus {
		output := Sanitize(corpus[0], opts...)
		assert.Equal(corpus[1], output)
	}
}
