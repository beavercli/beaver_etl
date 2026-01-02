package matcher

import (
	"fmt"
	"path/filepath"

	"github.com/beavercli/beaver_etl/internal/manifest"
	"github.com/bmatcuk/doublestar"
)

// the core responsibilities of the matcher is
// to parse the filesystem of the project and
// return list of matched files.
//
// The matcher should also resolve the overlapping files according to the manifest rules
// currently the resolution is defined by the order of the sources in the manifest
//
// The matcher should explicitly define what field is nil and what is auto
// so the following steps can enirch the data in case of the auto

type Matcher interface {
	MatchFiles(rootRepo string, m *manifest.Manifest) ([]Match, error)
}

type Match struct {
	Path         Path
	Title        *Title
	Tags         []Tag
	Contributors []Contributor
	Link         *Link
	Language     *Language
}
type Path string
type Title string
type Tag string
type Contributor struct {
	Name     string
	LastName string
	Email    string
}
type Link string
type Language string

func MatchFiles(root string, manifestPath string, m *manifest.Manifest) ([]Match, error) {
	absRepo, err := absPath(root)
	if err != nil {
		return nil, err
	}
	absManifest, err := absPath(manifestPath)
	if err != nil {
		return nil, err
	}
	manifestDir := filepath.Dir(absManifest)

	matches := make([]Match, 0, len(m.Sources))
	for _, source := range m.Sources {
		switch source.Type {
		case manifest.File:
			fs := source.Spec.(*manifest.FileSource)
			match, err := fileToMatch(absRepo, manifestDir, fs)
			if err != nil {
				return nil, err
			}
			matches = append(matches, match)
			break
		case manifest.Pattern:
			ps := source.Spec.(*manifest.PatternSource)
			ms, err := patternToMatches(absRepo, manifestDir, ps)
			if err != nil {
				return nil, err
			}
			matches = append(matches, ms...)
			break
		default:
			return nil, fmt.Errorf("Source type %s is not supported", source.Type)
		}
	}

	return nil, nil
}

func toTags(tags []manifest.Tag) []Tag {
	ts := make([]Tag, len(tags))
	for i, t := range tags {
		ts[i] = Tag(t)
	}
	return ts
}

func toContributor(c manifest.Contributor) Contributor {
	return Contributor{
		Name:     string(c.Name),
		LastName: string(c.LastName),
		Email:    string(c.Email),
	}
}

func toContributors(contribs []manifest.Contributor) []Contributor {
	cs := make([]Contributor, len(contribs))
	for i, c := range contribs {
		cs[i] = toContributor(c)
	}
	return cs
}
func fileToMatch(repo string, dir string, fs *manifest.FileSource) (Match, error) {
	path, err := relToPath(repo, joinPath(dir, string(fs.Path)))
	if err != nil {
		return Match{}, err
	}

	return Match{
		Path:         Path(path),
		Title:        (*Title)(fs.Title),
		Link:         (*Link)(fs.Link),
		Language:     (*Language)(fs.Language),
		Tags:         toTags(fs.Tags),
		Contributors: toContributors(fs.Contributors),
	}, nil
}

func patternToMatches(repo string, dir string, ps *manifest.PatternSource) ([]Match, error) {
	// collect all candidates in Include
	candidates := make([]string, 0, len(ps.Include))
	for _, p := range ps.Include {
		pattern := joinPath(dir, filepath.FromSlash(string(p)))

		matches, err := doublestar.Glob(pattern)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, matches...)
	}

	excludePatterns := make([]string, len(ps.Exclude))
	for i, p := range ps.Exclude {
		excludePatterns[i] = joinPath(dir, filepath.FromSlash(string(p)))
	}

	// filter the collect candidates with the Exclude patterns
	paths := make([]string, 0, len(candidates))
	for _, c := range candidates {
		m, err := matchAny(c, excludePatterns)
		if err != nil {
			return nil, err
		}

		if m {
			continue
		}

		paths = append(paths, c)
	}
	matches := make([]Match, len(paths))
	for i, p := range paths {
		path, err := relToPath(repo, p)

		if err != nil {
			return nil, err
		}

		matches[i] = Match{
			Path: Path(path),
			// TODO:
		}
	}

	return nil, nil
}

func matchAny(c string, patterns []string) (bool, error) {
	for _, p := range patterns {
		m, err := doublestar.Match(p, c)
		if err != nil {
			return false, err
		}
		if m {
			return true, nil
		}
	}
	return false, nil
}
