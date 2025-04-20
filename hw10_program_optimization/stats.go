package hw10programoptimization

import (
	"bufio"
	"bytes"
	"io"

	"github.com/buger/jsonparser"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	stat := make(DomainStat)
	domainSuffix := []byte{}
	if len(domain) > 0 {
		domainSuffix = []byte("." + domain)
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		addr, fieldType, _, err := jsonparser.Get(scanner.Bytes(), "Email")
		if err != nil || fieldType != jsonparser.String {
			continue
		}
		if atCharIdx := fixAddr(addr, domainSuffix); atCharIdx != -1 {
			stat[string(addr[atCharIdx+1:])]++
		}
	}
	return stat, nil
}

// fixAddr replaces uppercase letters with lowercase and
// returns index of the latest character '@'
// if the character is not the last one in buf
// and if lowercase version of buf ends with domain characters. Otherwise,
// the function returns -1, meanwhile, buf can be left modified.
func fixAddr(buf []byte, suffix []byte) (atCharIdx int) {
	atCharIdx = -1
	for i := range buf {
		if buf[i] >= 'A' && buf[i] <= 'Z' {
			buf[i] += 0x20
		}
		if buf[i] == '@' {
			atCharIdx = i
		}
	}
	if !bytes.HasSuffix(buf, suffix) || atCharIdx == len(buf)-1 {
		atCharIdx = -1
	}
	return
}
