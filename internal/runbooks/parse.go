package runbooks

import (
	"html"
	"html/template"
	"io/fs"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// Track is one workshop runbook (e.g. 101) parsed from markdown.
type Track struct {
	ID       string
	Title    string
	Intro    string
	Sections []Section
}

// Section is an H2 block (What is Grok Build?, How do I work with an agent?, …).
type Section struct {
	ID    string
	Title string
	Beats []Beat
}

// Beat is an H3 card (title + detail + optional fenced prompt examples).
type Beat struct {
	ID      string
	Title   string
	Detail  string
	Prompts []string
}

// Load reads name (e.g. "101.md") from fsys and parses it.
func Load(fsys fs.FS, name string) (*Track, error) {
	b, err := fs.ReadFile(fsys, name)
	if err != nil {
		return nil, err
	}
	track, err := Parse(string(b))
	if err != nil {
		return nil, err
	}
	track.ID = strings.TrimSuffix(name, ".md")
	return track, nil
}

// Parse turns workshop markdown (H1 / H2 / H3 / fenced prompts) into cards.
func Parse(md string) (*Track, error) {
	md = strings.ReplaceAll(md, "\r\n", "\n")
	lines := strings.Split(md, "\n")
	track := &Track{}
	var section *Section
	var beat *Beat
	var fence []string
	inFence := false
	usedIDs := map[string]int{}

	flushBeat := func() {
		if beat == nil {
			return
		}
		beat.Detail = strings.TrimSpace(beat.Detail)
		if beat.Title == "" && beat.Detail == "" && len(beat.Prompts) == 0 {
			beat = nil
			return
		}
		if section == nil {
			beat = nil
			return
		}
		base := slug(beat.Title)
		if base == "" {
			base = section.ID + "-notes"
		}
		beat.ID = uniqueID(usedIDs, base)
		section.Beats = append(section.Beats, *beat)
		beat = nil
	}

	flushSection := func() {
		flushBeat()
		if section == nil {
			return
		}
		if section.ID == "" {
			section.ID = uniqueID(usedIDs, slug(section.Title))
		}
		track.Sections = append(track.Sections, *section)
		section = nil
	}

	ensureBeat := func() {
		if beat == nil {
			beat = &Beat{}
		}
	}

	for _, line := range lines {
		if inFence {
			if strings.HasPrefix(line, "```") {
				ensureBeat()
				beat.Prompts = append(beat.Prompts, strings.TrimRight(strings.Join(fence, "\n"), "\n"))
				fence = nil
				inFence = false
				continue
			}
			fence = append(fence, line)
			continue
		}
		if strings.HasPrefix(line, "```") {
			inFence = true
			fence = nil
			continue
		}
		if strings.HasPrefix(line, "# ") && track.Title == "" && section == nil {
			track.Title = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			continue
		}
		if strings.HasPrefix(line, "## ") {
			flushSection()
			title := strings.TrimSpace(strings.TrimPrefix(line, "## "))
			section = &Section{Title: title, ID: uniqueID(usedIDs, slug(title))}
			continue
		}
		if strings.HasPrefix(line, "### ") {
			flushBeat()
			title := strings.TrimSpace(strings.TrimPrefix(line, "### "))
			beat = &Beat{Title: title}
			continue
		}
		if section == nil {
			if track.Intro != "" {
				track.Intro += "\n"
			}
			track.Intro += line
			continue
		}
		ensureBeat()
		if beat.Detail != "" {
			beat.Detail += "\n"
		}
		beat.Detail += line
	}
	if inFence && len(fence) > 0 {
		ensureBeat()
		beat.Prompts = append(beat.Prompts, strings.TrimRight(strings.Join(fence, "\n"), "\n"))
	}
	flushSection()
	track.Intro = strings.TrimSpace(track.Intro)
	return track, nil
}

// SafeFileName maps a track id to a markdown filename ("101" → "101.md").
func SafeFileName(id string) (string, bool) {
	if id == "" || strings.Contains(id, "..") {
		return "", false
	}
	for _, r := range id {
		if r > unicode.MaxASCII || !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_') {
			return "", false
		}
	}
	return id + ".md", true
}

var (
	nonSlug = regexp.MustCompile(`[^a-z0-9]+`)
	boldRe  = regexp.MustCompile(`\*\*(.+?)\*\*`)
	codeRe  = regexp.MustCompile("`([^`]+)`")
	emRe    = regexp.MustCompile(`\*(.+?)\*`)
)

func slug(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "—", " ")
	s = strings.ReplaceAll(s, "–", " ")
	s = strings.ReplaceAll(s, "→", " ")
	s = nonSlug.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if len(s) > 56 {
		s = strings.Trim(s[:56], "-")
	}
	return s
}

func uniqueID(used map[string]int, base string) string {
	if base == "" {
		base = "section"
	}
	used[base]++
	if used[base] == 1 {
		return base
	}
	return base + "-" + strconv.Itoa(used[base])
}

// InlineHTML renders a short workshop note: paragraphs, **bold**, `code`, *italic*.
func InlineHTML(s string) template.HTML {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	var b strings.Builder
	for _, para := range strings.Split(s, "\n\n") {
		p := strings.Join(strings.Fields(strings.ReplaceAll(para, "\n", " ")), " ")
		if p == "" {
			continue
		}
		p = html.EscapeString(p)
		p = codeRe.ReplaceAllString(p, "<code>$1</code>")
		p = boldRe.ReplaceAllString(p, "<strong>$1</strong>")
		p = emRe.ReplaceAllString(p, "<em>$1</em>")
		b.WriteString("<p>")
		b.WriteString(p)
		b.WriteString("</p>\n")
	}
	return template.HTML(b.String())
}
