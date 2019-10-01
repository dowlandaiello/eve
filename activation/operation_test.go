// Package activation implements a generic activation net.
package activation

import "testing"

/* BEGIN EXPORTED METHODS */

// TestNewOperation tests the functionality of the NewOperator method.
func TestNewOperation(t *testing.T) {
	operator := Add          // Get a copy of the add operator
	var applicant int64 = 10 // Get a copy of the applicant we'll use to initialize the operation
	var parameter int64 = 20 // Get a copy of the parameter we'll use to initialize the operation

	op := NewOperation(operator, applicant, parameter) // Initialize a new operation

	// Check operators not equivalent
	if op.Operator != operator {
		t.Fatalf("operation instance has operator %d, expected %d", op.Operator, operator) // Panic
	} else if op.Applicant != 10 {
		t.Fatalf("operation instance has applicant %d, expected %d", op.Applicant, applicant) // Panic
	} else if op.Parameter != 20 {
		t.Fatalf("operation instance has parameter %d, expected %d", op.Parameter, parameter) // Panic
	}
}

// TestExecuteOperation tests the functionality of the Execute operation method.
func TestExecuteOperation(t *testing.T) {
	operator := Add          // Get a copy of the add operator
	var applicant int64 = 10 // Get a copy of the applicant we'll use to initialize the operation
	var parameter int64 = 20 // Get a copy of the parameter we'll use to initialize the operation
	desiredResult := 30      // Get a copy of the result we'll use later to scrutinize somebody

	op := NewOperation(operator, applicant, parameter) // Initialize a new operation

	// Check does not perform basic arithmetic operations
	if result := op.Execute(); result != 30 {
		t.Fatalf("operation execution resulted in %d, expected %d", result, desiredResult) // Panic
	}
}

/* END EXPORTED METHODS */
