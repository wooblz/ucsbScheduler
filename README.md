# 📅 UCSB Class Scheduler

> A fast, full-stack course schedule visualizer and iCalendar (`.ics`) generator built for UC Santa Barbara Gauchos.

[![Go Version](https://img.shields.io/github/go-mod/go-version/wooblz/ucsbScheduler?style=flat-square)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Full--Text%20Search-336791?style=flat-square&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![React](https://img.shields.io/badge/React-18-61DAFB?style=flat-square&logo=react&logoColor=black)](https://react.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)

<!-- IMAGE PLACEHOLDER 1: HERO / MAIN PREVIEW -->
<!-- Recommended: Take a full-width desktop screenshot showing selected classes and the weekly schedule grid populated. -->

<p align="center">
  <img width="1774" height="901" alt="image" src="https://github.com/user-attachments/assets/78ff85c2-9790-4187-84ae-4bba60dedf41" />
</p>

---

## 📖 Overview

UCSB Class Scheduler streamlines class registration and planning for UC Santa Barbara students. Instead of manually copying lecture hours, section times, room locations, and final exam dates into your calendar, UCSB Scheduler lets you search courses in real-time, preview your weekly grid, and export a standards-compliant `.ics` file that imports directly into **Google Calendar**, **Apple Calendar**, **Outlook**, and more.

### Key Highlights
- **Direct UCSB API Integration**: Ingests curriculum data and final exam schedules directly from the official UCSB Developer API.
- **Fast Full-Text Search**: Powered by PostgreSQL `tsvector` weighted indexing (`course_id` and `title`) and ILIKE fallback for instant query results.
- **Weekly Schedule Visualizer**: Live interactive timetable grid showing lectures, discussion sections, and final exams.
- **Standards-Compliant iCal Export**: Generates `.ics` files featuring exact recurrence rules (`RRULE: FREQ=WEEKLY`), Pacific Time offsets, and individual final exam dates.
- **Light & Dark Mode**: Seamless theme switching with persistent CSS variables.
- **Automated Catalog Sync**: Built-in background worker that refreshes course and section databases every 24 hours.

---

## 📸 Screenshots & Demos

### Weekly Schedule Grid & Light / Dark Mode
<!-- IMAGE PLACEHOLDER 2: LIGHT vs DARK MODE -->
<!-- Recommended: A side-by-side image or a GIF showing the Light and Dark theme toggle. -->
<p align="center">
  <img width="1439" height="852" alt="image" src="https://github.com/user-attachments/assets/5bfd39b6-9899-4e3f-9bc2-2fc17ab0d766" />

</p>

### Interactive Search & Section Selection
<!-- IMAGE PLACEHOLDER 3: SEARCH & SELECTION -->
<!-- Recommended: Crop showing search results for a course (e.g. CMPSC 16) with lecture times and section enroll buttons. -->
<p align="center">
  <img width="1309" height="375" alt="image" src="https://github.com/user-attachments/assets/3e1a61a2-f068-4ce2-818a-7a7ffe16e593" />

</p>

### Final Result in Calendar App
<!-- IMAGE PLACEHOLDER 4: IMPORTED CALENDAR VIEW -->
<!-- Recommended: A screenshot showing the imported schedule in Google Calendar or Apple Calendar. -->
<p align="center">
  <img width="1274" height="498" alt="image" src="https://github.com/user-attachments/assets/84d8b7b9-b54d-44af-8fcd-adaf06d003c5" />

</p>

---

## 🛠️ Architecture & Tech Stack

```text
┌────────────────┐        HTTP / JSON         ┌─────────────────────────────────┐
│                ├───────────────────────────►│            Go Server            │
│  React 18 SPA  │                            │         (net/http API)          │
│ (static/index) │◄───────────────────────────┤                                 │
└────────────────┘   Live Preview & .ics DL   └───────┬─────────────────┬───────┘
                                                      │                 │
                                            Curriculum Sync        PostgreSQL
                                                      ▼                 ▼
                                              ┌───────────────┐ ┌───────────────┐
                                              │   UCSB API    │ │ Full-Text     │
                                              │ Curriculums & │ │ Search & JSON │
                                              │    Finals     │ │ Aggregation   │
                                              └───────────────┘ └───────────────┘
