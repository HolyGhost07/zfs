package zfs

import (
	"io"
	"strings"

	"github.com/theairkit/runcmd"
)

const (
	FS   = "filesystem"
	SNAP = "snapshot"
)

type Zfs struct {
	runcmd.Runner
}

var std = NewZfs(runcmd.NewLocalRunner())

func NewZfs(r runcmd.Runner, e error) *Zfs {
	if e != nil {
		return nil
	}
	return &Zfs{r}
}

func CreateSnapshot(fs, snapName string) error {
	return std.CreateSnapshot(fs, snapName)
}

func (this *Zfs) CreateSnapshot(fs, snapName string) error {
	_, err := this.Command("zfs snapshot " + fs + "@" + snapName).Run()
	return err
}

func DestroyFs(fs string) error {
	return std.DestroyFs(fs)
}

func (this *Zfs) DestroyFs(fs string) error {
	cmd := this.Command("zfs destroy -r " + fs)
	_, err := cmd.Run()
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
		out, err := this.Command(cmd).Run()
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

	if recursive {
		cmd = cmd + " -r"
	}
	out, err := this.Command(cmd + " " + fsName).Run()
	if err != nil {
		return nil, err
	}
	for _, fs := range out {
		fsList = append(fsList, fs)
	}
	return fsList, nil
}

func ListByProperty(property string) ([]string, error) {
	return std.ListByProperty(property)
}

func (this *Zfs) ListByProperty(property string) ([]string, error) {
	fsList := make([]string, 0)
	out, err := this.Command("zfs get -H -o name -s local " + property).Run()
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
	_, err := this.Command("zfs rename " + oldName + " " + newName).Run()
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
	sendCmd := this.Command(cmd)
	if err := sendCmd.Start(); err != nil {
		return err
	}
	_, err := io.Copy(cw.Stdin(), sendCmd.Stdout())
	return err
}

func RecvSnapshot(fs, snap string) (runcmd.CmdWorker, error) {
	return std.RecvSnapshot(fs, snap)
}

func (this *Zfs) RecvSnapshot(fs, snap string) (runcmd.CmdWorker, error) {
	cmd := this.Command("zfs recv -F " + fs + "@" + snap)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}
