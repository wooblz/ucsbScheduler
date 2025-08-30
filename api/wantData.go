package api

import "github.com/wooblz/ucsbScheduler/models"

var Solution1 = []models.Class{
    {
        CourseID:    "CMPSC     5A ",
        Title:       "INTRO DATA SCI 1",
        SubjectArea: "CMPSC   ",
        ClassSections: []models.Section{
            {
                TimeLocations: []models.TimeLocation{
                    {
                        Room:      "1610",
                        Building:  "BRDA",
                        Days:      "M W    ",
                        BeginTime: "15:30",
                        EndTime:   "16:45",
                    },
                },
            },
        },
    },
}
var Solution2 = []models.Class{
    {
        CourseID:    "CMPSC    16  ",
        Title:       "PROBLEM SOLVING I",
        SubjectArea: "CMPSC   ",
        ClassSections: []models.Section{
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "1101", Building: "ILP", Days: "M W    ", BeginTime: "14:00", EndTime: "15:15"},
                },
            },
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "3525", Building: "PHELP", Days: "   R   ", BeginTime: "09:00", EndTime: "09:50"},
                },
            },
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "3525", Building: "PHELP", Days: "   R   ", BeginTime: "10:00", EndTime: "10:50"},
                },
            },
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "3525", Building: "PHELP", Days: "   R   ", BeginTime: "11:00", EndTime: "11:50"},
                },
            },
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "3525", Building: "PHELP", Days: "   R   ", BeginTime: "12:00", EndTime: "12:50"},
                },
            },
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "3525", Building: "PHELP", Days: "   R   ", BeginTime: "13:00", EndTime: "13:50"},
                },
            },
        },
    },
    {
        CourseID:    "CMPSC    24  ",
        Title:       "PROBLEM SOLVING II",
        SubjectArea: "CMPSC   ",
        ClassSections: []models.Section{
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "1701", Building: "TD-W", Days: "M W    ", BeginTime: "11:00", EndTime: "12:15"},
                },
            },
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "3525", Building: "PHELP", Days: "   R   ", BeginTime: "14:00", EndTime: "14:50"},
                },
            },
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "3525", Building: "PHELP", Days: "   R   ", BeginTime: "15:00", EndTime: "15:50"},
                },
            },
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "3525", Building: "PHELP", Days: "   R   ", BeginTime: "16:00", EndTime: "16:50"},
                },
            },
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "3525", Building: "PHELP", Days: "   R   ", BeginTime: "17:00", EndTime: "17:50"},
                },
            },
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "3525", Building: "PHELP", Days: "   R   ", BeginTime: "18:00", EndTime: "18:50"},
                },
            },
        },
    },
    {
        CourseID:    "MATH    190PS",
        Title:       "PROBLEM SOLVING",
        SubjectArea: "MATH    ",
        ClassSections: []models.Section{
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "1508", Building: "PHELP", Days: "M W F  ", BeginTime: "09:00", EndTime: "09:50"},
                },
            },
        },
    },
    {
        CourseID:    "MATH CS 101A ",
        Title:       "PROBLEM SOLVING I",
        SubjectArea: "MATH CS ",
        ClassSections: []models.Section{
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "164B", Building: "CRST", Days: " T R   ", BeginTime: "14:00", EndTime: "15:15"},
                },
            },
        },
    },
}
var Solution3 = []models.Class{
    {
        CourseID:    "CMPSC    16  ",
        Title:       "PROBLEM SOLVING I",
        SubjectArea: "CMPSC   ",
        Room:        "1101",
        Building:    "ILP",
        Days:        "M W    ",
        BeginTime:   "14:00",
        EndTime:     "15:15",
        ClassSections: []models.Section{
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "3525", Building: "PHELP", Days: "   R   ", BeginTime: "09:00", EndTime: "09:50"},
                },
            },
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "3525", Building: "PHELP", Days: "   R   ", BeginTime: "10:00", EndTime: "10:50"},
                },
            },
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "3525", Building: "PHELP", Days: "   R   ", BeginTime: "11:00", EndTime: "11:50"},
                },
            },
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "3525", Building: "PHELP", Days: "   R   ", BeginTime: "12:00", EndTime: "12:50"},
                },
            },
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "3525", Building: "PHELP", Days: "   R   ", BeginTime: "13:00", EndTime: "13:50"},
                },
            },
        },
    },
    {
        CourseID:    "CMPSC    24  ",
        Title:       "PROBLEM SOLVING II",
        SubjectArea: "CMPSC   ",
        Room:        "1701",
        Building:    "TD-W",
        Days:        "M W    ",
        BeginTime:   "11:00",
        EndTime:     "12:15",
        ClassSections: []models.Section{
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "3525", Building: "PHELP", Days: "   R   ", BeginTime: "14:00", EndTime: "14:50"},
                },
            },
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "3525", Building: "PHELP", Days: "   R   ", BeginTime: "15:00", EndTime: "15:50"},
                },
            },
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "3525", Building: "PHELP", Days: "   R   ", BeginTime: "16:00", EndTime: "16:50"},
                },
            },
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "3525", Building: "PHELP", Days: "   R   ", BeginTime: "17:00", EndTime: "17:50"},
                },
            },
            {
                TimeLocations: []models.TimeLocation{
                    {Room: "3525", Building: "PHELP", Days: "   R   ", BeginTime: "18:00", EndTime: "18:50"},
                },
            },
        },
    },
    {
        CourseID:    "MATH    190PS",
        Title:       "PROBLEM SOLVING",
        SubjectArea: "MATH    ",
        Room:        "1508",
        Building:    "PHELP",
        Days:        "M W F  ",
        BeginTime:   "09:00",
        EndTime:     "09:50",
        ClassSections: []models.Section{},
    },
    {
        CourseID:    "MATH CS 101A ",
        Title:       "PROBLEM SOLVING I",
        SubjectArea: "MATH CS ",
        Room:        "164B",
        Building:    "CRST",
        Days:        " T R   ",
        BeginTime:   "14:00",
        EndTime:     "15:15",
        ClassSections: []models.Section{},
    },
}

var FinalSolution1 = models.Final {
    HasFinals: true,
    Comments:  "",
    ExamDay:   "R",
    ExamDate:  "20250320",
    BeginTime: "12:00",
    EndTime:   "15:00",
}
