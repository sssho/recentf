package lib

import (
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

type WscriptShell struct {
	Shell  *ole.IUnknown
	Wshell *ole.IDispatch
}

func NewWscriptShell() (*WscriptShell, error) {
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED|ole.COINIT_SPEED_OVER_MEMORY)
	shell, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return nil, err
	}
	wshell, err := shell.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		shell.Release()
		return nil, err
	}
	return &WscriptShell{shell, wshell}, nil
}

func (w *WscriptShell) Close() {
	w.Shell.Release()
	w.Wshell.Release()
	ole.CoUninitialize()
}

func (w *WscriptShell) ShortcutInfo(path string) (string, string, error) {
	shortcut, err := oleutil.CallMethod(w.Wshell, "CreateShortcut", path)
	if err != nil {
		return "", "", err
	}
	shortcutDispath := shortcut.ToIDispatch()

	targetPath, err := shortcutDispath.GetProperty("TargetPath")
	if err != nil {
		return "", "", err
	}

	args, err := shortcutDispath.GetProperty("Arguments")
	if err != nil {
		return "", "", err
	}
	return targetPath.ToString(), args.ToString(), nil
}

func ResolveShortcut(path string) (string, string, error) {
	w, err := NewWscriptShell()
	if err != nil {
		return "", "", err
	}
	defer w.Close()
	targetPath, args, err := w.ShortcutInfo(path)
	if err != nil {
		return "", "", err
	}
	return targetPath, args, nil
}
