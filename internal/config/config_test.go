package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{
			name:    "load valid config",
			want:    36,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := LoadConfig()

			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(config.Abbreviations) != tt.want {
				t.Errorf("LoadConfig() got %d abbreviations, want %d", len(config.Abbreviations), tt.want)
			}
		})
	}
}

func TestAbbreviationStructure(t *testing.T) {
	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(config.Abbreviations) == 0 {
		t.Fatal("No abbreviations found in config")
	}

	abbr := config.Abbreviations[0]
	if abbr.Name == "" {
		t.Error("First abbreviation should have a name")
	}
	if abbr.Abbr == "" && (abbr.Options == nil || abbr.Options.Regex == "") {
		t.Error("Abbreviation should have either abbr or regex")
	}
	if abbr.Snippet == "" {
		t.Error("Abbreviation should have a snippet")
	}
}
