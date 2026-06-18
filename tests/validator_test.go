package tests

import (
	"testing"

	vfs "github.com/lordofscripts/govault/internal/fs"
)

func TestValidator_ValidateFilename(t *testing.T) {
	CountriesISO := vfs.KeyValuePairArray{
		vfs.KeyValuePair{Key: "NL", Value: "Nederland"},
		vfs.KeyValuePair{Key: "NLD", Value: "Nederland"},
		vfs.KeyValuePair{Key: "PA", Value: "Panama"},
		vfs.KeyValuePair{Key: "PAN", Value: "Panama"},
		vfs.KeyValuePair{Key: "US", Value: "U.S.A."},
		vfs.KeyValuePair{Key: "USA", Value: "U.S.A."},
	}
	ServerRole := vfs.KeyValuePairArray{
		vfs.KeyValuePair{Key: "PROD", Value: "Production"},
		vfs.KeyValuePair{Key: "STAGE", Value: "Staging"},
	}
	ImagePrefixes := vfs.KeyValuePairArray{
		vfs.KeyValuePair{Key: "IMG", Value: "Standard image"},
		vfs.KeyValuePair{Key: "PIC", Value: "Drawing"},
		vfs.KeyValuePair{Key: "IMGX", Value: "Encrypted image"},
		vfs.KeyValuePair{Key: "SCN", Value: "Scanned picture"},
	}
	ImageExt := vfs.KeyValuePairArray{
		vfs.KeyValuePair{Key: "jpg", Value: "JPEG"},
		vfs.KeyValuePair{Key: "png", Value: "PNG"},
	}

	// better use NewValidator() followed with the fluent API methods.
	v := &vfs.Validator{
		Lists: map[string]vfs.StringMap{
			"L1": {"LOG": "Logfile", "ERR": "Errorfile"},
		},
		Regexes: vfs.StringMap{
			"R1": `[a-z]{3}`,
			"R2": `\d{6}`,
			//"R3": `jpg|jpeg|png|tiff|gif|psd|pdn`, // rule for image file extensions
		},
		Predefined: map[string]vfs.KeyValuePairArray{
			"C": CountriesISO,
			"S": ServerRole,
			//"I": ImagePrefixes,
			"X": ImageExt,
		},
	}

	// exercise the fluent API methods
	v.AddRuleItem(3, `jpg|jpeg|png|tiff|gif|psd|pdn|xcf`). // becomes $R3
								AddRuleItem(4, `\p{L}\d[a-d]`).                                                    // becomes $R4 (Unicode character followed by digit and single lowercase letter)
								AddPredefinedItem('I', ImagePrefixes).                                             // becomes %I
								AddListItem(2, vfs.StringMap{"PAS": "Passport", "DNI": "National Identification"}) // becomes $L2

	testCases := []struct {
		Template string
		Filename string
		Pass     bool
	}{
		// Test case: List, Date, Time, Regex, Wildcard
		{Template: "app_$L1_%D_%T_$R1_*.txt", Filename: "app_LOG_20231025_143005_abc_session99.txt", Pass: true},
		{Template: "app_$L1_%D_%T_$R1_*.txt", Filename: "app_ERR_20231025_143005_abc_session99.txt", Pass: true},
		{Template: "%S_%D_%T_$L1_$R1_*.csv", Filename: "PROD_20231201_235959_ERR_xyz_export_v1.csv", Pass: true},
		{Template: "%S_%D_%T_$L1_$R1_*.csv", Filename: "STAGE_20231201_235959_LOG_abc_export_v1.csv", Pass: true},
		{Template: "%S_%D_%T_$L1_$R1_*.csv", Filename: "STAGE_20231201_235959_REPORT_42_export_v1.csv", Pass: false},
		{Template: "%C_$L2_%D_%T_$R1_*.pdf", Filename: "NL_PAS_20231025_143005_abc_session99.pdf", Pass: true},
		{Template: "%C_$L2_%D_%T_$R1_*.pdf", Filename: "NLD_DNI_20231025_143005_abc_session99.pdf", Pass: true},
		{Template: "%C_$L2_%D_%T_$R1_*.pdf", Filename: "PA_DNI_20231025_143005_abc_session99.pdf", Pass: true},
		{Template: "%C_$L2_%D_%T_$R1_*.pdf", Filename: "PAN_PAS_20231025_143005_abc_session99.pdf", Pass: true},
		{Template: "%C_$L2_%D_%T_$R1_*.pdf", Filename: "US_PAS_20231025_143005_abc_session99.pdf", Pass: true},
		{Template: "%C_$L2_%D_%T_$R1_*.pdf", Filename: "USA_DNI_20231025_143005_abc_session99.pdf", Pass: true},
		{Template: "%C_$L2_%D_%T_$R1_*.pdf", Filename: "DE_LOG_20231025_143005_abc_session99.pdf", Pass: false},

		{Template: "%C_$L2_%D_%T_$R1_*.pdf", Filename: "NL_PAS_20231025_143005_Σ3c_session99.pdf", Pass: false},
		{Template: "%C_$L2_%D_%T_$R1_*.pdf", Filename: "NLD_DNI_20231025_143005_Ω2m_session99.pdf", Pass: false},
		{Template: "%C_$L2_%D_%T_$R4_*.pdf", Filename: "NLD_DNI_20231025_143005_Ω2m_session99.pdf", Pass: false},
		{Template: "%C_$L2_%D_%T_$R4_*.pdf", Filename: "NL_PAS_20231025_143005_Σ3c_session99.pdf", Pass: true},
		{Template: "%C_$L2_%D_%T_$R4_*.pdf", Filename: "NLD_DNI_20231025_143005_Ω2d_session99.pdf", Pass: true},
		// Test case: images Prefix, Date, Time, Rule, Free-form
		{Template: "IMG_%D_%T_$R2*.jpg", Filename: "IMG_20231025_143005_123456.jpg", Pass: true},
		{Template: "IMG_%D_%T_$R2*.jpg", Filename: "IMG_20231025_143005_123456_free.jpg", Pass: true},
		{Template: "IMG_%D_%T_$R2*.$R3", Filename: "IMG_20231025_143005_123456_free.gif", Pass: true},
		// Test case: image prefixes
		{Template: "%I_%D_%T_$R2*.jpg", Filename: "IMG_20231025_143005_123456_free.jpg", Pass: true},
		{Template: "%I_%D_%T_$R2*.jpg", Filename: "IMGX_20231025_143005_123456_free.jpg", Pass: true},
		{Template: "%I_%D_%T_$R2*.jpg", Filename: "PIC_20231025_143005_123456_free.jpg", Pass: true},
		{Template: "%I_%D_%T_$R2*.jpg", Filename: "SCN_20231025_143005_123456_free.jpg", Pass: true},
		{Template: "%I_%D_%T_$R2*.jpg", Filename: "FAX_20231025_143005_123456_free.jpg", Pass: false},

		{Template: "%I_%D_%T*.$R3", Filename: "IMG_20250513_143731_ed.jpg", Pass: true},
		{Template: "%I_%D_%T*.%X", Filename: "PIC_20250513_143731_ed.png", Pass: true},
	}

	for nr, tcase := range testCases {
		isValid, err := v.ValidateFilename(tcase.Template, tcase.Filename)
		if err != nil {
			t.Errorf("#%d %v\n", nr+1, err)
		} else if isValid != tcase.Pass {
			t.Errorf("#%d validation failed, exp:%t got:%t\n", nr+1, tcase.Pass, isValid)
		}

	}
}
