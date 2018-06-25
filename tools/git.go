package tools

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type Git struct {
	Env   *Env
	mutex sync.Mutex
}

func NewGit(env *Env) *Git {
	return &Git{Env: env}
}

func (g *Git) Clone(url string) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	_, err := g.Env.ExecLines("git", "clone", url, ".")
	return err
}

func (g *Git) Fetch() error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	_, err := g.Env.ExecLines("git", "fetch")
	return err
}

func (g *Git) Pull() error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	_, err := g.Env.ExecLines("git", "pull", "--all")
	return err
}

func (g *Git) ListRemotes() (*RemoteList, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	res := &RemoteList{}
	res.References = make(map[string]string)
	lines, err := g.Env.ExecLines("git", "ls-remote")
	for _, l := range lines {
		if strings.HasPrefix(l, "From") {
			res.Origin = l[5:]
		} else {
			hashAndHead := strings.Split(l, "\t")
			res.References[strings.TrimSpace(hashAndHead[1])] = strings.TrimSpace(hashAndHead[0])
		}
	}
	return res, err
}

func (g *Git) ListBranches() (branches []string, e error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	lines, err := g.Env.ExecLines("git", "branch", "-a")
	for i, l := range lines {
		clean := strings.Replace(l, "*", "", -1)
		clean = strings.TrimSpace(clean)
		lines[i] = clean
	}
	return lines, err
}

func (g *Git) Clean(files bool, directories bool) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	cmd := "-"
	if files {
		cmd += "f"
	}
	if directories {
		cmd += "d"
	}
	if len(cmd) == 1 {
		cmd = ""
	}

	_, err := g.Env.ExecLines("git", "clean", cmd)
	return err
}

func (g *Git) Checkout(branch string) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	_, err := g.Env.ExecLines("git", "checkout", branch)
	return err
}

//returns the current head commit hash
func (g *Git) GetHead() (head string, e error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	lines, err := g.Env.ExecLines("git", "rev-parse", "HEAD")
	head = ""
	if len(lines) > 0 {
		head = strings.TrimSpace(lines[0])
	}
	return head, err

}

func (g *Git) ResetHard() error {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	_, err := g.Env.ExecLines("git", "reset", "--hard")
	return err
}

func (g *Git) Show(hash string) (value *GitObject, e error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	format := "%H%n%h%n%ae%n%ce%n%at%n%cD%n%s%n%B"
	lines, err := g.Env.ExecLines("git", "--no-pager", "show", "--quiet", "--pretty=format:"+format, hash)

	if err != nil {
		return nil, err
	} else {

		time, e := strconv.ParseInt(lines[4], 10, 64)
		if e != nil {
			return nil, e
		}
		msg := make([]string, 0)
		for i := 6; i < len(lines); i++ {
			msg = append(msg, lines[i])
		}

		return &GitObject{Hash: hash, CommitHash: lines[0], AuthorEmail: lines[1], CommitterEmail: lines[2], UnixTimestamp: time, CommitterDate: lines[4], Subject: lines[5], Body: msg}, nil
	}
}

func (g *Git) Log(hash string) (values []*GitObject, e error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	format := "%H%n%h%n%ae%n%ce%n%at%n%cD%n%s"
	out, err := g.Env.ExecLines("git", "--no-pager", "log", "--quiet", "--pretty=format:"+format, hash)

	if err != nil {
		return nil, err
	} else {
		res := make([]*GitObject, 0)
		for i := 0; i < len(out); i += 7 {
			time, e := strconv.ParseInt(out[i+4], 10, 64)
			if e != nil {
				return nil, e
			}
			obj := &GitObject{Hash: hash, CommitHash: out[i+0], AbbreviatedCommitHash: out[i+1], AuthorEmail: out[i+2], CommitterEmail: out[i+3], UnixTimestamp: time, CommitterDate: out[i+5], Subject: out[i+6]}
			res = append(res, obj)
		}

		return res, nil
	}
}

//returns all branches which contain the given commit
func (g *Git) GetBranches(hash string) (branches []string, e error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	lines, err := g.Env.ExecLines("git", "branch", "--contains", hash)
	for i, l := range lines {
		clean := strings.Replace(l, "*", "", -1)
		clean = strings.TrimSpace(clean)
		lines[i] = clean
	}
	return lines, err
}

//list all files of a tree-ish
func (g *Git) ListFiles(id string) (files []*GitFile, e error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// git --no-pager ls-tree --long -r master
	out, err := g.Env.ExecLines("git", "--no-pager", "ls-tree", "--long", "-r", id)

	//e.g. <mode> SP <type> SP <object> SP <object size> TAB <file>
	res := make([]*GitFile, 0)
	stopChars := []rune{' ', '\t'}
	for _, line := range out {

		tokens := make([]string, 0)

		text, off := scanUntil(line, 0, stopChars)
		tokens = append(tokens, text)

		text, off = scanUntil(line, off, stopChars)
		tokens = append(tokens, text)

		text, off = scanUntil(line, off, stopChars)
		tokens = append(tokens, text)

		text, off = scanUntil(line, off, stopChars)
		tokens = append(tokens, text)

		text, off = scanUntil(line, off, []rune{})
		tokens = append(tokens, text)

		if len(tokens) != 5 {
			return nil, fmt.Errorf("unsupported line format: " + line)
		}

		mode, e := strconv.ParseInt(tokens[0], 10, 32)
		if e != nil {
			return nil, e
		}

		size, e := strconv.ParseInt(tokens[3], 10, 64)
		if e != nil {
			return nil, e
		}

		res = append(res, &GitFile{Mode: int(mode), Type: tokens[1], Object: tokens[2], Size: size, Name: strings.TrimSpace(tokens[4])})

	}
	return res, err
}

func scanUntil(str string, offsetRune int, stopChars []rune) (string, int) {
	sb := strings.Builder{}
	for i, r := range str {
		if i < offsetRune {
			continue
		}
		isStopChar := false
		for _, stop := range stopChars {
			if r == stop {
				isStopChar = true
				break
			}
		}
		//just consume stop chars until we find something useful
		if sb.Len() == 0 && isStopChar {
			continue
		}
		if !isStopChar {
			sb.WriteRune(r)
		} else {
			//got a stop char, return
			return sb.String(), i
		}
	}
	return sb.String(), len(str)
}

//returns file bytes of an object identified by name
func (g *Git) GetFile(id string, name string) ([]byte, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	//git --no-pager show master:app/app/src/main/res/drawable-xxhdpi/icon.png
	buf, err := g.Env.ExecBytes("git", "--no-pager", "show", id+":"+name)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

//e.g. 100644 blob e7b4def49cb53d9aa04228dd3edb14c9e635e003      15	dir/settings.gradle
type GitFile struct {
	Mode   int
	Type   string
	Object string
	Size   int64
	Name   string
}

func (g *GitFile) String() string {
	return strconv.Itoa(g.Mode) + " - " + g.Type + " - " + g.Object + " - " + strconv.FormatInt(g.Size, 10) + " - " + g.Name
}

type FormatField string

type RemoteList struct {
	//e.g. https://github.com/...
	Origin string
	//e.g. HEAD -> c282fe12fc065314d20571aa7f2a3b3efb6d6324
	References map[string]string
}

type GitObject struct {
	//hash of the actual object
	Hash string

	//%H: commit hash
	CommitHash string

	//%h: abbreviated commit hash
	AbbreviatedCommitHash string

	//%ae: author email
	AuthorEmail string

	//%ce: committer email
	CommitterEmail string

	//%at: author date, UNIX timestamp
	UnixTimestamp int64

	//%cD: committer date, RFC2822 style
	CommitterDate string

	//%s: subject
	Subject string

	//%B: raw body (unwrapped subject and body)
	Body []string
}

type ByGitObjectTimeDesc []*GitObject

func (s ByGitObjectTimeDesc) Len() int {
	return len(s)
}
func (s ByGitObjectTimeDesc) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByGitObjectTimeDesc) Less(i, j int) bool {
	return s[i].UnixTimestamp > s[j].UnixTimestamp
}
