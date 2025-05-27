<h1 align="center"><strong><em>BioTree</em></strong></h1>

**BioTree** is a powerful and intuitive web application built with **Go**, **JavaScript**, **HTML**, and **CSS**, that allows users to upload documents in `.pdf`, `.txt`, or `.docx` format and generates dynamic **SVG visualizations** based on the content. Whether you're analyzing biological data, scientific text, or general structured information, BioTree offers a fast, clean, and interactive way to visualize it.

---

## Features

-  Upload support for **PDF**, **TXT**, and **DOCX** files  
-  Clean **SVG output** for scalable and shareable visuals  
-  Written in **Go** (backend), **JavaScript**, **HTML**, and **CSS** (frontend)  
-  Lightweight and fast – ideal for local and web deployments  
-  Simple and clean UI with responsive design

---

##  Input Formats

BioTree currently supports the following file types for input:

- `.pdf` – Portable Document Format
- `.txt` – Plain text files
- `.docx` – Microsoft Word Open XML Format

These are parsed and converted into a structured format to produce a meaningful SVG chart or visualization.

---

##  Output

The output is a valid **SVG string**, rendered in the browser or downloadable as an image file. The SVGs are:

- Scalable (perfect for high-resolution)
- Styled and structured
- Reflective of the hierarchy or data within the document (e.g., sections, keywords, or tree structures)

---

##  How It Works

1. **Frontend (JavaScript + HTML + CSS)**  
   - File upload UI  
   - Preview and interactive SVG rendering  
   - Responsive layout with smooth transitions

2. **Backend (Go)**  
   - File handling and content extraction (PDF, DOCX, TXT)  
   - Data parsing and transformation logic  
   - SVG generation logic  
   - REST API for file upload and SVG response

---

##  Getting Started

### Prerequisites

- Go 1.20+
- Node.js (optional for frontend build tools)
- `libreoffice` (if you're using a DOCX to TXT converter externally)

### Installation

```bash
git clone https://github.com/Vinolia-E/BioTree.git
cd BioTree
go run main.go
