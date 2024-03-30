// Copyright 2011 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/md5"
	"crypto/rsa"
	"fmt"
	"log"

	"github.com/kayrus/putty"
	"github.com/lxn/walk"

	. "github.com/lxn/walk/declarative"
)

type PuttyKey struct {
	key            *putty.Key
	bitlen         int
	pubfingerprint string
	checked        bool
}

func NewPuttyKey(keyfile string) (*PuttyKey, error) {
	puttyKey, err := putty.NewFromFile(keyfile)
	if err != nil {
		return nil, err
	}

	if puttyKey.Encryption != "none" {
		return nil, fmt.Errorf("Keys with passwords are unsuported")
	}

	_, err = puttyKey.ParseRawPrivateKey(nil)
	if err != nil {
		return nil, err
	}

	pubkey, _ := puttyKey.ParseRawPublicKey()
	var lpubkeylen int = -1
	switch key := pubkey.(type) {
	case *rsa.PublicKey:
		lpubkeylen = key.N.BitLen()

	case *ecdsa.PublicKey:
		lpubkeylen = key.Params().BitSize

	case *dsa.PublicKey:
		lpubkeylen = key.P.BitLen()
	}

	fp := md5.Sum(puttyKey.PublicKey)
	lsfp := ""

	for i, b := range fp {
		lsfp += fmt.Sprintf("%02x", b)
		if i < len(fp)-1 {
			lsfp += ":"
		}
	}

	return &PuttyKey{puttyKey, lpubkeylen, lsfp, false}, nil
}

type PuttyKeysModel struct {
	walk.TableModelBase
	items []*PuttyKey
}

func NewPuttyKeysModel(_keyfiles []string) *PuttyKeysModel {
	keysmodel := &PuttyKeysModel{items: []*PuttyKey{}}

	for _, lkey := range _keyfiles {
		litem, err := NewPuttyKey(lkey)
		if err != nil {
			continue
		}

		keysmodel.items = append(keysmodel.items, litem)
	}

	return keysmodel
}

func (m *PuttyKeysModel) RowCount() int {
	return len(m.items)
}

func (m *PuttyKeysModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.key.Algo

	case 1:
		return item.bitlen

	case 2:
		return item.pubfingerprint

	case 3:
		return item.key.Comment
	}

	panic("unexpected col")
}

func (m *PuttyKeysModel) Checked(row int) bool {
	return m.items[row].checked
}

// Called by the TableView when the user toggled the check box of a given row.
func (m *PuttyKeysModel) SetChecked(row int, checked bool) error {
	m.items[row].checked = checked

	return nil
}

func (m *PuttyKeysModel) AddItem(keyfile string) {
	item, err := NewPuttyKey(keyfile)
	if err != nil {
		return
	}

	m.items = append(m.items, item)
	m.PublishRowsInserted(len(m.items), len(m.items))
}

func (m *PuttyKeysModel) RemoveItem(itemIndex int) {
	if itemIndex < 0 {
		return
	}

	m.items = append(m.items[:itemIndex], m.items[itemIndex+1:]...)
	m.PublishRowsRemoved(itemIndex, itemIndex)
}

func main() {
	keysmodel := NewPuttyKeysModel([]string{"./ruslan.ppk", "./ruslan-20032019.ppk"})

	ico, _ := walk.NewIconFromResourceId(2)
	mw, _ := walk.NewMainWindow()

	// Create the notify icon and make sure we clean it up on exit.
	ni, err := walk.NewNotifyIcon(mw)
	if err != nil {
		log.Fatal(err)
	}
	defer ni.Dispose()

	// Set the icon and a tool tip text.
	if err := ni.SetIcon(ico); err != nil {
		log.Fatal(err)
	}
	if err := ni.SetToolTip("Click for info or use the context menu to exit."); err != nil {
		log.Fatal(err)
	}

	var lldldg *walk.Dialog
	var tv *walk.TableView

	viewKeysAction := walk.NewAction()
	ni.ContextMenu().Actions().Add(viewKeysAction)
	viewKeysAction.SetText("View Keys")
	viewKeysAction.Triggered().Attach(func() {
		if lldldg != nil {
			lldldg.SetFocus()
		}

		Dialog{
			Title:      "PAgent Key List",
			MinSize:    Size{600, 400},
			AssignTo:   &lldldg,
			Persistent: true,
			Icon:       ico,
			Layout:     VBox{},
			Children: []Widget{
				TableView{
					AssignTo:         &tv,
					HeaderHidden:     true,
					AlternatingRowBG: false,
					CheckBoxes:       false,
					ColumnsOrderable: false,
					MultiSelection:   false,
					Columns: []TableViewColumn{
						{Title: "Key Type"},
						{Title: "PubKey Length", Width: 50},
						{Title: "PubKey FingerPrint"},
						{Title: "Key Description", Width: 150},
					},
					Model: keysmodel,
					OnBoundsChanged: func() {
						b := tv.Bounds()
						c := tv.Columns()

						lwidth := b.Width - c.At(0).Width() - c.At(1).Width() - c.At(3).Width()
						c.At(2).SetWidth(lwidth)
					},
					OnSelectedIndexesChanged: func() {
						fmt.Printf("SelectedIndexes: %v\n", tv.SelectedIndexes())
					},
				},
				Composite{
					Layout: Grid{Columns: 2, Alignment: AlignHCenterVCenter},
					Children: []Widget{
						PushButton{
							Text: "Add Key",
							OnClicked: func() {
								dlg := new(walk.FileDialog)

								dlg.FilePath = "d:\\src\\walk\\examples\\tableview"
								dlg.Filter = "Putty ppk Files (*.ppk"
								dlg.Title = "Select an key file"

								if ok, err := dlg.ShowOpen(lldldg); err != nil {
									return
								} else if !ok {
									return
								}

								keysmodel.AddItem(dlg.FilePath)
							},
						},
						PushButton{
							Text: "Remove Key",
							OnClicked: func() {
								keysmodel.RemoveItem(tv.CurrentIndex())
							},
						},
					},
				},
			},
		}.Create(nil)

		lldldg.Show()
		lldldg.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
			lldldg = nil
		})
	})

	// We put an exit action into the context menu.
	exitAction := walk.NewAction()
	ni.ContextMenu().Actions().Add(exitAction)
	exitAction.SetText("E&xit")
	exitAction.Triggered().Attach(func() {
		walk.App().Exit(0)
	})

	// The notify icon is hidden initially, so we have to make it visible.
	if err := ni.SetVisible(true); err != nil {
		log.Fatal(err)
	}

	// Run the message loop.
	mw.Run()
}
