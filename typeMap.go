package autoconfig

import (
	"reflect"
	"sort"

	qt "github.com/mappu/miqt/qt6"
)

// mapKvPair is a temporary struct that allows editing k/v pairs of a map.
// It implements Renderer{} interface so it can be used with openDialogFor().
type mapKvPair struct {
	Key   reflect.Value
	Value reflect.Value
}

func (mapKvPair) Render(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {

	kField := rv.Field(0).Interface().(reflect.Value)
	vField := rv.Field(1).Interface().(reflect.Value)

	kSaver := handle_any(area, &kField, tag, "Key")
	vSaver := handle_any(area, &vField, tag, "Value")

	return func() {
		kSaver()
		vSaver()
	}
}

func handle_map(area *qt.QFormLayout, rv *reflect.Value, tag reflect.StructTag, label string) SaveFunc {

	// If there is a struct tag applied to the map, it will be not used here
	// at all, but it will be propagated into the renderer for both key+value.

	itemList := qt.NewQTreeWidget2()
	itemList.SetHeaderHidden(false)
	itemList.SetColumnCount(2)
	itemList.SetHeaderLabels([]string{"Key", "Value"})
	itemList.SetUniformRowHeights(true)
	itemList.SetRootIsDecorated(false)
	itemList.SetSelectionMode(qt.QAbstractItemView__ContiguousSelection)

	var currentOrderingKeys []reflect.Value

	refreshListContent := func() {
		itemList.Clear()
		currentOrderingKeys = nil

		// TODO the map order is shuffled on every refresh - should this sort
		// somehow?

		iter := rv.MapRange()
		for iter.Next() {
			kField := iter.Key()
			vField := iter.Value()

			// Keep track of the key values in the current rendering order
			// That means we can use Qt selection indexes to refer to any
			// map key type
			keyCopy := reflect.New(kField.Type()).Elem()
			keyCopy.Set(kField)

			currentOrderingKeys = append(currentOrderingKeys, keyCopy)

			// TODO vField is not addressible here, so formatValue() fails to find
			// (*T) String() if vField has type T
			// Although it works if Stringer is implemented on the value receiver

			listItem := qt.NewQTreeWidgetItem2([]string{formatValue(&kField), formatValue(&vField)})
			itemList.AddTopLevelItem(listItem)
		}
	}
	refreshListContent()

	// Adding

	addButton := qt.NewQToolButton2()
	setIcon(addButton.QAbstractButton, "list-add", "+", "Add...")
	addButton.SetAutoRaise(true)
	addButton.OnClicked(func() {

		newKey := reflect.New(rv.Type().Key() /* T */) // pointer-to-T, not a T
		if defaulter, ok := newKey.Interface().(Resetter); ok {
			defaulter.Reset()
		}

		newValue := reflect.New(rv.Type().Elem() /* T */) // pointer-to-T, not a T
		if defaulter, ok := newValue.Interface().(Resetter); ok {
			defaulter.Reset()
		}

		pair := mapKvPair{newKey.Elem(), newValue.Elem()}
		pairRv := reflect.ValueOf(&pair).Elem()

		openDialogFor(&pairRv, addButton.QWidget, tag, label, func() {

			// Maybe create map if it is nil
			if rv.IsNil() {
				rv.Set(reflect.MakeMap(rv.Type()))
			}

			// insert into map
			rv.SetMapIndex(pair.Key, pair.Value)

			// refresh list
			refreshListContent()
		})
	})

	// Editing (Slice or Array)

	editButton := qt.NewQToolButton2()

	editIndex := func(idx int) {

		curKey := currentOrderingKeys[idx]

		// Copy key data, in case the key was edited, then we would need to
		// delete the old key's entry
		// These are just shallow copies

		keyCopy := reflect.New(curKey.Type()).Elem()
		keyCopy.Set(curKey)

		curVal := rv.MapIndex(curKey) // Not addressible, can't be autoconfig-edited directly

		valCopy := reflect.New(curVal.Type()).Elem()
		valCopy.Set(curVal)

		pair := mapKvPair{keyCopy, valCopy}

		pairRv := reflect.ValueOf(&pair).Elem()

		openDialogFor(&pairRv, editButton.QWidget, tag, label, func() {
			// Move our copied values back into the map

			rv.SetMapIndex(curKey, reflect.Value{})

			rv.SetMapIndex(keyCopy, valCopy)

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

	// Deleting

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

			curKey := currentOrderingKeys[removeIdx]

			// Delete a map index by setting it to an unintiailized reflect.Value{}
			rv.SetMapIndex(curKey, reflect.Value{})
		}

		// re-render list
		refreshListContent()
	})

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

	addRowBoxAndButtons(area, label, itemList.QWidget, addButton, editButton, delButton)

	return func() {
	}
}
