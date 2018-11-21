package main

import (
	"testing"
)

func Example() {

	var sed seduko
	sed.init(nil)
	ok := sed.trySetCell(0, 3, 3, true)
	ok = ok && sed.trySetCell(0, 5, 2, true)
	ok = ok && sed.trySetCell(0, 7, 5, true)
	ok = ok && sed.trySetCell(1, 1, 6, true)
	ok = ok && sed.trySetCell(1, 8, 4, true)
	ok = ok && sed.trySetCell(2, 2, 1, true)
	ok = ok && sed.trySetCell(3, 6, 6, true)
	ok = ok && sed.trySetCell(3, 7, 4, true)
	ok = ok && sed.trySetCell(4, 0, 3, true)
	ok = ok && sed.trySetCell(4, 3, 9, true)
	ok = ok && sed.trySetCell(5, 0, 9, true)
	ok = ok && sed.trySetCell(5, 3, 8, true)
	ok = ok && sed.trySetCell(6, 1, 7, true)
	ok = ok && sed.trySetCell(6, 6, 3, true)
	ok = ok && sed.trySetCell(7, 4, 6, true)
	ok = ok && sed.trySetCell(7, 6, 1, true)
	ok = ok && sed.trySetCell(8, 0, 2, true)

	sed.print(false)

	// Output:
	// Unsolved:
	//  -------------------------------------
	//  |   |   |   | 3 |   | 2 |   | 5 |   |
	//  |   | 6 |   |   |   |   |   |   | 4 |
	//  |   |   | 1 |   |   |   |   |   |   |
	//  |---|---|---|---|---|---|---|---|---|
	//  |   |   |   |   |   |   | 6 | 4 |   |
	//  | 3 |   |   | 9 |   |   |   |   |   |
	//  | 9 |   |   | 8 |   |   |   |   |   |
	//  |---|---|---|---|---|---|---|---|---|
	//  |   | 7 |   |   |   |   | 3 |   |   |
	//  |   |   |   |   | 6 |   | 1 |   |   |
	//  | 2 |   |   |   |   |   |   |   |   |
	//  -------------------------------------
}

func TestCreateValidSeduko(t *testing.T) {

	var sed seduko
	sed.init(nil)
	ok := sed.trySetCell(0, 3, 3, true)
	ok = ok && sed.trySetCell(0, 5, 2, true)
	ok = ok && sed.trySetCell(0, 7, 5, true)
	ok = ok && sed.trySetCell(1, 1, 6, true)
	ok = ok && sed.trySetCell(1, 8, 4, true)
	ok = ok && sed.trySetCell(2, 2, 1, true)
	ok = ok && sed.trySetCell(3, 6, 6, true)
	ok = ok && sed.trySetCell(3, 7, 4, true)
	ok = ok && sed.trySetCell(4, 0, 3, true)
	ok = ok && sed.trySetCell(4, 3, 9, true)
	ok = ok && sed.trySetCell(5, 0, 9, true)
	ok = ok && sed.trySetCell(5, 3, 8, true)
	ok = ok && sed.trySetCell(6, 1, 7, true)
	ok = ok && sed.trySetCell(6, 6, 3, true)
	ok = ok && sed.trySetCell(7, 4, 6, true)
	ok = ok && sed.trySetCell(7, 6, 1, true)
	ok = ok && sed.trySetCell(8, 0, 2, true)

	if !ok {
		t.Fail()
	}
}

func TestCreateInvalidSeduko(t *testing.T) {

	var sed seduko
	sed.init(nil)
	ok := sed.trySetCell(0, 3, 3, true)
	ok = ok && sed.trySetCell(0, 5, 2, true)
	ok = ok && sed.trySetCell(0, 7, 5, true)
	ok = ok && sed.trySetCell(1, 1, 6, true)
	ok = ok && sed.trySetCell(1, 8, 4, true)
	ok = ok && sed.trySetCell(2, 2, 1, true)
	ok = ok && sed.trySetCell(3, 6, 6, true)
	ok = ok && sed.trySetCell(3, 7, 4, true)
	ok = ok && sed.trySetCell(4, 0, 3, true)
	ok = ok && sed.trySetCell(4, 3, 9, true)
	ok = ok && sed.trySetCell(5, 0, 9, true)
	ok = ok && sed.trySetCell(5, 3, 8, true)
	ok = ok && sed.trySetCell(6, 1, 7, true)
	ok = ok && sed.trySetCell(6, 6, 3, true)
	ok = ok && sed.trySetCell(7, 4, 6, true)
	ok = ok && sed.trySetCell(7, 6, 1, true)
	ok = ok && sed.trySetCell(8, 0, 2, true)
	if !ok {
		t.Fail()
	}

	// duplicate in column 3
	if sed.trySetCell(6, 3, 3, true) {
		t.Fail()
	}
	if sed.cells[6][3].value != 0 {
		t.Fail()
	}

	// duplicate in row 8
	if sed.trySetCell(8, 8, 2, true) {
		t.Fail()
	}
	if sed.cells[8][8].value != 0 {
		t.Fail()
	}

	// duplicate last region
	if sed.trySetCell(8, 8, 1, true) {
		t.Fail()
	}
	if sed.cells[8][8].value != 0 {
		t.Fail()
	}
}
