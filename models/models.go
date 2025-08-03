package models

type PartialApi struct {
    Message string `json:"message"`
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
    TimeLocations []TimeLocation `json:"timeLocations"`
}
type TimeLocation struct  {
    Room         string `json:"room"`
	Building     string `json:"building"`
	Days         string `json:"days"`
	BeginTime    string `json:"beginTime"`
	EndTime      string `json:"endTime"`
}
type Final struct {
	HasFinals bool   `json:"hasFinals"`
	Comments  string `json:"comments"`
	ExamDay   string `json:"examDay"`
	ExamDate  string `json:"examDate"`
	BeginTime string `json:"beginTime"`
	EndTime   string `json:"endTime"`
}
