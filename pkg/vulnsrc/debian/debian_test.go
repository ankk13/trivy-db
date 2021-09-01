package debian_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aquasecurity/trivy-db/pkg/db"
	"github.com/aquasecurity/trivy-db/pkg/dbtest"
	"github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/aquasecurity/trivy-db/pkg/vulnsrc/debian"
)

func TestVulnSrc_Update(t *testing.T) {
	type wantKV struct {
		key   []string
		value interface{}
	}
	tests := []struct {
		name       string
		dir        string
		wantValues []wantKV
		wantErr    string
	}{
		{
			name: "happy path",
			dir:  filepath.Join("testdata", "happy"),
			wantValues: []wantKV{
				// Ref. https://security-tracker.debian.org/tracker/CVE-2021-33560
				{
					key: []string{"advisory-detail", "CVE-2021-33560", "debian 9", "libgcrypt20"},
					value: types.Advisory{
						VendorIDs:    []string{"DLA-2691-1"},
						FixedVersion: "1.7.6-2+deb9u4",
					},
				},
				{
					key: []string{"advisory-detail", "CVE-2021-33560", "debian 10", "libgcrypt20"},
					value: types.Advisory{
						FixedVersion: "1.8.4-5+deb10u1",
					},
				},
				{
					key: []string{"advisory-detail", "CVE-2021-33560", "debian 11", "libgcrypt20"},
					value: types.Advisory{
						FixedVersion: "1.8.7-6",
					},
				},
				{
					key: []string{"advisory-detail", "CVE-2021-33560", "debian 11", "libgcrypt20"},
					value: types.Advisory{
						FixedVersion: "1.8.7-6",
					},
				},
				{
					key: []string{"advisory-detail", "DSA-3714-1", "debian 8", "akonadi"},
					value: types.Advisory{
						VendorIDs:    []string{"DSA-3714-1"},
						FixedVersion: "1.13.0-2+deb8u2",
					},
				},
			},
		},
		{
			name:    "sad broken distributions",
			dir:     filepath.Join("testdata", "broken-distributions"),
			wantErr: "failed to decode Debian distribution JSON",
		},
		{
			name:    "sad broken packages",
			dir:     filepath.Join("testdata", "broken-packages"),
			wantErr: "failed to decode testdata/broken-packages/",
		},
		{
			name:    "sad broken CVE",
			dir:     filepath.Join("testdata", "broken-cve"),
			wantErr: "json decode error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := dbtest.InitTestDB(t, nil)
			dbPath := db.Path(tmpDir)

			vs := debian.NewVulnSrc()

			err := vs.Update(tt.dir)
			if tt.wantErr != "" {
				require.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			require.NoError(t, err)
			require.NoError(t, db.Close())

			for _, want := range tt.wantValues {
				dbtest.JSONEq(t, dbPath, want.key, want.value)
			}
		})
	}
}
