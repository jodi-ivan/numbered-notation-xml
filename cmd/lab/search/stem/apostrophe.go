package stem

// order matters: specific rules first, generic rule last
var ElisionRules = []struct{ Name, Pattern, Replace string }{
	// 1. specific exceptions
	{"elision_cahaya", `(?i)\bcah['’]ya`, `cahaya`},        // cah'ya -> cahaya
	{"elision_sejahtera", `(?i)sejah['’]tra`, `sejahtera`}, // misplaced apostrophe

	// 2. leading apostrophe
	{"elision_nggak", `(?i)(^|[^\p{L}])['’]ngg?a(k?)\b`, `${1}enggak`}, // 'nggak -> enggak
	{"elision_meng", `(?i)(^|[^\p{L}])['’]ng`, `${1}meng`},             // 'ngajak -> mengajak
	{"elision_aku", `(?i)(^|[^\p{L}])['’]ku\b`, `${1}aku`},             // 'Ku -> aku
	{"elision_akan", `(?i)(^|[^\p{L}])['’]kan\b`, `${1}akan`},          // 'kan -> akan

	// 3. generic: letter ' consonant -> letter e consonant
	{"elision_generic", `(?i)(\p{L})['’]([bcdfghjklmnpqrstvwxyz])`, `${1}e${2}`},
}
