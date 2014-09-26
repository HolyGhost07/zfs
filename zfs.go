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
	PROPERTY   = "zbackup:"
	fsNotexist = "dataset does not exist"
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

func CreateSnap(fs, snapName string) error {
	return std.CreateSnap(fs, snapName)
}

func (this *Zfs) CreateSnap(fs, snapName string) error {
	command, err := this.Command(
		"zfs snapshot " + fs + "@" + snapName,
	)
	if err != nil {
		return err
	}
	_, err = command.Run()
	return err
}

func DestroyFs(fs string) error {
	return std.DestroyFs(fs)
}

func (this *Zfs) DestroyFs(fs string) error {
	command, err := this.Command(
		"zfs destroy -r " + fs,
	)
	if err != nil {
		return err
	}
	_, err = command.Run()
	return err
}

func ExistFs(fs, fsType string) (bool, error) {
	return std.ExistFs(fs, fsType)
}

func (this *Zfs) ExistFs(fs, fsType string) (bool, error) {
	if _, err := this.ListFs(fs, fsType, false); err != nil {
		if strings.Contains(err.Error(), fsNotexist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func ExistSnap(remote, snapshot string) (bool, error) {
	return std.ExistSnap(remote+"@"+snapshot, SNAP)
}

func (this *Zfs) ExistSnap(remote, snapshot string) (bool, error) {
	return this.ExistFs(remote+"@"+snapshot, SNAP)
}

func ListFs(fsName, fsType string, recursive bool) ([]string, error) {
	return std.ListFs(fsName, fsType, recursive)
}

func (this *Zfs) ListFs(fsName, fsType string, recursive bool) ([]string, error) {
	fsList := make([]string, 0)
	c := "zfs list -H -o name -t " + fsType
	if strings.HasSuffix(fsName, "*") {
		command, err := this.Command(c)
		if err != nil {
			return nil, err
		}
		out, err := command.Run()
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
		c = c + " -r"
	case fsType == SNAP:
		c = c + " -d1"
	}
	command, err := this.Command(c + " " + fsName)
	if err != nil {
		return nil, err
	}
	out, err := command.Run()
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
	command, err := this.Command(
		"zfs get -H -o name -s local " + property,
	)
	if err != nil {
		return nil, err
	}
	out, err := command.Run()
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
	command, err := this.Command(
		"zfs rename " + oldName + " " + newName,
	)
	if err != nil {
		return err
	}
	_, err = command.Run()
	return err
}

func SendSnap(fs, snapCurr, snapNew string, cw runcmd.CmdWorker) error {
	return std.SendSnap(fs, snapCurr, snapNew, cw)
}
func (this *Zfs) SendSnap(fs, snapCurr, snapNew string, cw runcmd.CmdWorker) error {
	cmd := "zfs send -i " + fs + "@" + snapCurr + " " + fs + "@" + snapNew
	if snapNew == "" {
		cmd = "zfs send " + fs + "@" + snapCurr
	}
	command, err := this.Command(cmd)
	if err != nil {
		return err
	}
	if err := command.Start(); err != nil {
		return err
	}
	_, err = io.Copy(cw.StdinPipe(), command.StdoutPipe())
	return err
}

func RecvSnap(fs, snap string) (runcmd.CmdWorker, error) {
	return std.RecvSnap(fs, snap)
}

func (this *Zfs) RecvSnap(fs, snap string) (runcmd.CmdWorker, error) {
	command, err := this.Command(
		"zfs recv " + fs + "@" + snap,
	)
	if err != nil {
		return nil, err
	}
	err = command.Start()
	return command, nil
}

func SetProperty(property, value, fs string) error {
	return std.SetProperty(property, value, fs)
}

func (this *Zfs) SetProperty(property, value, fs string) error {
	command, err := this.Command(
		"zfs set " + property + "=" + value + " " + fs,
	)
	if err != nil {
		return err
	}
	if _, err = command.Run(); err != nil {
		return err
	}
	out, err := this.Property(property, fs)
	if err != nil {
		return err
	}
	if out != value {
		return errors.New("cannot set property: " + property)
	}
	return nil
}

func Property(property, fs string) (string, error) {
	return std.Property(property, fs)
}

func (this *Zfs) Property(property, fs string) (string, error) {
	command, err := this.Command(
		"zfs get -H -o value " + property + " " + fs,
	)
	if err != nil {
		return "", err
	}
	out, err := command.Run()
	if err != nil {
		return "", err
	}
	if len(out) > 1 {
		return "", errors.New("property is multivalue: " + strings.Join(out, "\n"))
	}
	return out[0], nil
}

func RecentSnapshot(fs string) (string, error) {
	return std.RecentSnapshot(fs)
}

func (this *Zfs) RecentSnapshot(fs string) (string, error) {
	recent := ""
	command, err := this.Command(
		"zfs list -Hrt snapshot -o name -S creation " + fs,
	)
	if err != nil {
		return "", err
	}
	snapList, err := command.Run()
	if err != nil {
		return "", err
	}
	for i := 0; i < len(snapList); i++ {
		prop, err := this.Property(PROPERTY, snapList[i])
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
