package libs

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Scrambler -
func Scrambler(src string) string {
	ln := len(src)
	var temp []string
	for _, c := range src {
		temp = append(temp, fmt.Sprintf("%c", c))
	}
	perm := rand.Perm(ln)
	res := make([]string, ln)
	for i, v := range perm {
		res[v] = temp[i]
	}

	return strings.Join(res, "")
}

// GuuID -
func GuuID(data string) string {
	u := uuid.NewMD5(uuid.NameSpaceDNS, []byte(data+time.Now().String()))
	d, _ := u.MarshalText()
	return string(d)
}
