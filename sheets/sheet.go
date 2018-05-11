package sheets

import "gopkg.in/Iwark/spreadsheet.v2"
import "fmt"

// Sheet represents a single Google Spreadsheet and allows certain actions to
// be taken on this spreadsheet.
//
// A Note On Google Sheets Primative Naming:
// 	The Google Sheets API calls an overall document a Spreadsheet.
//
//	This Spreadsheet then contains multiple Sheets. It is in these Sheets
// 	that you enter data.
//
//	To avoid confusion from these similar names this project refers to
//	"Google Spreadsheet" documents as "Sheets".
//
//	In turn "Google Sheets" (the objects you enter the data onto) are
//	called "Pages".
//
// 	The project terminology is a bit confusing because it calls "Google
//	Spreadsheets" (Which contain multiple sub-documents) "Sheets". And it
//	calls "Google Sheets" (The objects you enter data into) "Pages".
//
//	Be aware.
type Sheet struct {
	// svc is the Google API spreadsheet service
	svc *spreadsheet.Service

	// apiSheet is the Google API object used to interact with the
	// spreadsheet
	apiSheet spreadsheet.Spreadsheet

	// apiPage is the Google API object used to interact with the
	// spreadsheet page which will be manipulated
	apiPage *spreadsheet.Sheet
}

// NewSheet creates a new Sheet instance
func NewSheet(svc *spreadsheet.Service, id string, name string) (*Sheet,
	error) {

	// Get apiSheet
	apiSheet, err := svc.FetchSpreadsheet(id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving the sheet API "+
			"object: %s", err.Error())
	}

	// Get apiPage
	apiPage, err := apiSheet.SheetByTitle(name)
	if err != nil {
		return nil, fmt.Errorf("error retrieving the page API object:"+
			" %s", err.Error())
	}

	// Make instance
	return &Sheet{
		svc:      svc,
		apiSheet: apiSheet,
		apiPage:  apiPage,
	}, nil
}
