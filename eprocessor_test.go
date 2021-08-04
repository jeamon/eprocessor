package main

import "testing"

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
