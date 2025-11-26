package autoconfig

import (
	"reflect"
	"sort"

	qt "github.com/mappu/miqt/qt6"
)

func handle_slice(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {

	// HorizontalLayout
	// - QTreeWidget
	// - VerticalLayout
	//   - buttons x3
	//   - vspacer

	hbox := qt.NewQHBoxLayout2()

	itemList := qt.NewQTreeWidget2()
	itemList.SetHeaderHidden(true)
	itemList.SetUniformRowHeights(true)
	itemList.SetRootIsDecorated(false)
	itemList.SetSelectionMode(qt.QAbstractItemView__ContiguousSelection)
	hbox.AddWidget(itemList.QWidget)

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

	vbox := qt.NewQVBoxLayout2()
	// TODO if some buttons have icons created and some do not, the widths do
	// not match - should justify/align the widths(!)

	addButton := qt.NewQToolButton2()
	setIcon(addButton.QAbstractButton, "list-add", "+", "Add...")
	addButton.SetAutoRaise(true)
	addButton.OnClicked(func() {

		newElem := reflect.New(rv.Type().Elem() /* T */) // pointer-to-T, not a T
		if defaulter, ok := newElem.Interface().(Resetter); ok {
			defaulter.Reset()
		}

		openDialogFor(&newElem, addButton.QWidget, label, func() {

			// insert into slice
			maybeChangedRv := reflect.Append(*rv, newElem.Elem())
			rv.Set(maybeChangedRv)

			// refresh list
			refreshListContent()
		})
	})
	vbox.AddWidget(addButton.QWidget)

	editIndex := func(idx int) {
		curVal := rv.Index(idx)

		openDialogFor(&curVal, addButton.QWidget, label, func() {
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

	editButton := qt.NewQToolButton2()
	setIcon(editButton.QAbstractButton, "document-edit-symbolic", "\u270e" /* pencil emoji */, "Edit...")
	editButton.SetAutoRaise(true)
	editButton.OnClicked(func() {
		curIdx := itemList.CurrentIndex()
		if curIdx == nil {
			return // nothing selected
		}
		editIndex(curIdx.Row())
	})

	vbox.AddWidget(editButton.QWidget)

	delButton := qt.NewQToolButton2()
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

	vbox.AddWidget(delButton.QWidget)

	refreshButtonsEnabled := func() {
		selCt := len(itemList.SelectedItems())
		editButton.SetEnabled(selCt == 1)
		delButton.SetEnabled(selCt > 0)
	}
	refreshButtonsEnabled()
	itemList.OnSelectionChanged(func(super func(*qt.QItemSelection, *qt.QItemSelection), selected, deselected *qt.QItemSelection) {
		refreshButtonsEnabled()
		super(selected, deselected)
	})

	valign := qt.NewQSpacerItem4(0, 0, qt.QSizePolicy__Minimum, qt.QSizePolicy__MinimumExpanding)
	vbox.AddSpacerItem(valign)

	hbox.AddLayout(vbox.QLayout)

	area.AddRow4(label+`:`, hbox.QLayout)

	return func() {
	}
}
