package zfs

import (
	"errors"
	"io"
	"strings"

	"github.com/theairkit/runcmd"
)

const (
	FS       = "filesystem"
	SNAP     = "snapshot"
	PROPERTY = "zbackup:"
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
	c, err := this.Command("zfs snapshot " + fs + "@" + snap)
	if err != nil {
		return err
	}
	_, err = c.Run()
	return err
}

func DestroyFs(fs string) error {
	return std.DestroyFs(fs)
}

func (this *Zfs) DestroyFs(fs string) error {
	c, err := this.Command("zfs destroy -r " + fs)
	if err != nil {
		return err
	}
	_, err = c.Run()
	return err
}

func ExistFs(fs, fsType string) (bool, error) {
	return std.ExistFs(fs, fsType)
}

func (z *Zfs) ExistFs(fs, fsType string) (bool, error) {
	if _, err := z.ListFs(fs, fsType, "", false); err != nil {
		if strings.Contains(err.Error(), "dataset does not exist") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func ExistSnap(remote, snapshot string) (bool, error) {
	return std.ExistSnap(remote+"@"+snapshot, SNAP)
}

func (z *Zfs) ExistSnap(remote, snapshot string) (bool, error) {
	return z.ExistFs(remote+"@"+snapshot, SNAP)
}

func ListFs(fs, fsType, sortProp string, recursive bool) ([]string, error) {
	return std.ListFs(fs, fsType, sortProp, recursive)
}

func (this *Zfs) ListFs(fs, fsType, sortProp string, recursive bool) ([]string, error) {
	r := ""
	if recursive {
		r = "-r"
	}
	cmd := "zfs list -Ho name -t " + fsType + " " + r + " " + fs
	if sortProp != "" {
		cmd += " -S " + sortProp
	}
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

func RecentSnap(fs, property string) (string, error) {
	return std.RecentSnap(fs, property)
}

func (this *Zfs) RecentSnap(fs, property string) (string, error) {
	c, err := this.Command("zfs list -Hro name -t snapshot -S creation " + fs)
	if err != nil {
		return "", err
	}
	out, err := c.Run()
	if err != nil {
		return "", err
	}
	for _, snap := range out {
		val, err := this.Property(snap, property)
		if err != nil {
			return "", nil
		}
		if val == "true" {
			return snap, nil
		}
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

func RenameFs(oldName, newName string) error {
	return std.RenameFs(oldName, newName)
}

func (this *Zfs) RenameFs(oldName, newName string) error {
	c, err := this.Command("zfs rename " + oldName + " " + newName)
	if err != nil {
		return err
	}
	_, err = c.Run()
	return err
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
	val, err := this.Property(fs, property)
	if err != nil {
		return err
	}
	if val != value {
		return errors.New("cannot set property: " + property)
	}
	return nil
}

func SendSnap(fs, snapCurr, snapNew string, cw runcmd.CmdWorker) error {
	return std.SendSnap(fs, snapCurr, snapNew, cw)
}
func (this *Zfs) SendSnap(fs, snapCurr, snapNew string, cw runcmd.CmdWorker) error {
	cmd := "zfs send -i " + fs + "@" + snapCurr + " " + fs + "@" + snapNew
	if snapNew == "" {
		cmd = "zfs send " + fs + "@" + snapCurr
	}
	sendCmd, err := this.Command(cmd)
	if err != nil {
		return err
	}
	if err := sendCmd.Start(); err != nil {
		return err
	}
	_, err = io.Copy(cw.StdinPipe(), sendCmd.StdoutPipe())
	return err
}
