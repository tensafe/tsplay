package tsplay_core

import (
	"regexp"
	"sort"
	"strings"
)

type PageObservationSelector struct {
	Value  string `json:"value"`
	Kind   string `json:"kind,omitempty"`
	Score  int    `json:"score,omitempty"`
	Reason string `json:"reason,omitempty"`
}

var observedSelectorLongDigitsPattern = regexp.MustCompile(`\d{3,}`)
var observedSelectorHexPattern = regexp.MustCompile(`[0-9a-f]{6,}`)

func normalizeObservedSelectorDiagnostics(element *PageObservationElement) {
	if element == nil {
		return
	}
	diagnostics := observedSelectorDiagnostics(*element)
	if len(diagnostics) == 0 {
		element.PrimarySelector = ""
		element.SelectorRationale = ""
		element.SelectorDetails = nil
		return
	}
	element.SelectorDetails = diagnostics
	element.SelectorCandidates = make([]string, 0, len(diagnostics))
	for _, diagnostic := range diagnostics {
		element.SelectorCandidates = append(element.SelectorCandidates, diagnostic.Value)
	}
	element.PrimarySelector = diagnostics[0].Value
	element.SelectorRationale = diagnostics[0].Reason
}

func preferredObservedSelector(element PageObservationElement) string {
	diagnostics := observedSelectorDiagnostics(element)
	if len(diagnostics) == 0 {
		return ""
	}
	return diagnostics[0].Value
}

func observedSelectorDiagnostics(element PageObservationElement) []PageObservationSelector {
	if len(element.SelectorCandidates) == 0 {
		return nil
	}

	seen := map[string]bool{}
	diagnostics := make([]PageObservationSelector, 0, len(element.SelectorCandidates))
	for _, selector := range element.SelectorCandidates {
		selector = strings.TrimSpace(selector)
		if selector == "" || seen[selector] {
			continue
		}
		seen[selector] = true
		diagnostics = append(diagnostics, PageObservationSelector{
			Value:  selector,
			Kind:   observedSelectorKind(selector),
			Score:  scoreObservedSelectorCandidateForElement(element, selector),
			Reason: observedSelectorReason(selector),
		})
	}

	sort.SliceStable(diagnostics, func(i, j int) bool {
		if diagnostics[i].Score != diagnostics[j].Score {
			return diagnostics[i].Score > diagnostics[j].Score
		}
		if diagnostics[i].Kind != diagnostics[j].Kind {
			return diagnostics[i].Kind < diagnostics[j].Kind
		}
		return diagnostics[i].Value < diagnostics[j].Value
	})
	return diagnostics
}

func scoreObservedSelectorCandidate(selector string) int {
	return scoreObservedSelectorCandidateForElement(PageObservationElement{}, selector)
}

func scoreObservedSelectorCandidateForElement(element PageObservationElement, selector string) int {
	kind := observedSelectorKind(selector)
	score := 40
	switch kind {
	case "data-testid":
		score = 100
	case "data-test":
		score = 98
	case "data-cy":
		score = 96
	case "tag-data-testid":
		score = 95
	case "tag-data-test":
		score = 93
	case "tag-data-cy":
		score = 91
	case "tag-id":
		score = 92
	case "id":
		score = 90
	case "name+type":
		score = 86
	case "name":
		score = 83
	case "href":
		score = 80
	case "tag-aria-label":
		score = 78
	case "aria-label":
		score = 76
	case "role":
		score = 74
	case "placeholder":
		score = 70
	case "text":
		score = 64
	case "xpath":
		score = 10
	}

	if selectorLooksGeneratedID(selector) {
		score -= 18
	}
	if kind == "text" && len([]rune(observedSelectorQuotedValue(selector))) > 48 {
		score -= 10
	}
	if kind == "href" && selectorHasDynamicHref(selector) {
		score -= 14
	}
	if kind == "placeholder" && strings.TrimSpace(element.Label) != "" {
		score -= 4
	}
	if score < 0 {
		return 0
	}
	return score
}

func observedSelectorKind(selector string) string {
	switch {
	case strings.HasPrefix(selector, `[`+"data-testid="):
		return "data-testid"
	case strings.HasPrefix(selector, `[`+"data-test="):
		return "data-test"
	case strings.HasPrefix(selector, `[`+"data-cy="):
		return "data-cy"
	case strings.Contains(selector, `[`+"data-testid="):
		return "tag-data-testid"
	case strings.Contains(selector, `[`+"data-test="):
		return "tag-data-test"
	case strings.Contains(selector, `[`+"data-cy="):
		return "tag-data-cy"
	case strings.HasPrefix(selector, "#"):
		return "id"
	case strings.Contains(selector, "#"):
		return "tag-id"
	case strings.Contains(selector, "[name=") && strings.Contains(selector, "[type="):
		return "name+type"
	case strings.Contains(selector, "[name="):
		return "name"
	case strings.HasPrefix(selector, "a[href="):
		return "href"
	case strings.HasPrefix(selector, "role="):
		return "role"
	case strings.Contains(selector, "[placeholder="):
		return "placeholder"
	case strings.HasPrefix(selector, "[aria-label="):
		return "aria-label"
	case strings.Contains(selector, "[aria-label="):
		return "tag-aria-label"
	case strings.HasPrefix(selector, "text="):
		return "text"
	case strings.HasPrefix(selector, "xpath="):
		return "xpath"
	default:
		return "other"
	}
}

func observedSelectorReason(selector string) string {
	switch observedSelectorKind(selector) {
	case "data-testid", "data-test", "data-cy", "tag-data-testid", "tag-data-test", "tag-data-cy":
		return "Uses an explicit test attribute, which is usually the most stable selector across layout and text changes."
	case "id", "tag-id":
		if selectorLooksGeneratedID(selector) {
			return "Uses an id, but it looks generated, so treat it as a fallback rather than a long-term stable selector."
		}
		return "Uses a stable id, which is typically stronger than visible-text selectors."
	case "name", "name+type":
		return "Targets a named form field, which is usually stable for inputs and selects."
	case "href":
		return "Targets a link destination, which can stay stable when the href is static."
	case "aria-label", "tag-aria-label":
		return "Uses an accessibility label, which stays readable and resilient when layout changes."
	case "role":
		return "Uses accessible role and name, which is readable and often more stable than raw text."
	case "placeholder":
		return "Falls back to placeholder text because no stronger test attribute or stable id was available."
	case "text":
		return "Falls back to visible text because no stronger stable attribute was found."
	case "xpath":
		return "XPath is kept only as a last resort when stronger selectors are unavailable."
	default:
		return "Fallback selector candidate kept for manual review."
	}
}

func selectorLooksGeneratedID(selector string) bool {
	id := observedSelectorIDValue(selector)
	if id == "" {
		return false
	}
	if len(id) > 40 || strings.ContainsAny(id, ": ") {
		return true
	}
	lower := strings.ToLower(id)
	for _, marker := range []string{"react", "radix", "headlessui", "chakra", "mantine", "mui", "ember", "auto", "generated"} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return observedSelectorLongDigitsPattern.MatchString(lower) || observedSelectorHexPattern.MatchString(lower)
}

func selectorHasDynamicHref(selector string) bool {
	value := strings.ToLower(observedSelectorQuotedValue(selector))
	return strings.Contains(value, "?") || strings.Contains(value, "#") || observedSelectorLongDigitsPattern.MatchString(value)
}

func observedSelectorIDValue(selector string) string {
	switch {
	case strings.HasPrefix(selector, "#"):
		return strings.TrimPrefix(selector, "#")
	case strings.Contains(selector, "#"):
		parts := strings.SplitN(selector, "#", 2)
		if len(parts) == 2 {
			return parts[1]
		}
	}
	return ""
}

func observedSelectorQuotedValue(selector string) string {
	start := strings.IndexByte(selector, '"')
	end := strings.LastIndexByte(selector, '"')
	if start >= 0 && end > start {
		return selector[start+1 : end]
	}
	start = strings.IndexByte(selector, '\'')
	end = strings.LastIndexByte(selector, '\'')
	if start >= 0 && end > start {
		return selector[start+1 : end]
	}
	return ""
}
