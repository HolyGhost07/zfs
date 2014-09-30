package zfs

import (
	"errors"
	"io"
	"strings"

	"github.com/theairkit/runcmd"
)

var (
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

func CreateFs(fs string) error {
	return std.CreateFs(fs)
}

func (this *Zfs) CreateFs(fs string) error {
	c, err := this.Command("zfs create " + fs)
	if err != nil {
		return err
	}
	_, err = c.Run()
	return err
}

func CreateSnap(fs, snap string) error {
	return std.CreateSnap(fs, snap)
}

func (this *Zfs) CreateSnap(fs, snap string) error {
	return this.CreateFs(fs + "@" + snap)
}

func DestroyFs(fs string) error {
	return std.DestroyFs(fs)
}

func (this *Zfs) DestroyFs(fs string) error {
	c, err := this.Command("zfs destroy " + fs)
	if err != nil {
		return err
	}
	_, err = c.Run()
	return err
}

func DestroySnap(fs, snap string) error {
	return std.DestroySnap(fs, snap)
}

func (this *Zfs) DestroySnap(fs, snap string) error {
	return this.DestroyFs(fs + "@" + snap)
}

func RenameFs(fsOld, fsNew string) error {
	return std.RenameFs(fsOld, fsNew)
}

func (this *Zfs) RenameFs(fsOld, fsNew string) error {
	c, err := this.Command("zfs rename " + fsOld + " " + fsNew)
	if err != nil {
		return err
	}
	_, err = c.Run()
	return err
}

func RenameSnap(snapOld, snapNew string) error {
	return std.RenameSnap(snapOld, snapNew)
}

func (this *Zfs) RenameSnap(snapOld, snapNew string) error {
	return this.RenameFs(snapOld, snapNew)
}

func ExistFs(fs string) (bool, error) {
	return std.ExistFs(fs)
}

func (this *Zfs) ExistFs(fs string) (bool, error) {
	if _, err := this.ListFs(fs, false); err != nil {
		if strings.Contains(err.Error(), DATANOE) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func ExistSnap(snap string) (bool, error) {
	return std.ExistSnap(snap)
}

func (this *Zfs) ExistSnap(snap string) (bool, error) {
	if _, err := this.ListSnap(snap, false); err != nil {
		if strings.Contains(err.Error(), DATANOE) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func ListFs(fs string, recursive bool) ([]string, error) {
	return std.ListFs(fs, recursive)
}

func (this *Zfs) ListFs(fs string, recursive bool) ([]string, error) {
	r := ""
	if recursive {
		r = "-r"
	}
	cmd := "zfs list -Ho name -t filesystem " + r + " " + fs
	c, err := this.Command(cmd)
	if err != nil {
		return nil, err
	}
	return c.Run()
}

func ListSnap(snap string, recursive bool) ([]string, error) {
	return std.ListSnap(snap, recursive)
}

func (this *Zfs) ListSnap(snap string, recursive bool) ([]string, error) {
	r := ""
	if recursive {
		r = "-r"
	}
	cmd := "zfs list -Ho name -t snapshot " + r + " " + snap
	c, err := this.Command(cmd)
	if err != nil {
		return nil, err
	}
	return c.Run()
}

func Property(fs, property string) (string, error) {
	return std.Property(fs, property)
}

func (this *Zfs) Property(fs, property string) (string, error) {
	c, err := this.Command("zfs get -H -o value " + property + " " + fs)
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
	c, err := this.Command("zfs set " + property + "=" + value + " " + fs)
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
	c, err := this.Command("zfs list -Hro name -t snapshot -S creation " + snap)
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
	c, err := this.Command("zfs recv -F " + fs + "@" + snap)
	if err != nil {
		return nil, err
	}
	err = c.Start()
	return c, nil
}

func SendSnap(fs, snapCurr, snapNew string, cw runcmd.CmdWorker) error {
	return std.SendSnap(fs, snapCurr, snapNew, cw)
}
func (this *Zfs) SendSnap(fs, snapCurr, snapNew string, cw runcmd.CmdWorker) error {
	cmd := "zfs send -i " + fs + "@" + snapCurr + " " + fs + "@" + snapNew
	if snapNew == "" {
		cmd = "zfs send " + fs + "@" + snapCurr
	}
	c, err := this.Command(cmd)
	if err != nil {
		return err
	}
	if err := c.Start(); err != nil {
		return err
	}
	_, err = io.Copy(cw.StdinPipe(), c.StdoutPipe())
	return err
}
