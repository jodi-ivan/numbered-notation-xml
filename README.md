# 🎼 Kidung Jemaat Digital — Numbered Notation Engine

*A MusicXML → Numbered Notation renderer with layout logic built from scratch.*

This project aims to digitize and render **Kidung Jemaat** hymns into numbered notation (not angka) based on the Indonesian Yamuger (yayasan musik gereja), following the visual style of the official publication as closely as possible.

The goal is to create a searchable, accessible, and accurate digital version of Kidung Jemaat, starting from **MusicXML** sources and manually curated metadata.

---

## ✨ Current Progress

### ✔ **Content Digitization**
* 41 / 478 hymns fully digitized
* Each hymn translated from MusicXML to custom numbered-notation format
* Metadata and layout rules preserved

### ✔ **Rendering Engine**
Currently the project exposes only the **core engine**, responsible for:

* Parsing MusicXML
* Mapping pitch → numbered notation
* Detecting ties, slurs, note lengths
* Rendering line layout
* Horizontal spacing rules
* Line-break logic
* Multi-verse alignment (first verse baseline complete)

The engine outputs SVG-based notation modeled after the original KJ print style.

---
## 🖼 Screenshot
**SVG modified to have background color
![NumbericNotation](files/var/www/assets/5.svg) 

---

## 🏁 Getting started
- Clone the repository
- Checkout [releases tab](https://github.com/jodi-ivan/numbered-notation-xml/releases)  
- Download these two files:  
    - **kidung-jemaat.db** : the metadata of the music that cannot be stored in the musicxml
    - **musicxml.zip** : musicxml files that needed for the app to run
- Place them somewhere in the drive
- Adjust config in the `files/etc/numbered-mutation-xml/config.ini`
- run the app from `cmd/rest/app.go`
- open browser and open `http//localhost:[port]/kidung-jemaat/render/001` (zero leading numbering, currently from 001 to 022)
> 💡 Alternatively you can download the `goldenfiles.zip` to see the final render looks like. 
---

## 🔧 Features in Progress
### 🔹 Lyric Processing
* Automatic syllable alignment per note
* Verse-by-verse indentation and block formatting
* Handle elisions, melismas, and overlapping notation

### 🔹 Better Content Pipeline
* Improve/automate extraction from MusicXML
* Hybrid manual/automatic verification

### 🔹 Content mangement and discovery
* For searchabilty and categorization

## 📌 Next Features on the Roadmap

### 🎵 4-part SATB Support
* Soprano / Alto / Tenor / Bass layering
* Multi-staff layout
* Proper vertical alignment

### 🎼 Full Musical Notation (Optional Mode)
Switch between:
- Numbered notation
- Traditional staff notation

### 🧩 API + Web Viewer
* Server-side engraving → SVG output
* Web client with simple lyrics + score viewer
* Mobile-friendly optimized layout (future)
---

## 🧠 Why This Project Exists

I started digitizing Kidung Jemaat for personal use, to create a digital hymnal with clean notation and better searchability.
Along the way, it grew into a more general exploration of:

* music engraving algorithms
* layout engines
* type-setting rules for classical hymnbooks
* MusicXML parsing
* text-notation alignment

Development pauses occasionally when I research better automation approaches


## 🤓 References (sites that I used for content and references music theory)
- https://www.hooktheory.com/cheat-sheet
- https://alkitab.sabda.org/resource.php?res=kidung_jemaat
- https://www.musicca.com/dictionary/scales 
