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
	r := ""
	if recursive {
		r = "-r "
	}
	if strings.Contains(fsName, "*") {
		allFs, err := this.Command("zfs list -H -t " + fsType + " -o name -r").Run()
		if err != nil {
			return fsList, err
		}
		for _, fs := range allFs {
			if strings.Contains(fs, strings.Trim(fsName, "*")) {
				fsList = append(fsList, fs)
			}
		}
		return fsList, err
	}
	return this.Command("zfs list -H -t " + fsType + " -o name " + r + fsName).Run()
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
	cmd := this.Command(sendCmd(fs, snapCurr, snapNew))
	if err := cmd.Start(); err != nil {
		return err
	}
	_, err := io.Copy(cw.Stdin(), cmd.Stdout())
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

func sendCmd(fs, snapCurr, snapNew string) string {
	if snapNew == "" {
		return "zfs send " + fs + "@" + snapCurr
	}
	return "zfs send -i " + fs + "@" + snapCurr + " " + fs + "@" + snapNew
}
