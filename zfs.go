package zfs

import (
	"fmt"
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
	fmt.Println("zfs snapshot " + fs + "@" + snapName)
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
	fmt.Println("zfs destroy -r " + fs)
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

	fmt.Println("zfs list -H -o name -t " + fsType)

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

	fmt.Println(cmd + " " + fsName)
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
	fmt.Println("zfs get -H -o name -s local " + property)
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
	fmt.Println("zfs rename " + oldName + " " + newName)
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
	fmt.Println(cmd)
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
	c, err := this.Command("zfs recv -F " + fs + "@" + snap)
	if err != nil {
		return nil, err
	}
	fmt.Println(c)
	err = c.Start()
	return c, nil
}
