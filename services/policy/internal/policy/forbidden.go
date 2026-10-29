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

// No leading word boundary: compounds such as "airstrike" or "counterintercept" match too.
var forbiddenNameRE = regexp.MustCompile(`(?i)(strike|engag|intercept|designat|hostile)[a-z]*|t-k-`)

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

// forbiddenBaseline counts the hits that predate the check, per file|word. They are owner
// follow-ups, not accepted vocabulary: JADC2-MAP section 6 F1 (the command envelope names
// cmd.intercept.* as a protected class) and the registry example rule
// 'rule.hostile-drone-critical'. Any extra occurrence, including in these files, fails.
var forbiddenBaseline = map[string]int{
	"command/envelope.schema.json|intercept": 1,
	"registry/workflow.schema.json|hostile":  2,
}

// NewForbiddenHits returns the hits beyond the baseline count for their file|word. When a
// file|word exceeds its count, every hit for it is returned, so the reviewer sees them all.
func NewForbiddenHits(hits []ForbiddenHit) []ForbiddenHit {
	n := map[string]int{}
	for _, h := range hits {
		n[h.File+"|"+h.Word]++
	}
	var fresh []ForbiddenHit
	for _, h := range hits {
		if k := h.File + "|" + h.Word; n[k] > forbiddenBaseline[k] {
			fresh = append(fresh, h)
		}
	}
	return fresh
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
