package main

import (
	"reflect"
	"testing"
)

func TestExtractFilename(t *testing.T) {
	// url is valid web link provided - expected is the expected filename
	casesTable := []struct {
		url  string
		want string
	}{
		{"https://some-link/some-parts/xxxx/data.csv", "data.csv"},
		{"https://some-link/some-parts/another-parts/infos.csv", "infos.csv"},
		{"https://some-link/some-parts/another-parts/data", "data"},
	}

	for _, c := range casesTable {
		got := ExtractFilename(c.url)
		if got != c.want {
			t.Errorf("Extraction for %q was incorrect, got: %q, wanted %q", c.url, got, c.want)
		}
	}
}

func TestRemoveMemoField(t *testing.T) {

	ImportDate := "08/04/2021"

	input := [][]string{

		{"Date", "Name", "Address", "Address2", "City", "State", "Zipcode", "Telephone", "Mobile", "Amount", "Processor", "Memo"},

		{"01/04/2016", "Jerome AMON", "Poland Street", "Poland Street 2", "Warsaw", "PL", "38002", "", "000-000-0000", "$90", "Stripe", "memo infos"},

		{"01/04/2017", "Jerome AMON", "Poland Street", "Poland Street 2", "Warsaw", "PL", "38002", "", "000-000-0000", "$90", "Stripe", "memo infos"},

		{"01/04/2018", "Abou AMON", "Poland Street", "Poland Street 2", "Warsaw", "PL", "38002", "", "000-000-0000", "$90", "Stripe", "memo infos"},

		{"01/04/2019", "Abou AMON", "Poland Street", "Poland Street 2", "Krakow", "PL", "38002", "", "000-000-0000", "$90", "Stripe", "memo infos"},
	}

	// last field of each record which is Memo field should be removed.
	expected := [][]string{

		{"Date", "Name", "Address", "Address2", "City", "State", "Zipcode", "Telephone", "Mobile", "Amount", "Processor", "08/04/2021"},

		{"01/04/2016", "Jerome AMON", "Poland Street", "Poland Street 2", "Warsaw", "PL", "38002", "", "000-000-0000", "$90", "Stripe", "08/04/2021"},

		{"01/04/2017", "Jerome AMON", "Poland Street", "Poland Street 2", "Warsaw", "PL", "38002", "", "000-000-0000", "$90", "Stripe", "08/04/2021"},

		{"01/04/2018", "Abou AMON", "Poland Street", "Poland Street 2", "Warsaw", "PL", "38002", "", "000-000-0000", "$90", "Stripe", "08/04/2021"},

		{"01/04/2019", "Abou AMON", "Poland Street", "Poland Street 2", "Krakow", "PL", "38002", "", "000-000-0000", "$90", "Stripe", "08/04/2021"},
	}
	// process the input data.
	RemoveMemoField(&input, ImportDate)
	// compare input data and expected state.
	if reflect.DeepEqual(input, expected) == false {
		t.Errorf("after processing. got input content different from expected content.")
	}
}

func TestReplaceEmptyValues(t *testing.T) {

	input := [][]string{

		{"01/04/2016", "Jerome AMON", "Poland Street", "", "Warsaw", "PL", "38002", "", "000-000-0000", "$90", "Stripe", "08/04/2021"},

		{"01/04/2017", "Jerome AMON", "Poland Street", "", "Warsaw", "PL", "38002", "", "000-000-0000", "$90", "Stripe", "08/04/2021"},

		{"01/04/2018", "Abou AMON", "Poland Street", "", "Warsaw", "PL", "38002", "", "000-000-0000", "$90", "Stripe", "08/04/2021"},

		{"01/04/2019", "Abou AMON", "Poland Street", "missing", "Krakow", "PL", "38002", "", "000-000-0000", "$90", "Stripe", "08/04/2021"},
	}

	// counting from 1. we expect to see fields 4th (Address2) and 8th (Telephone) got value "missing".
	expected := [][]string{

		{"01/04/2016", "Jerome AMON", "Poland Street", "missing", "Warsaw", "PL", "38002", "missing", "000-000-0000", "$90", "Stripe", "08/04/2021"},

		{"01/04/2017", "Jerome AMON", "Poland Street", "missing", "Warsaw", "PL", "38002", "missing", "000-000-0000", "$90", "Stripe", "08/04/2021"},

		{"01/04/2018", "Abou AMON", "Poland Street", "missing", "Warsaw", "PL", "38002", "missing", "000-000-0000", "$90", "Stripe", "08/04/2021"},

		{"01/04/2019", "Abou AMON", "Poland Street", "missing", "Krakow", "PL", "38002", "missing", "000-000-0000", "$90", "Stripe", "08/04/2021"},
	}
	// process the input data.
	ReplaceEmptyValues(&input)
	// compare input data and expected state.
	if reflect.DeepEqual(input, expected) == false {
		t.Errorf("after processing. got input content different from expected content.")
	}
}

func TestRemoveDuplicateRecords(t *testing.T) {

	input := [][]string{

		{"01/04/2016", "Jerome AMON", "Poland Street", "missing", "Warsaw", "PL", "38002", "missing", "000-000-0000", "$90", "Stripe", "08/04/2021"},
		{"01/04/2016", "Jerome AMON", "Poland Street", "missing", "Warsaw", "PL", "38002", "missing", "000-000-0000", "$90", "Stripe", "08/04/2021"},
		{"01/04/2016", "Jerome AMON", "Poland Street", "missing", "Warsaw", "PL", "38002", "missing", "000-000-0000", "$90", "Stripe", "08/04/2021"},

		{"01/04/2017", "Jerome AMON", "Poland Street", "missing", "Warsaw", "PL", "38002", "missing", "000-000-0000", "$90", "Stripe", "08/04/2021"},
		{"01/04/2017", "Jerome AMON", "Poland Street", "missing", "Warsaw", "PL", "38002", "missing", "000-000-0000", "$90", "Stripe", "08/04/2021"},

		{"01/04/2018", "Abou AMON", "Poland Street", "missing", "Warsaw", "PL", "38002", "missing", "000-000-0000", "$90", "Stripe", "08/04/2021"},
		{"01/04/2018", "Abou AMON", "Poland Street", "missing", "Warsaw", "PL", "38002", "missing", "000-000-0000", "$90", "Stripe", "08/04/2021"},

		{"01/04/2019", "Abou AMON", "Poland Street", "missing", "Krakow", "PL", "38002", "missing", "000-000-0000", "$90", "Stripe", "08/04/2021"},
		{"01/04/2019", "Abou AMON", "Poland Street", "missing", "Krakow", "PL", "38002", "missing", "000-000-0000", "$90", "Stripe", "08/04/2021"},

		{"01/04/2016", "Abou AMON", "Poland Street", "missing", "Warsaw", "PL", "38002", "missing", "000-000-0000", "$90", "Stripe", "08/04/2021"},
	}

	mapOfRecords := make(map[Record]struct{})
	want := 5
	got := RemoveDuplicateRecords(&input, mapOfRecords)

	if got != want {
		t.Errorf("got %d, wanted: %d", got, want)
	}
}

func TestToJson(t *testing.T) {

	r := &Record{
		Date:       "08/02/2019",
		Name:       "Jerome A.",
		Address:    "0000 Krakow",
		Address2:   "missing",
		City:       "Krakow",
		State:      "Lesser Poland",
		Zipcode:    "00-000",
		Telephone:  "000-000-0000",
		Mobile:     "504-319-6911",
		Amount:     "$14",
		Processor:  "PayPal",
		ImportDate: "08/02/2021",
	}

	got := r.RecordToJson()
	want := `{"date":"08/02/2019","name":"Jerome A.","address":"0000 Krakow","address2":"missing","city":"Krakow","state":"Lesser Poland","zipcode":"00-000","telephone":"000-000-0000","mobile":"504-319-6911","amount":"$14","processor":"PayPal","importdate":"08/02/2021"}`

	if got != want {
		t.Errorf("got %q, wanted: %q", got, want)
	}
}

func ExamplePause_exit() {
	Pause("exit")
	// Output:
	// {:} Press [Enter] key to exit
}

func ExamplePause_continue() {
	Pause("continue")
	// Output:
	// {:} Press [Enter] key to continue
}

func ExampleBanner() {
	Banner()
	// Output:
	// ///////////////////////////////////////////////////////////////////////////////////
	// @@@@@@@@@@@@@@@@@@@@ E-COMPANY TOOL // CSV FILE PROCESSOR v1.0 @@@@@@@@@@@@@@@@@@@@
	// ///////////////////////////////////////////////////////////////////////////////////
}
