package models

type OfficerDelta struct {
	Officer     []Officer `json:"officers"`
	CreatedTime string    `json:"CreatedTime"`
	DeltaAt     string    `json:"delta_at"`
}

type Officer struct {
	CompanyNumber                         string         `json:"company_number"`
	ChangedAt                             string         `json:"changed_at"`
	Kind                                  string         `json:"kind"`
	InternalId                            string         `json:"internal_id"`
	AppointmentDate                       string         `json:"appointment_date"`
	Title                                 string         `json:"title"`
	CorporateInd                          bool           `json:"corporate_ind"`
	Surname                               string         `json:"surname"`
	Forename                              string         `json:"forename"`
	MiddleName                            string         `json:"middle_name"`
	DateOfBirth                           string         `json:"date_of_birth"`
	ServiceAddressSameAsRegisteredAddress bool           `json:"service_address_same_as_registered_address"`
	Nationality                           string         `json:"nationality"`
	Occupation                            string         `json:"occupation"`
	OfficerId                             string         `json:"officer_id"`
	SecureDirector                        bool           `json:"secure_director"`
	OfficerDetailId                       string         `json:"officer_detail_id"`
	OfficerRole                           string         `json:"officer_role"`
	UsualResidentialCountry               string         `json:"usual_residential_country"`
	PreviousNameArray                     PreviousName   `json:"previous_name_array"`
	Identification                        Identification `json:"identification"`
	ServiceAddress                        Address        `json:"service_address"`
	UsualResidentialAddress               Address        `json:"usual_residential_address"`
}

type PreviousName struct {
	PreviousSurname   string `json:"previous_surname"`
	PreviousForename  string `json:"previous_forename"`
	PreviousTimestamp string `json:"previous_timestamp"`
}

type Identification struct {
	EEA Eea `json:"eea"`
}

type Eea struct {
	PlaceRegistered    string `json:"place_registered"`
	RegistrationNumber string `json:"registration_number"`
}

type Address struct {
	Premise                 string `json:"premise"`
	AddressLine1            string `json:"address_line_1"`
	AddressLine2            string `json:"address_line_2"`
	Locality                string `json:"locality"`
	CareOf                  string `json:"care_of"`
	Region                  string `json:"region"`
	PoBox                   string `json:"po_box"`
	SuppliedCompanyName     string `json:"supplied_company_name"`
	Country                 string `json:"country"`
	PostalCode              string `json:"postal_code"`
	UsualCountryOfResidence string `json:"usual_country_of_residence"`
}