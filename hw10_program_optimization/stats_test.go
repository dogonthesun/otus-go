//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":6,"Name":"Janice Rose","Username":"KeithHart","Email":"Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":7,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":8,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@domain.sale","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
{"Id":9,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@domain.forsale","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}
`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("empty buffer", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(""), "ru")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

	t.Run("empty domain for filtering", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com":   2,
			"linktype.com":    1,
			"browsedrive.gov": 1,
			"teklist.net":     1,
			"domain.sale":     1,
			"domain.forsale":  1,
		}, result)
	})

	t.Run("same domain suffix", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "sale")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"domain.sale": 1}, result)
	})
}
