//
// A seduko type with methods to read and solve it
//
package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Cell represents a single cell on the seduko grid
type puzzleCell struct {
	value int  // current value of this cell (0 means not set)
	fixed bool // is this a fixed value cell (i.e. supplied in definition)
}

// PuzzleStatus represents a Seduko status code
type puzzleStatus int

const (
	statusNotSolved  puzzleStatus = iota // No attempt made to solve
	statusSolved     puzzleStatus = iota // Solution has been found
	statusNoSolution puzzleStatus = iota // No solution exists
)

const gridSize = 9

// Represents a seduko grid
type bits uint16
type seduko struct {
	cells         [][]puzzleCell // Grid of cells
	status        puzzleStatus
	rows          []bits        // bitmap of set values per row
	columns       []bits        // bitmap of set values per column
	regions       []bits        // bitmap of set values per region
	solveDuration time.Duration // how long this took to solve (wall clock time)
}

// Error type for an invalid seduko definition
type invalidSeduko struct {
	message     string
	nestedError error
}

func (err *invalidSeduko) Error() string {
	if len(err.message) == 0 && err.nestedError != nil {
		return err.nestedError.Error()
	}
	return err.message
}

// Check if this is a valid value for this cell
// If valid, set it and return true
// If not valid return false
func (sed *seduko) trySetCell(row, col, val int, fixed bool) bool {
	value := uint(val)
	bit := bits(1 << value)
	if (sed.rows[row] & bit) != 0 {
		return false
	}
	if (sed.columns[col] & bit) != 0 {
		return false
	}
	region := ((row / 3) * 3) + (col / 3)
	if (sed.regions[region] & bit) != 0 {
		return false
	}
	sed.rows[row] |= bit
	sed.columns[col] |= bit
	sed.regions[region] |= bit
	sed.cells[row][col].value = val
	sed.cells[row][col].fixed = fixed
	return true
}

// Clear a set cell value
func (sed *seduko) undoSetCell(row, col, value int) {
	bit := bits(1 << uint(value))
	region := ((row / 3) * 3) + (col / 3)

	sed.rows[row] &= ^bit
	sed.columns[col] &= ^bit
	sed.regions[region] &= ^bit
	sed.cells[row][col].value = 0
}

//
// Initialise the Seduko and read its definition from the supplied scanner
// Returns an error if the defintion is not valid
//
func (sed *seduko) init(scanner *bufio.Scanner) error {
	sed.status = statusNotSolved

	// set bitmaps used to track which values have been used
	sed.rows = make([]bits, gridSize, gridSize)
	sed.columns = make([]bits, gridSize, gridSize)
	sed.regions = make([]bits, gridSize, gridSize)

	// create the grid
	sed.cells = make([][]puzzleCell, gridSize, gridSize)
	for i := 0; i < gridSize; i++ {
		sed.cells[i] = make([]puzzleCell, gridSize, gridSize)
	}

	if scanner != nil {
		row := 0
		for scanner.Scan() {
			line := scanner.Text()
			nums := strings.Split(line, ",")
			if len(nums) != gridSize {
				return &invalidSeduko{message: fmt.Sprintf("Invalid seduko definition: Must have exactly %d columns, have %d on line %d", gridSize, len(nums), row+1)}
			}
			for col := 0; col < gridSize; col++ {
				if len(nums[col]) != 0 {
					val, err := strconv.ParseInt(nums[col], 10, 64)
					if err != nil {
						return &invalidSeduko{nestedError: err}
					}
					if !sed.trySetCell(row, col, int(val), true) {
						return &invalidSeduko{message: fmt.Sprintf("Number %d at cell (%d,%d) is not valid", val, row+1, col+1)}
					}
				}
			}
			if row == gridSize-1 {
				break // don't read any more lines
			}
			row++
		}

		if row != gridSize-1 {
			return &invalidSeduko{message: fmt.Sprintf("Only %d rows read, must have %d!", row, gridSize)}
		}
	}

	return nil

	/*
		// initialise with a test puzzle
		ok := sed.trySetCell(0, 3, 8, true)
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
		return nil
	*/
}

func (sed *seduko) print(showTiming bool) {
	switch sed.status {
	case statusSolved:
		fmt.Printf("Solved:\n")
	case statusNoSolution:
		fmt.Printf("No Solution Found:\n")
	case statusNotSolved:
		fmt.Printf("Unsolved:\n")
	}
	for row := 0; row < len(sed.cells); row++ {
		if row == 0 {
			fmt.Printf(" -------------------------------------\n")
		} else if row%3 == 0 {
			fmt.Printf(" |---|---|---|---|---|---|---|---|---|\n")
		}
		for col := 0; col < len(sed.cells[0]); col++ {
			if val := sed.cells[row][col].value; val == 0 {
				fmt.Print(" |  ")
			} else {
				fmt.Printf(" | %d", val)
			}
		}
		fmt.Printf(" |\n")
	}
	fmt.Printf(" -------------------------------------\n\n")
	if showTiming && sed.status == statusSolved {
		fmt.Printf("Solved in %g ms\n", float64(sed.solveDuration)/float64(time.Millisecond))
	}
}

// Solve searches for a solution to the Seduko
// If one is found store it and set the status to statusSolved
func (sed *seduko) solve() error {

	start := time.Now()
	found, err := sed.findSolution(0, 0)
	sed.solveDuration = time.Since(start)

	if err != nil {
		return err
	}
	if found {
		sed.status = statusSolved
	} else {
		sed.status = statusNoSolution
	}
	return nil
}

// Return the next cell in the grid (or (10,1) if we move past the end)
func nextCell(row, col int) (int, int) {
	if row >= gridSize || col >= gridSize || row < 0 || col < 0 {
		panic(fmt.Sprintf("nextCell: Invalid (row,col) : (%d,%d)", row, col))
	}
	col++
	if col == gridSize {
		col = 0
		row++
	}
	return row, col
}

//
// Find a as solution for the puzzle starting at the given cell using
// a simple depth first search with back-tracking
// Assume all earlier cells have a (possibly provisional )value set
//
func (sed *seduko) findSolution(row, col int) (found bool, err error) {

	// If we reach past the end we have found a valid solution for all cells
	if row >= gridSize {
		return true, nil // solution found!
	}

	// If this is a fixed cell just move onto the next one
	var cell = &(sed.cells[row][col])
	if cell.fixed {
		nextR, nextC := nextCell(row, col)
		return sed.findSolution(nextR, nextC)
	}

	// otherwise try each possible value looking for a valid solution
	for val := 1; val <= gridSize; val++ {
		if sed.trySetCell(row, col, val, false) {
			nextR, nextC := nextCell(row, col)
			found, err := sed.findSolution(nextR, nextC)
			if found || err != nil {
				return found, err // have a solution or an error occured
			}
			// rollback this value
			sed.undoSetCell(row, col, val)
		}
	}

	// if we get here no valid value exists for cell (row,col) - need to back-track
	return false, nil
}
