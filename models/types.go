package models

import "text/template"

var (
	// Datum contains all data from csv file
	D Data

	// T contains all parsed templates
	T *template.Template

	err error

	//represents ftp uri
	FTPURI string

	//represents ftp login
	FTPLogin string

	//represents ftp password
	FTPPassword string

	//represents server port
	Port string
)

// Data represents data that is passed to templates
type Data struct {
	Datum []CDRModified
	TC    TotalCharged
}

// TotalCharged represents 'Charged duration in minutes summed up by categories'
type TotalCharged struct {
	FixedToMobile    float64
	International    float64
	National         float64
	IntercapitalCity float64
}

// CDRModified represents a data from a cdr file, but with some columns removed
type CDRModified struct {
	Five        string  `csv:"Connect Datetime" bson:"Connect Datetime"`
	Six         string  `csv:"Disconnect Datetime" bson:"Disconnect Datetime"`
	Ten         int     `csv:"Charged Duration (Seconds)" bson:"Charged Duration (Seconds)"`
	Eleven      float64 `csv:"Charged Duration (Minutes)" bson:"Charged Duration (Minutes)"`
	Thirteen    int     `csv:"Calling Number" bson:"Calling Number"`
	Nineteen    int     `csv:"Called Number" bson:"Called Number"`
	Twenty      string  `csv:"Called Number Location" bson:"Called Number Location"`
	TwentyOne   string  `csv:"Location Pair Category" bson:"Location Pair Category"`
	TwentyTwo   float64 `csv:"Charged Amount" bson:"Charged Amount"`
	TwentyThree string  `csv:"Currency Code" bson:"Currency Code"`
	TwentyFive  int     `csv:"Completion Code ID" bson:"Completion Code ID"`
	TwentySix   string  `csv:"Completion Code Name" bson:"Completion Code Name"`
}

// LoginInfo represents a username-password data from a CDR MongoDB database
type LoginInfo struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
}
