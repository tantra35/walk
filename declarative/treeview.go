// Copyright 2012 The Walk Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package declarative

import (
	"github.com/lxn/walk"
)

type TreeView struct {
	AssignTo             **walk.TreeView
	Name                 string
	Disabled             bool
	Hidden               bool
	Font                 Font
	MinSize              Size
	MaxSize              Size
	StretchFactor        int
	Row                  int
	RowSpan              int
	Column               int
	ColumnSpan           int
	ContextMenuActions   []*walk.Action
	Model                walk.TreeModel
	OnCurrentItemChanged walk.EventHandler
	OnItemCollapsed      walk.TreeItemEventHandler
	OnItemExpanded       walk.TreeItemEventHandler
}

func (tv TreeView) Create(parent walk.Container) error {
	w, err := walk.NewTreeView(parent)
	if err != nil {
		return err
	}

	return InitWidget(tv, w, func() error {
		if err := w.SetModel(tv.Model); err != nil {
			return err
		}

		if tv.OnCurrentItemChanged != nil {
			w.CurrentItemChanged().Attach(tv.OnCurrentItemChanged)
		}

		if tv.OnItemCollapsed != nil {
			w.ItemCollapsed().Attach(tv.OnItemCollapsed)
		}

		if tv.OnItemExpanded != nil {
			w.ItemExpanded().Attach(tv.OnItemExpanded)
		}

		if tv.AssignTo != nil {
			*tv.AssignTo = w
		}

		return nil
	})
}

func (tv TreeView) WidgetInfo() (name string, disabled, hidden bool, font *Font, minSize, maxSize Size, stretchFactor, row, rowSpan, column, columnSpan int, contextMenuActions []*walk.Action) {
	return tv.Name, tv.Disabled, tv.Hidden, &tv.Font, tv.MinSize, tv.MaxSize, tv.StretchFactor, tv.Row, tv.RowSpan, tv.Column, tv.ColumnSpan, tv.ContextMenuActions
}
