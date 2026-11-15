package openaichat

import (
	"maps"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/NggKhaiz/FCC/gocore/providerpolicy"
)

// Body keys that carry the output-token budget across OpenAI-compatible
// policies, in the order the reference inspects them.
var outputTokenFields = []string{"max_completion_tokens", "max_tokens"}

var outputTokenFieldSet = map[string]struct{}{
	"max_completion_tokens": {},
	"max_tokens":            {},
}

var structuredTextFields = []string{"message", "detail", "error"}

var structuredContainerFields = []string{"error", "errors", "detail", "details"}

// Python's field grammar has Unicode-aware lookarounds around optional
// quotes. RE2 has no lookaround, so compile each quote combination separately
// and check its two field boundaries after matching. Keeping variants separate
// also lets a match start inside a quote glued to a preceding word, exactly as
// Python's backtracking regex does.
type capPattern struct {
	pattern *regexp.Regexp
	field   int
	cap     int
}

var fieldBoundCapPatterns = buildCapPatterns()

func buildCapPatterns() [][]capPattern {
	quotes := []string{"[`'\"]", ""}
	patterns := make([][]capPattern, 4)
	for _, leading := range quotes {
		for _, trailing := range quotes {
			field := `(?P<field>` + leading + `(?:max_completion_tokens|max_tokens)` + trailing + `)`
			capValue := "[`'\"]?(?P<cap>[0-9]+)[`'\"]?"
			variants := []string{
				`range\s+of\s+` + field + `\s+should\s+be\s+\[\s*1\s*,\s*` + capValue + `\s*\]`,
				field + `\s*(?::\s*)?(?:(?:(?:must|should)\s+be\s+)?(?:less|smaller)\s+than\s+or\s+equal\s+to|(?:(?:must|should)\s+be\s*)?<=|(?:(?:must|should)\s+be\s+)?at\s+most|(?:must|should)\s+not\s+exceed|maximum(?:\s+allowed)?(?:\s+value)?\s+(?:is|of))\s*` + capValue,
				`maximum(?:\s+allowed)?(?:\s+value)?\s+for\s+` + field + `\s+is\s+` + capValue,
				`maximum(?:\s+allowed)?(?:\s+value)?\s+of\s+` + capValue + `\s+for\s+` + field,
			}
			for i, raw := range variants {
				p := regexp.MustCompile(raw)
				patterns[i] = append(patterns[i], capPattern{p, p.SubexpIndex("field"), p.SubexpIndex("cap")})
			}
		}
	}
	return patterns
}

// ParseOutputTokenCap returns the allowed output-token maximum named in a
// 400 rejection, or nil when the error does not bind such a cap.
func ParseOutputTokenCap(err *providerpolicy.ProviderError) *int {
	if !isCapBadRequest(err) {
		return nil
	}
	for _, text := range capErrorTexts(err) {
		if cap := parseCap(text); cap != nil {
			return cap
		}
	}
	return nil
}

// ClampOutputTokens returns a shallow clone with output-token fields
// clamped to cap, or nil when nothing exceeds the cap.
func ClampOutputTokens(body map[string]any, cap int) map[string]any {
	var clamped map[string]any
	for _, field := range outputTokenFields {
		value := asInt(body[field])
		if value != nil && *value > cap {
			if clamped == nil {
				clamped = maps.Clone(body)
			}
			clamped[field] = cap
		}
	}
	return clamped
}

func isCapBadRequest(err *providerpolicy.ProviderError) bool {
	return err.Kind == "BadRequestError" || (err.HasStatus && err.StatusCode == 400)
}

func parseCap(rawText string) *int {
	if !strings.Contains(rawText, "max_tokens") && !strings.Contains(rawText, "max_completion_tokens") {
		return nil
	}
	text := normalizeCapText(rawText)
	hasQuotes := strings.ContainsAny(text, "`'\"")
	for _, variants := range fieldBoundCapPatterns {
		if !hasQuotes {
			// The final variant omits both optional quotes. Most provider
			// errors have no quotes and need only one RE2 search per grammar.
			variants = variants[len(variants)-1:]
		}
		bestStart := len(text) + 1
		bestVariant := len(variants)
		var best []int
		for variant, candidate := range variants {
			for offset := 0; offset < len(text); {
				match := candidate.pattern.FindStringSubmatchIndex(text[offset:])
				if match == nil {
					break
				}
				for i := range match {
					if match[i] >= 0 {
						match[i] += offset
					}
				}
				start, end := match[2*candidate.field], match[2*candidate.field+1]
				if !pythonWordBefore(text, start) && !pythonWordAfter(text, end) {
					if match[0] < bestStart || match[0] == bestStart && variant < bestVariant {
						bestStart, bestVariant, best = match[0], variant, match
					}
					break
				}
				// Advance one rune from the match start, not its end: another
				// valid match may overlap the rejected one.
				_, width := utf8.DecodeRuneInString(text[match[0]:])
				offset = match[0] + width
			}
		}
		if best != nil {
			group := variants[bestVariant].cap
			value, err := strconv.Atoi(text[best[2*group]:best[2*group+1]])
			if err == nil && value > 0 {
				return &value
			}
		}
	}
	return nil
}

// CPython 3.14 uses Unicode 16, while the project's Go 1.22 bootstrap
// uses Unicode 15. These additional alphanumeric ranges are Unicode 16's
// delta, derived from Python str.isalnum() against Go's letter/number tables.
var pythonUnicode16WordAdditions = [][2]rune{
	{0x1c89, 0x1c8a}, {0xa7cb, 0xa7cd}, {0xa7da, 0xa7dc},
	{0x105c0, 0x105f3}, {0x10d40, 0x10d65}, {0x10d6f, 0x10d85},
	{0x10ec2, 0x10ec4}, {0x11380, 0x11389}, {0x1138b, 0x1138b},
	{0x1138e, 0x1138e}, {0x11390, 0x113b5}, {0x113b7, 0x113b7},
	{0x113d1, 0x113d1}, {0x113d3, 0x113d3}, {0x116d0, 0x116e3},
	{0x11bc0, 0x11be0}, {0x11bf0, 0x11bf9}, {0x13460, 0x143fa},
	{0x16100, 0x1611d}, {0x16130, 0x16139}, {0x16d40, 0x16d6c},
	{0x16d70, 0x16d79}, {0x18cff, 0x18cff}, {0x1ccf0, 0x1ccf9},
	{0x1e5d0, 0x1e5ed}, {0x1e5f0, 0x1e5fa}, {0x2ebf0, 0x2ee5d},
}

func pythonWordRune(r rune) bool {
	if r < utf8.RuneSelf {
		return r == '_' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9'
	}
	if unicode.IsLetter(r) || unicode.IsNumber(r) {
		return true
	}
	i := sort.Search(len(pythonUnicode16WordAdditions), func(i int) bool {
		return pythonUnicode16WordAdditions[i][1] >= r
	})
	return i < len(pythonUnicode16WordAdditions) && pythonUnicode16WordAdditions[i][0] <= r
}

func pythonWordBefore(text string, index int) bool {
	if index == 0 {
		return false
	}
	r, _ := utf8.DecodeLastRuneInString(text[:index])
	return pythonWordRune(r)
}

func pythonWordAfter(text string, index int) bool {
	if index == len(text) {
		return false
	}
	r, _ := utf8.DecodeRuneInString(text[index:])
	return pythonWordRune(r)
}

// Python regex \s and \d use Unicode semantics. Normalizing these two
// character classes preserves word boundaries while allowing RE2 to match.
func normalizeCapText(text string) string {
	needsNormalization := false
	for i := 0; i < len(text); i++ {
		c := text[i]
		if c >= utf8.RuneSelf || c == '\v' || c >= 0x1c && c <= 0x1f {
			needsNormalization = true
			break
		}
	}
	if !needsNormalization {
		return text
	}
	var result strings.Builder
	result.Grow(len(text))
	var lastNonSpace rune
	for _, r := range text {
		if unicode.IsSpace(r) || r >= 0x1c && r <= 0x1f {
			result.WriteByte(' ')
			continue
		}
		// In the range grammar the lower bound is the literal ASCII "1",
		// not a \d expression. Keep digits immediately after '[' intact.
		if r >= '0' && r <= '9' {
			result.WriteRune(r)
		} else if digit, ok := decimalDigit(r); ok && lastNonSpace != '[' {
			result.WriteRune('0' + digit)
		} else {
			result.WriteRune(r)
		}
		lastNonSpace = r
	}
	return result.String()
}

func decimalDigit(r rune) (rune, bool) {
	// Unicode 16 decimal digits absent from the Go 1.22 Unicode 15 tables.
	for _, start := range [...]rune{0x10d40, 0x116d0, 0x116da, 0x11bf0, 0x16130, 0x16d70, 0x1ccf0, 0x1e5f1} {
		if r >= start && r < start+10 {
			return r - start, true
		}
	}
	for _, item := range unicode.Digit.R16 {
		if r >= rune(item.Lo) && r <= rune(item.Hi) && (r-rune(item.Lo))%rune(item.Stride) == 0 {
			return ((r - rune(item.Lo)) / rune(item.Stride)) % 10, true
		}
	}
	for _, item := range unicode.Digit.R32 {
		if r >= rune(item.Lo) && r <= rune(item.Hi) && (r-rune(item.Lo))%rune(item.Stride) == 0 {
			return ((r - rune(item.Lo)) / rune(item.Stride)) % 10, true
		}
	}
	return 0, false
}

// capErrorTexts keeps unstructured text separate from structured
// validation messages, matching the reference _error_texts helper.
func capErrorTexts(err *providerpolicy.ProviderError) []string {
	var texts []string
	switch body := err.Body.(type) {
	case nil:
		texts = []string{err.Message}
	case string:
		texts = []string{body}
	case *providerpolicy.JSONObject, map[string]any, []any:
		texts = structuredErrorTexts(err.Body)
	default:
		texts = nil
	}
	lowered := make([]string, len(texts))
	for i, text := range texts {
		// Python's Unicode lowercase expands U+0130 to i + combining dot;
		// Go's simple case mapping drops that dot and changes \w boundaries.
		lowered[i] = strings.ToLower(strings.ReplaceAll(text, "\u0130", "i\u0307"))
	}
	return lowered
}

// structuredErrorTexts walks structured error bodies without letting text
// nested below an explicitly scoped validation error escape that scope.
func structuredErrorTexts(body any) []string {
	texts := make([]string, 0)
	pending := []any{body}
	for len(pending) > 0 {
		value := pending[len(pending)-1]
		pending = pending[:len(pending)-1]

		if items, ok := value.([]any); ok {
			for _, item := range items {
				switch item.(type) {
				case string:
					texts = append(texts, item.(string))
				case *providerpolicy.JSONObject, map[string]any, []any:
					pending = append(pending, item)
				}
			}
			continue
		}

		obj := asObject(value)
		if obj == nil {
			continue
		}
		paramValue := obj["param"]
		outputField := normalizeOutputField(paramValue)
		hasParam := paramValue != nil
		messages := make([]string, 0)
		for _, field := range structuredTextFields {
			if message, ok := obj[field].(string); ok {
				messages = append(messages, message)
			}
		}
		if outputField != "" {
			for _, message := range messages {
				texts = append(texts, outputField+" "+message)
			}
		}
		if hasParam {
			continue
		}
		texts = append(texts, messages...)
		for _, field := range structuredContainerFields {
			switch nested := obj[field].(type) {
			case *providerpolicy.JSONObject, map[string]any, []any:
				pending = append(pending, nested)
			}
		}
	}
	return texts
}

func normalizeOutputField(value any) string {
	text, ok := value.(string)
	if !ok {
		return ""
	}
	field := strings.ToLower(strings.Trim(text, "`'\" "))
	if _, ok := outputTokenFieldSet[field]; ok {
		return field
	}
	return ""
}

// asObject accepts both decoded fixture maps and ordered JSON objects.
func asObject(value any) map[string]any {
	switch typed := value.(type) {
	case map[string]any:
		return typed
	case *providerpolicy.JSONObject:
		return typed.Fields
	default:
		return nil
	}
}
