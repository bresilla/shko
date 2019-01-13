//source https://github.com/zyedidia/highlight

package shko

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	yaml "gopkg.in/yaml.v2"
)

// DetectFiletype will use the list of syntax definitions provided and the filename and first line of the file
// to determine the filetype of the file
// It will return the corresponding syntax definition for the filetype
func DetectFiletype(defs []*Def, filename string, firstLine []byte) *Def {
	for _, d := range defs {
		if d.ftdetect[0].MatchString(filename) {
			return d
		}
		if len(d.ftdetect) > 1 {
			if d.ftdetect[1].MatchString(string(firstLine)) {
				return d
			}
		}
	}

	emptyDef := new(Def)
	emptyDef.FileType = "Unknown"
	emptyDef.rules = new(rules)
	return emptyDef
}

// A Group represents a syntax group
type Group uint8

// Groups contains all of the groups that are defined
// You can access them in the map via their string name
var Groups map[string]Group
var numGroups Group

// String returns the group name attached to the specific group
func (g Group) String() string {
	for k, v := range Groups {
		if v == g {
			return k
		}
	}
	return ""
}

// A Def is a full syntax definition for a language
// It has a filetype, information about how to detect the filetype based
// on filename or header (the first line of the file)
// Then it has the rules which define how to highlight the file
type Def struct {
	FileType string
	ftdetect []*regexp.Regexp
	rules    *rules
}

// A Pattern is one simple syntax rule
// It has a group that the rule belongs to, as well as
// the regular expression to match the pattern
type pattern struct {
	group Group
	regex *regexp.Regexp
}

// rules defines which patterns and regions can be used to highlight
// a filetype
type rules struct {
	regions  []*region
	patterns []*pattern
	includes []string
}

// A region is a highlighted region (such as a multiline comment, or a string)
// It belongs to a group, and has start and end regular expressions
// A region also has rules of its own that only apply when matching inside the
// region and also rules from the above region do not match inside this region
// Note that a region may contain more regions
type region struct {
	group      Group
	limitGroup Group
	parent     *region
	start      *regexp.Regexp
	end        *regexp.Regexp
	skip       *regexp.Regexp
	rules      *rules
}

func init() {
	Groups = make(map[string]Group)
}

// ParseDef parses an input syntax file into a highlight Def
func ParseDef(input []byte) (s *Def, err error) {
	// This is just so if we have an error, we can exit cleanly and return the parse error to the user
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	var rules map[interface{}]interface{}
	if err = yaml.Unmarshal(input, &rules); err != nil {
		return nil, err
	}

	s = new(Def)

	for k, v := range rules {
		if k == "filetype" {
			filetype := v.(string)

			s.FileType = filetype
		} else if k == "detect" {
			ftdetect := v.(map[interface{}]interface{})
			if len(ftdetect) >= 1 {
				syntax, err := regexp.Compile(ftdetect["filename"].(string))
				if err != nil {
					return nil, err
				}

				s.ftdetect = append(s.ftdetect, syntax)
			}
			if len(ftdetect) >= 2 {
				header, err := regexp.Compile(ftdetect["header"].(string))
				if err != nil {
					return nil, err
				}

				s.ftdetect = append(s.ftdetect, header)
			}
		} else if k == "rules" {
			inputRules := v.([]interface{})

			rules, err := parseRules(inputRules, nil)
			if err != nil {
				return nil, err
			}

			s.rules = rules
		}
	}

	return s, err
}

// ResolveIncludes will sort out the rules for including other filetypes
// You should call this after parsing all the Defs
func ResolveIncludes(defs []*Def) {
	for _, d := range defs {
		resolveIncludesInDef(defs, d)
	}
}

func resolveIncludesInDef(defs []*Def, d *Def) {
	for _, lang := range d.rules.includes {
		for _, searchDef := range defs {
			if lang == searchDef.FileType {
				d.rules.patterns = append(d.rules.patterns, searchDef.rules.patterns...)
				d.rules.regions = append(d.rules.regions, searchDef.rules.regions...)
			}
		}
	}
	for _, r := range d.rules.regions {
		resolveIncludesInRegion(defs, r)
		r.parent = nil
	}
}

func resolveIncludesInRegion(defs []*Def, region *region) {
	for _, lang := range region.rules.includes {
		for _, searchDef := range defs {
			if lang == searchDef.FileType {
				region.rules.patterns = append(region.rules.patterns, searchDef.rules.patterns...)
				region.rules.regions = append(region.rules.regions, searchDef.rules.regions...)
			}
		}
	}
	for _, r := range region.rules.regions {
		resolveIncludesInRegion(defs, r)
		r.parent = region
	}
}

func parseRules(input []interface{}, curRegion *region) (*rules, error) {
	rules := new(rules)

	for _, v := range input {
		rule := v.(map[interface{}]interface{})
		for k, val := range rule {
			group := k

			switch object := val.(type) {
			case string:
				if k == "include" {
					rules.includes = append(rules.includes, object)
				} else {
					// Pattern
					r, err := regexp.Compile(object)
					if err != nil {
						return nil, err
					}

					groupStr := group.(string)
					if _, ok := Groups[groupStr]; !ok {
						numGroups++
						Groups[groupStr] = numGroups
					}
					groupNum := Groups[groupStr]
					rules.patterns = append(rules.patterns, &pattern{groupNum, r})
				}
			case map[interface{}]interface{}:
				// region
				region, err := parseRegion(group.(string), object, curRegion)
				if err != nil {
					return nil, err
				}
				rules.regions = append(rules.regions, region)
			default:
				return nil, fmt.Errorf("Bad type %T", object)
			}
		}
	}

	return rules, nil
}

func parseRegion(group string, regionInfo map[interface{}]interface{}, prevRegion *region) (*region, error) {
	var err error

	region := new(region)
	if _, ok := Groups[group]; !ok {
		numGroups++
		Groups[group] = numGroups
	}
	groupNum := Groups[group]
	region.group = groupNum
	region.parent = prevRegion

	region.start, err = regexp.Compile(regionInfo["start"].(string))

	if err != nil {
		return nil, err
	}

	region.end, err = regexp.Compile(regionInfo["end"].(string))

	if err != nil {
		return nil, err
	}

	// skip is optional
	if _, ok := regionInfo["skip"]; ok {
		region.skip, err = regexp.Compile(regionInfo["skip"].(string))

		if err != nil {
			return nil, err
		}
	}

	// limit-color is optional
	if _, ok := regionInfo["limit-group"]; ok {
		groupStr := regionInfo["limit-group"].(string)
		if _, ok := Groups[groupStr]; !ok {
			numGroups++
			Groups[groupStr] = numGroups
		}
		groupNum := Groups[groupStr]
		region.limitGroup = groupNum

		if err != nil {
			return nil, err
		}
	} else {
		region.limitGroup = region.group
	}

	region.rules, err = parseRules(regionInfo["rules"].([]interface{}), region)

	if err != nil {
		return nil, err
	}

	return region, nil
}

// RunePos returns the rune index of a given byte index
// This could cause problems if the byte index is between code points
func runePos(p int, str string) int {
	if p < 0 {
		return 0
	}
	if p >= len(str) {
		return utf8.RuneCountInString(str)
	}
	return utf8.RuneCountInString(str[:p])
}

func combineLineMatch(src, dst LineMatch) LineMatch {
	for k, v := range src {
		if g, ok := dst[k]; ok {
			if g == 0 {
				dst[k] = v
			}
		} else {
			dst[k] = v
		}
	}
	return dst
}

// A State represents the region at the end of a line
type State *region

// LineStates is an interface for a buffer-like object which can also store the states and matches for every line
type LineStates interface {
	Line(n int) string
	LinesNum() int
	State(lineN int) State
	SetState(lineN int, s State)
	SetMatch(lineN int, m LineMatch)
}

// A Highlighter contains the information needed to highlight a string
type Highlighter struct {
	lastRegion *region
	Def        *Def
}

// NewHighlighter returns a new highlighter from the given syntax definition
func NewHighlighter(def *Def) *Highlighter {
	h := new(Highlighter)
	h.Def = def
	return h
}

// LineMatch represents the syntax highlighting matches for one line. Each index where the coloring is changed is marked with that
// color's group (represented as one byte)
type LineMatch map[int]Group

func findIndex(regex *regexp.Regexp, skip *regexp.Regexp, str []rune, canMatchStart, canMatchEnd bool) []int {
	regexStr := regex.String()
	if strings.Contains(regexStr, "^") {
		if !canMatchStart {
			return nil
		}
	}
	if strings.Contains(regexStr, "$") {
		if !canMatchEnd {
			return nil
		}
	}

	var strbytes []byte
	if skip != nil {
		strbytes = skip.ReplaceAllFunc([]byte(string(str)), func(match []byte) []byte {
			res := make([]byte, utf8.RuneCount(match))
			return res
		})
	} else {
		strbytes = []byte(string(str))
	}

	match := regex.FindIndex(strbytes)
	if match == nil {
		return nil
	}
	// return []int{match.Index, match.Index + match.Length}
	return []int{runePos(match[0], string(str)), runePos(match[1], string(str))}
}

func findAllIndex(regex *regexp.Regexp, str []rune, canMatchStart, canMatchEnd bool) [][]int {
	regexStr := regex.String()
	if strings.Contains(regexStr, "^") {
		if !canMatchStart {
			return nil
		}
	}
	if strings.Contains(regexStr, "$") {
		if !canMatchEnd {
			return nil
		}
	}
	matches := regex.FindAllIndex([]byte(string(str)), -1)
	for i, m := range matches {
		matches[i][0] = runePos(m[0], string(str))
		matches[i][1] = runePos(m[1], string(str))
	}
	return matches
}

func (h *Highlighter) highlightRegion(highlights LineMatch, start int, canMatchEnd bool, lineNum int, line []rune, curRegion *region, statesOnly bool) LineMatch {
	// highlights := make(LineMatch)

	if start == 0 {
		if !statesOnly {
			if _, ok := highlights[0]; !ok {
				highlights[0] = curRegion.group
			}
		}
	}

	loc := findIndex(curRegion.end, curRegion.skip, line, start == 0, canMatchEnd)
	if loc != nil {
		if !statesOnly {
			highlights[start+loc[0]] = curRegion.limitGroup
		}
		if curRegion.parent == nil {
			if !statesOnly {
				highlights[start+loc[1]] = 0
				h.highlightRegion(highlights, start, false, lineNum, line[:loc[0]], curRegion, statesOnly)
			}
			h.highlightEmptyRegion(highlights, start+loc[1], canMatchEnd, lineNum, line[loc[1]:], statesOnly)
			return highlights
		}
		if !statesOnly {
			highlights[start+loc[1]] = curRegion.parent.group
			h.highlightRegion(highlights, start, false, lineNum, line[:loc[0]], curRegion, statesOnly)
		}
		h.highlightRegion(highlights, start+loc[1], canMatchEnd, lineNum, line[loc[1]:], curRegion.parent, statesOnly)
		return highlights
	}

	if len(line) == 0 || statesOnly {
		if canMatchEnd {
			h.lastRegion = curRegion
		}

		return highlights
	}

	firstLoc := []int{len(line), 0}

	var firstRegion *region
	for _, r := range curRegion.rules.regions {
		loc := findIndex(r.start, nil, line, start == 0, canMatchEnd)
		if loc != nil {
			if loc[0] < firstLoc[0] {
				firstLoc = loc
				firstRegion = r
			}
		}
	}
	if firstLoc[0] != len(line) {
		highlights[start+firstLoc[0]] = firstRegion.limitGroup
		h.highlightRegion(highlights, start, false, lineNum, line[:firstLoc[0]], curRegion, statesOnly)
		h.highlightRegion(highlights, start+firstLoc[1], canMatchEnd, lineNum, line[firstLoc[1]:], firstRegion, statesOnly)
		return highlights
	}

	fullHighlights := make([]Group, len([]rune(string(line))))
	for i := 0; i < len(fullHighlights); i++ {
		fullHighlights[i] = curRegion.group
	}

	for _, p := range curRegion.rules.patterns {
		matches := findAllIndex(p.regex, line, start == 0, canMatchEnd)
		for _, m := range matches {
			for i := m[0]; i < m[1]; i++ {
				fullHighlights[i] = p.group
			}
		}
	}
	for i, h := range fullHighlights {
		if i == 0 || h != fullHighlights[i-1] {
			// if _, ok := highlights[start+i]; !ok {
			highlights[start+i] = h
			// }
		}
	}

	if canMatchEnd {
		h.lastRegion = curRegion
	}

	return highlights
}

func (h *Highlighter) highlightEmptyRegion(highlights LineMatch, start int, canMatchEnd bool, lineNum int, line []rune, statesOnly bool) LineMatch {
	if len(line) == 0 {
		if canMatchEnd {
			h.lastRegion = nil
		}
		return highlights
	}

	firstLoc := []int{len(line), 0}
	var firstRegion *region
	for _, r := range h.Def.rules.regions {
		loc := findIndex(r.start, nil, line, start == 0, canMatchEnd)
		if loc != nil {
			if loc[0] < firstLoc[0] {
				firstLoc = loc
				firstRegion = r
			}
		}
	}
	if firstLoc[0] != len(line) {
		if !statesOnly {
			highlights[start+firstLoc[0]] = firstRegion.limitGroup
		}
		h.highlightEmptyRegion(highlights, start, false, lineNum, line[:firstLoc[0]], statesOnly)
		h.highlightRegion(highlights, start+firstLoc[1], canMatchEnd, lineNum, line[firstLoc[1]:], firstRegion, statesOnly)
		return highlights
	}

	if statesOnly {
		if canMatchEnd {
			h.lastRegion = nil
		}

		return highlights
	}

	fullHighlights := make([]Group, len(line))
	for _, p := range h.Def.rules.patterns {
		matches := findAllIndex(p.regex, line, start == 0, canMatchEnd)
		for _, m := range matches {
			for i := m[0]; i < m[1]; i++ {
				fullHighlights[i] = p.group
			}
		}
	}
	for i, h := range fullHighlights {
		if i == 0 || h != fullHighlights[i-1] {
			// if _, ok := highlights[start+i]; !ok {
			highlights[start+i] = h
			// }
		}
	}

	if canMatchEnd {
		h.lastRegion = nil
	}

	return highlights
}

// HighlightString syntax highlights a string
// Use this function for simple syntax highlighting and use the other functions for
// more advanced syntax highlighting. They are optimized for quick rehighlighting of the same
// text with minor changes made
func (h *Highlighter) HighlightString(input string) []LineMatch {
	lines := strings.Split(input, "\n")
	var lineMatches []LineMatch

	for i := 0; i < len(lines); i++ {
		line := []rune(lines[i])
		highlights := make(LineMatch)

		if i == 0 || h.lastRegion == nil {
			lineMatches = append(lineMatches, h.highlightEmptyRegion(highlights, 0, true, i, line, false))
		} else {
			lineMatches = append(lineMatches, h.highlightRegion(highlights, 0, true, i, line, h.lastRegion, false))
		}
	}

	return lineMatches
}

// HighlightStates correctly sets all states for the buffer
func (h *Highlighter) HighlightStates(input LineStates) {
	for i := 0; i < input.LinesNum(); i++ {
		line := []rune(input.Line(i))
		// highlights := make(LineMatch)

		if i == 0 || h.lastRegion == nil {
			h.highlightEmptyRegion(nil, 0, true, i, line, true)
		} else {
			h.highlightRegion(nil, 0, true, i, line, h.lastRegion, true)
		}

		curState := h.lastRegion

		input.SetState(i, curState)
	}
}

// HighlightMatches sets the matches for each line in between startline and endline
// It sets all other matches in the buffer to nil to conserve memory
// This assumes that all the states are set correctly
func (h *Highlighter) HighlightMatches(input LineStates, startline, endline int) {
	for i := startline; i < endline; i++ {
		if i >= input.LinesNum() {
			break
		}

		line := []rune(input.Line(i))
		highlights := make(LineMatch)

		var match LineMatch
		if i == 0 || input.State(i-1) == nil {
			match = h.highlightEmptyRegion(highlights, 0, true, i, line, false)
		} else {
			match = h.highlightRegion(highlights, 0, true, i, line, input.State(i-1), false)
		}

		input.SetMatch(i, match)
	}
}

// ReHighlightStates will scan down from `startline` and set the appropriate end of line state
// for each line until it comes across the same state in two consecutive lines
func (h *Highlighter) ReHighlightStates(input LineStates, startline int) {
	// lines := input.LineData()

	h.lastRegion = nil
	if startline > 0 {
		h.lastRegion = input.State(startline - 1)
	}
	for i := startline; i < input.LinesNum(); i++ {
		line := []rune(input.Line(i))
		// highlights := make(LineMatch)

		// var match LineMatch
		if i == 0 || h.lastRegion == nil {
			h.highlightEmptyRegion(nil, 0, true, i, line, true)
		} else {
			h.highlightRegion(nil, 0, true, i, line, h.lastRegion, true)
		}
		curState := h.lastRegion
		lastState := input.State(i)

		input.SetState(i, curState)

		if curState == lastState {
			break
		}
	}
}

// ReHighlightLine will rehighlight the state and match for a single line
func (h *Highlighter) ReHighlightLine(input LineStates, lineN int) {
	line := []rune(input.Line(lineN))
	highlights := make(LineMatch)

	h.lastRegion = nil
	if lineN > 0 {
		h.lastRegion = input.State(lineN - 1)
	}

	var match LineMatch
	if lineN == 0 || h.lastRegion == nil {
		match = h.highlightEmptyRegion(highlights, 0, true, lineN, line, false)
	} else {
		match = h.highlightRegion(highlights, 0, true, lineN, line, h.lastRegion, false)
	}
	curState := h.lastRegion

	input.SetMatch(lineN, match)
	input.SetState(lineN, curState)
}
