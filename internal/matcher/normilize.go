package matcher

import "path/filepath"

func absPath(repoRoot string) (string, error) {
	abs, err := filepath.Abs(repoRoot)
	if err != nil {
		return "", err
	}

	abs = filepath.Clean(abs)
	abs = filepath.ToSlash(abs)

	return abs, nil
}

func relToPath(base string, path string) (string, error) {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return "", err
	}

	rel = filepath.Clean(rel)
	rel = filepath.ToSlash(rel)

	return rel, nil
}

func joinPath(dir string, path string) string {
	p := filepath.Join(dir, filepath.Clean(path))
	return p
}
