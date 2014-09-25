package zfs

import (
	"errors"
	"io"
	"strings"

	"github.com/theairkit/runcmd"
)

const (
	FS         = "filesystem"
	SNAP       = "snapshot"
	errStrange = "something goes wrong: no errors, but "
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

func CreateSnapshot(fs, snapName string) error {
	return std.CreateSnapshot(fs, snapName)
}

func (this *Zfs) CreateSnapshot(fs, snapName string) error {
	c, err := this.Command("zfs snapshot " + fs + "@" + snapName)
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
	if _, err := z.ListFs(fs, fsType, false); err != nil {
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

func ListFs(fsName, fsType string, recursive bool) ([]string, error) {
	return std.ListFs(fsName, fsType, recursive)
}

func (this *Zfs) ListFs(fsName, fsType string, recursive bool) ([]string, error) {
	fsList := make([]string, 0)
	cmd := "zfs list -H -o name -t " + fsType
	if strings.HasSuffix(fsName, "*") {
		c, err := this.Command(cmd)
		if err != nil {
			return nil, err
		}
		out, err := c.Run()
		if err != nil {
			return nil, err
		}
		for _, fs := range out {
			if strings.HasPrefix(fs, strings.Trim(fsName, "*")) {
				fsList = append(fsList, fs)
			}
		}
		return fsList, nil
	}
	switch {
	case recursive:
		cmd = cmd + " -r"
	case fsType == SNAP:
		cmd = cmd + " -d1"
	}
	c, err := this.Command(cmd + " " + fsName)
	if err != nil {
		return nil, err
	}
	out, err := c.Run()
	if err != nil {
		return nil, err
	}
	for _, fs := range out {
		fsList = append(fsList, fs)
	}
	return fsList, nil
}

func ListFsByProperty(property string) ([]string, error) {
	return std.ListFsByProperty(property)
}

func (this *Zfs) ListFsByProperty(property string) ([]string, error) {
	fsList := make([]string, 0)
	c, err := this.Command("zfs get -H -o name -s local " + property)
	if err != nil {
		return nil, err
	}
	out, err := c.Run()
	if err != nil {
		return nil, err
	}
	for _, fs := range out {
		fsList = append(fsList, fs)
	}
	return fsList, nil
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

func SendSnapshot(fs, snapCurr, snapNew string, cw runcmd.CmdWorker) error {
	return std.SendSnapshot(fs, snapCurr, snapNew, cw)
}
func (this *Zfs) SendSnapshot(fs, snapCurr, snapNew string, cw runcmd.CmdWorker) error {
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

func RecvSnapshot(fs, snap string) (runcmd.CmdWorker, error) {
	return std.RecvSnapshot(fs, snap)
}

func (this *Zfs) RecvSnapshot(fs, snap string) (runcmd.CmdWorker, error) {
	c, err := this.Command("zfs recv " + fs + "@" + snap)
	if err != nil {
		return nil, err
	}
	err = c.Start()
	return c, nil
}

func SetProperty(property, value, fs string) error {
	return std.SetProperty(property, value, fs)
}

func (this *Zfs) SetProperty(property, value, fs string) error {
	c, err := this.Command("zfs set " + property + "=" + value + " " + fs)
	if err != nil {
		return err
	}
	if _, err = c.Run(); err != nil {
		return err
	}
	out, err := this.Property(property, fs)
	if err != nil {
		return err
	}
	if out != value {
		return errors.New(errStrange + "cannot set property: " + property)
	}
	return nil
}

func Property(property, fs string) (string, error) {
	return std.Property(property, fs)
}

func (this *Zfs) Property(property, fs string) (string, error) {
	c, err := this.Command("zfs get -H -o value " + property + " " + fs)
	if err != nil {
		return "", err
	}
	out, err := c.Run()
	if err != nil {
		return "", err
	}
	if len(out) > 1 {
		return "", errors.New(errStrange + "property is multivalue: " + strings.Join(out, "\n"))
	}
	return out[0], nil
}

func RecentSnapshot(fs string) (string, error) {
	return std.RecentSnapshot(fs)
}

func (this *Zfs) RecentSnapshot(fs string) (string, error) {
	recent := ""
	c, err := this.Command("zfs list -Hrt snapshot -o name -S creation " + fs)
	if err != nil {
		return "", err
	}
	snapList, err := c.Run()
	if err != nil {
		return "", err
	}
	for i := 0; i < len(snapList); i++ {
		prop, err := this.Property("zbackup:", snapList[i])
		if err != nil {
			return "", nil
		}
		if prop == "true" {
			recent = snapList[i]
			break
		}
	}
	return recent, nil
}
