package integrations

// A Document is a document that can be downloaded and imported
type Document struct {
	DownloadURL string `json:"downloadURL"`
	DocumentID  string `json:"documentID"`
	CategoryID  string `json:"categoryID"`
	PatientID   string `json:"patientID"`
}
