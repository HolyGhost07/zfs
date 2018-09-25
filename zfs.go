package zfs

import (
	"errors"
	"io"
	"strings"

	"github.com/theairkit/runcmd"
)

var (
	FS      = "filesystem"
	SNAP    = "snapshot"
	DATANOE = "dataset does not exist"
)

type Zfs struct {
	runcmd.Runner
}

var std, _ = NewZfs(runcmd.NewLocalRunner())

func NewZfs(r runcmd.Runner, err error) (*Zfs, error) {
	if err != nil {
		return nil, err
	}
	return &Zfs{r}, nil
}

func CreateSnap(fs, snap string) error {
	return std.CreateSnap(fs, snap)
}

func (this *Zfs) CreateSnap(fs, snap string) error {
	c := this.Command("zfs snapshot " + fs + "@" + snap)
	err := c.CmdError()
	if err != nil {
		return err
	}
	_, err = c.Run()
	return err
}

func CreateFs(fs string) error {
	return std.CreateFs(fs)
}

func (this *Zfs) CreateFs(fs string) error {
	c := this.Command("zfs create " + fs)
	err := c.CmdError()
	if err != nil {
		return err
	}
	_, err = c.Run()
	return err
}

func Destroy(fs string) error {
	return std.Destroy(fs)
}

func (this *Zfs) Destroy(fs string) error {
	c := this.Command("zfs destroy " + fs)
	err := c.CmdError()
	if err != nil {
		return err
	}
	_, err = c.Run()
	return err
}

func RenameSnap(fs, snapOld, snapNew string) error {
	return std.RenameSnap(fs, snapOld, snapNew)
}

func (this *Zfs) RenameSnap(fs, snapOld, snapNew string) error {
	c := this.Command(
		"zfs rename " + fs + "@" + snapOld + " " + fs + "@" + snapNew,
	)
	err := c.CmdError()
	if err != nil {
		return err
	}
	_, err = c.Run()
	return err
}

func ExistFs(fs string) (bool, error) {
	return std.ExistFs(fs)
}

func (this *Zfs) ExistFs(fs string) (bool, error) {
	if _, err := this.List(fs, FS, false); err != nil {
		if strings.Contains(err.Error(), DATANOE) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func ExistSnap(fs, snap string) (bool, error) {
	return std.ExistSnap(fs, snap)
}

func (this *Zfs) ExistSnap(fs, snap string) (bool, error) {
	if _, err := this.List(fs+"@"+snap, SNAP, false); err != nil {
		if strings.Contains(err.Error(), DATANOE) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func List(fs, fsType string, recursive bool) ([]string, error) {
	return std.List(fs, fsType, recursive)
}

func (this *Zfs) List(fs, fsType string, recursive bool) ([]string, error) {
	// List fs by mask: zroot/blah*
	// Get all fs by recursive call List(), and return matches:
	if strings.HasSuffix(fs, "*") {
		list := make([]string, 0)
		out, err := this.List("", FS, false)
		if err != nil {
			return nil, err
		}
		for _, next := range out {
			if strings.Contains(next, strings.TrimRight(fs, "*")) {
				list = append(list, next)
			}
		}
		return list, nil
	}

	r := ""
	if recursive {
		r = "-r"
	}
	cmd := "zfs list -Ho name -t " + fsType + " " + r + " " + fs
	c := this.Command(cmd)
	err := c.CmdError()
	if err != nil {
		return nil, err
	}
	return c.Run()
}

func ListFsSnap(fs string) ([]string, error) {
	return std.ListFsSnap(fs)
}

func (this *Zfs) ListFsSnap(fs string) ([]string, error) {
	cmd := "zfs list -Ho name -d1 -t snapshot " + fs
	c := this.Command(cmd)
	err := c.CmdError()
	if err != nil {
		return nil, err
	}
	return c.Run()
}

func Property(fs, property string) (string, error) {
	return std.Property(fs, property)
}

func (this *Zfs) Property(fs, property string) (string, error) {
	c := this.Command("zfs get -H -o value " + property + " " + fs)
	err := c.CmdError()
	if err != nil {
		return "", err
	}
	out, err := c.Run()
	if err != nil {
		return "", err
	}
	return out[0], nil
}

func SetProperty(fs, property, value string) error {
	return std.SetProperty(fs, property, value)
}

func (this *Zfs) SetProperty(fs, property, value string) error {
	c := this.Command("zfs set " + property + "=" + value + " " + fs)
	err := c.CmdError()
	if err != nil {
		return err
	}
	if _, err = c.Run(); err != nil {
		return err
	}
	out, err := this.Property(fs, property)
	if err != nil {
		return err
	}
	if out != value {
		return errors.New("cannot set property: " + property)
	}
	return nil
}

func RecentSnap(snap, property string) (string, error) {
	return std.RecentSnap(snap, property)
}

func (this *Zfs) RecentSnap(snap, property string) (string, error) {
	c := this.Command("zfs list -Hro name -t snapshot -S creation " + snap)
	err := c.CmdError()
	if err != nil {
		return "", err
	}
	out, err := c.Run()
	if err != nil {
		return "", err
	}
	for _, snap := range out {
		if property != "" {
			out, err := this.Property(snap, property)
			if err != nil {
				return "", nil
			}
			if out == "true" {
				return snap, nil
			}
			continue
		}
		return snap, nil
	}
	return "", nil
}

func RecvSnap(fs, snap string) (runcmd.CmdWorker, error) {
	return std.RecvSnap(fs, snap)
}

func (this *Zfs) RecvSnap(fs, snap string) (runcmd.CmdWorker, error) {
	c := this.Command("zfs recv " + fs + "@" + snap)
	err := c.CmdError()
	if err != nil {
		return nil, err
	}
	err = c.Start()
	return c, nil
}

func SendSnap(fs, snapOld, snapNew string, cw runcmd.CmdWorker) (runcmd.CmdWorker, error) {
	return std.SendSnap(fs, snapOld, snapNew, cw)
}
func (this *Zfs) SendSnap(fs, snapOld, snapNew string, cw runcmd.CmdWorker) (runcmd.CmdWorker, error) {
	cmd := "zfs send -i " + fs + "@" + snapOld + " " + fs + "@" + snapNew
	if snapNew == "" {
		cmd = "zfs send " + fs + "@" + snapOld
	}
	sendWorker := this.Command(cmd)
	err := sendWorker.CmdError()
	if err != nil {
		return nil, err
	}

	sendWorkerStdout, err := sendWorker.StdoutPipe()
	if err != nil {
		return nil, err
	}

	cwStdin, err := cw.StdinPipe()
	if err != nil {
		return nil, err
	}

	if err := sendWorker.Start(); err != nil {
		return nil, err
	}
	_, err = io.Copy(cwStdin, sendWorkerStdout)
	return sendWorker, err
}
