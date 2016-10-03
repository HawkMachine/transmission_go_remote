package transmission

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Remote struct {
	address  string
	username string
	password string
	binPath  string
}

type Torrent struct {
	Id           int
	Name         string
	Hash         string
	Magnet       string
	State        string
	PercentDone  int
	ETASecs      int
	DateAdded    time.Time
	DateFinished time.Time
	DateStarted  time.Time
	LastActivity time.Time
}

func (r *Remote) runCmd(extra_args ...string) (string, error) {
	args := []string{r.address}
	if r.username != "" || r.password != "" {
		args = append(args, "-n", fmt.Sprintf("%s:%s", r.username, r.password))
	}
	args = append(args, extra_args...)
	cmd := exec.Command(r.binPath, args...)
	log.Printf("Running command: %v", cmd)
	bts, err := cmd.Output()
	return string(bts), err
}

func (r *Remote) List(torrent string) ([]*Torrent, error) {
	if torrent == "" {
		torrent = "all"
	}
	output, err := r.runCmd("-t", torrent, "-i")
	if err != nil {
		return nil, err
	}
	torrents := []*Torrent{}
	sections := strings.Split(output, "NAME")[1:] // Skip first empty element.
	var re *regexp.Regexp
	for _, section := range sections {
		t := &Torrent{}

		re = regexp.MustCompile(`(?m)Id: (\d+)`)
		items := re.FindStringSubmatch(section)
		if len(items) == 0 {
			return nil, fmt.Errorf("Failed to match Id")
		}
		id, err := strconv.Atoi(items[1])
		if err != nil {
			return nil, fmt.Errorf("Failed to parse id: %v", err)
		}
		t.Id = id

		re = regexp.MustCompile("(?m)Name: (.*)$")
		items = re.FindStringSubmatch(section)
		if len(items) == 0 {
			return nil, fmt.Errorf("Failed to match Name")
		}
		t.Name = items[1]

		re = regexp.MustCompile("(?m)Hash: (.*)$")
		items = re.FindStringSubmatch(section)
		if len(items) == 0 {
			return nil, fmt.Errorf("Failed to match Hash")
		}
		t.Hash = items[1]

		re = regexp.MustCompile("(?m)Magnet: (.*)$")
		items = re.FindStringSubmatch(section)
		if len(items) == 0 {
			return nil, fmt.Errorf("Failed to match Magnet")
		}
		t.Magnet = items[1]

		re = regexp.MustCompile("(?m)State: (.*)$")
		items = re.FindStringSubmatch(section)
		if len(items) == 0 {
			return nil, fmt.Errorf("Failed to match State")
		}
		t.State = items[1]

		re = regexp.MustCompile(`(?m)Percent Done: (\d+)(\.\d+)?%$`)
		items = re.FindStringSubmatch(section)
		if len(items) > 0 {
			pDone, err := strconv.Atoi(items[1])
			if err == nil {
				t.PercentDone = pDone
			}
		}

		re = regexp.MustCompile(`(?m)ETA:.*\((\d+) seconds\)`)
		items = re.FindStringSubmatch(section)
		if len(items) == 0 {
			return nil, fmt.Errorf("Failed to match ETA: %s", section)
		}
		eta, err := strconv.Atoi(items[1])
		if err != nil {
			return nil, fmt.Errorf("Failed to parse ETA: %v", err)
		}
		t.ETASecs = eta

		torrents = append(torrents, t)
	}
	return torrents, nil
}

func (r *Remote) ListAll() ([]*Torrent, error) {
	return r.List("")
}

func New(address, username, password, transmissionRemotePath string) (*Remote, error) {
	if transmissionRemotePath == "" {
		transmissionRemotePath = "transmission-remote"
	}
	return &Remote{
		address:  address,
		username: username,
		password: password,
		binPath:  transmissionRemotePath,
	}, nil
}
