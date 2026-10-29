package policy

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Forbidden-name check (onemind-xr docs/design/JADC2-MAP.md section 5.3, last paragraph;
// FEATURES I6): the stop-line names strike, engage, intercept, designate, hostile and the
// CoT strike prefix t-k- must not appear in omos services/schemas. TestForbiddenNames runs it
// in CI (build-policy.yml) against the tree and against a seeded bad schema.

var forbiddenNameRE = regexp.MustCompile(`(?i)\b(strike|engag|intercept|designat|hostile)[a-z]*|t-k-`)

// ForbiddenHit is one occurrence.
type ForbiddenHit struct {
	File string // path relative to the scanned root, forward slashes
	Line int
	Word string // strike | engage | intercept | designate | hostile | t-k-
	Text string
}

func (h ForbiddenHit) String() string {
	return fmt.Sprintf("%s:%d: %s: %s", h.File, h.Line, h.Word, strings.TrimSpace(h.Text))
}

// forbiddenBaseline lists hits that predate the check (file|word). They are owner follow-ups,
// not accepted vocabulary: JADC2-MAP section 6 F1 (the command envelope names cmd.intercept.*
// as a protected class) and the registry example rule 'rule.hostile-drone-critical'. A new
// file or word fails the check.
var forbiddenBaseline = map[string]bool{
	"command/envelope.schema.json|intercept": true,
	"registry/workflow.schema.json|hostile":  true,
}

// ScanForbiddenNames walks root for *.json files and returns every stop-line name.
func ScanForbiddenNames(root string) ([]ForbiddenHit, error) {
	var hits []ForbiddenHit
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		rel = filepath.ToSlash(rel)
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
		for n := 1; sc.Scan(); n++ {
			for _, m := range forbiddenNameRE.FindAllStringSubmatch(sc.Text(), -1) {
				w := strings.ToLower(m[1])
				switch w {
				case "":
					w = "t-k-"
				case "engag":
					w = "engage"
				case "designat":
					w = "designate"
				}
				hits = append(hits, ForbiddenHit{File: rel, Line: n, Word: w, Text: sc.Text()})
			}
		}
		return sc.Err()
	})
	return hits, err
}
