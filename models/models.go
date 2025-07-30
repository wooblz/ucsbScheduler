package models

type PartialApi struct {
    Total int `json:"total"`
    Classes []Class `json:"classes"`
}
type Class struct  {
    CourseID                 string `json:"courseId"` 
    Title                    string `json:"title"`
    SubjectArea              string `json:"subjectArea"`
    ClassSections []Section `json:"classSections"`
}
type Section struct  {
    TimeLocations []TimeLocation `json:"timeLocations`
}
type TimeLocation struct  {
    Room         string `json:"room"`
	Building     string `json:"building"`
	Days         string `json:"days"`
	BeginTime    string `json:"beginTime"`
	EndTime      string `json:"endTime"`
}
