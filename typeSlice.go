package autoconfig

import (
	"reflect"
	"sort"

	qt "github.com/mappu/miqt/qt6"
)

func handle_slice_or_array(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {

	// If there is a struct tag applied to the slice, it will be not used here
	// at all, but it will be propagated into the renderer for the value type
	// of that slice.

	buttons := make([]*qt.QToolButton, 0, 3)

	itemList := qt.NewQTreeWidget2()
	itemList.SetHeaderHidden(true)
	itemList.SetUniformRowHeights(true)
	itemList.SetRootIsDecorated(false)
	itemList.SetSelectionMode(qt.QAbstractItemView__ContiguousSelection)

	refreshListContent := func() {
		itemList.Clear()
		sliceItemsCt := rv.Len()
		for i := 0; i < sliceItemsCt; i++ {
			sliceElem := rv.Index(i)
			listItem := qt.NewQTreeWidgetItem2([]string{formatValue(&sliceElem)})
			itemList.AddTopLevelItem(listItem)
		}
	}
	refreshListContent()

	// Adding (Slice only)

	if rv.Kind() == reflect.Slice {

		addButton := qt.NewQToolButton2()
		setIcon(addButton.QAbstractButton, "list-add", "+", "Add...")
		addButton.SetAutoRaise(true)
		addButton.OnClicked(func() {

			newElem := reflect.New(rv.Type().Elem() /* T */) // pointer-to-T, not a T
			if defaulter, ok := newElem.Interface().(Resetter); ok {
				defaulter.Reset()
			}

			openDialogFor(&newElem, addButton.QWidget, tag, label, func() {

				// insert into slice
				maybeChangedRv := reflect.Append(*rv, newElem.Elem())
				rv.Set(maybeChangedRv)

				// refresh list
				refreshListContent()
			})
		})
		buttons = append(buttons, addButton)

	}

	// Editing (Slice or Array)

	editButton := qt.NewQToolButton2()

	editIndex := func(idx int) {
		curVal := rv.Index(idx)

		openDialogFor(&curVal, editButton.QWidget, tag, label, func() {
			// we have directly mutated inside the slice already

			// refresh list
			refreshListContent()
		})
	}
	itemList.OnDoubleClicked(func(idx *qt.QModelIndex) {
		if idx == nil {
			return // doubleclick was not on an item
		}
		editIndex(idx.Row())
	})

	setIcon(editButton.QAbstractButton, "document-edit-symbolic", "\u270e" /* pencil emoji */, "Edit...")
	editButton.SetAutoRaise(true)
	editButton.OnClicked(func() {
		curIdx := itemList.CurrentIndex()
		if curIdx == nil {
			return // nothing selected
		}
		editIndex(curIdx.Row())
	})

	buttons = append(buttons, editButton)

	// Deleting (Slice only)

	var delButton *qt.QToolButton = nil
	if rv.Kind() == reflect.Slice {
		delButton = qt.NewQToolButton2()
		setIcon(delButton.QAbstractButton, "edit-delete-symbolic", "\u00d7" /* &times; */, "Remove")
		delButton.SetAutoRaise(true)

		delButton.OnClicked(func() {
			selectedItems := itemList.SelectedItems()

			// extract indexes
			var selectedIndexes []int = make([]int, 0, len(selectedItems))
			for _, itm := range selectedItems {
				selectedIndexes = append(selectedIndexes, itemList.IndexOfTopLevelItem(itm))
			}

			// reverse list, so indexes remain stable as we pop them
			sort.Slice(selectedIndexes, func(i, j int) bool { return i > j })

			// remove each item
			for _, removeIdx := range selectedIndexes {
				updated := rv.Slice(0, removeIdx)
				afterPart := rv.Slice(removeIdx+1, rv.Len())
				updated = reflect.AppendSlice(updated, afterPart)
				rv.Set(updated)
			}

			// re-render list
			refreshListContent()
		})

		buttons = append(buttons, delButton)

	}

	refreshButtonsEnabled := func() {
		selCt := len(itemList.SelectedItems())
		editButton.SetEnabled(selCt == 1)
		if delButton != nil { // slice only
			delButton.SetEnabled(selCt > 0)
		}
	}
	refreshButtonsEnabled()
	itemList.OnSelectionChanged(func(super func(*qt.QItemSelection, *qt.QItemSelection), selected, deselected *qt.QItemSelection) {
		refreshButtonsEnabled()
		super(selected, deselected)
	})

	addRowBoxAndButtons(area, label, itemList.QWidget, buttons...)

	return func() {
	}
}
