package folder_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folder"
	"github.com/gofrs/uuid"
	// "github.com/stretchr/testify/assert"
)

// feel free to change how the unit test is structured
func Test_folder_GetFoldersByOrgID(t *testing.T) {
	t.Parallel()

	org1 := uuid.Must(uuid.NewV4())
	org2 := uuid.Must(uuid.NewV4())
	org3 := uuid.Must(uuid.NewV4())

	// prepare sample folders
	sampleFolders := []folder.Folder{
		{Name: "alpha", Paths: "alpha", OrgId: org1},
		{Name: "bravo", Paths: "alpha.bravo", OrgId: org1},
		{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: org1},
		{Name: "delta", Paths: "delta", OrgId: org2},
		{Name: "echo", Paths: "delta.echo", OrgId: org2},
		{Name: "foxtrot", Paths: "foxtrot", OrgId: org2},
	}

	tests := [...]struct {
		name    string
		orgID   uuid.UUID
		folders []folder.Folder
		want    []folder.Folder
	}{
		{
			name:    "Get folders for org1",
			orgID:   org1,
			folders: sampleFolders,
			want: []folder.Folder{
				{Name: "alpha", Paths: "alpha", OrgId: org1},
				{Name: "bravo", Paths: "alpha.bravo", OrgId: org1},
				{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: org1},
			},
		},
		{
			name:    "Get folders for org2",
			orgID:   org2,
			folders: sampleFolders,
			want: []folder.Folder{
				{Name: "delta", Paths: "delta", OrgId: org2},
				{Name: "echo", Paths: "delta.echo", OrgId: org2},
				{Name: "foxtrot", Paths: "foxtrot", OrgId: org2},
			},
		},
		{
			name:    "No folders for org3",
			orgID:   org3,
			folders: sampleFolders,
			want:    []folder.Folder{},
		},
	}
	for _, tt := range tests {
		tt := tt // ensures correct value in closure (anonymous function)
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(tt.folders)

			got := f.GetFoldersByOrgID(tt.orgID)

			if len(got) != len(tt.want) {
				t.Errorf("Expected %d folders, got %d", len(tt.want), len(got))
			}

			for _, expectedFolder := range tt.want {
				found := false
				for _, actualFolder := range got {
					if actualFolder.Name == expectedFolder.Name &&
						actualFolder.Paths == expectedFolder.Paths &&
						actualFolder.OrgId == expectedFolder.OrgId {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected folder %+v not found in result", expectedFolder)
				}
			}
		})
	}
}

func Test_folder_GetAllChildFolders(t *testing.T) {
	t.Parallel()

	org1 := uuid.Must(uuid.NewV4())
	org2 := uuid.Must(uuid.NewV4())

	// prepare sample folders
	sampleFolders := []folder.Folder{
		{Name: "alpha", Paths: "alpha", OrgId: org1},
		{Name: "bravo", Paths: "alpha.bravo", OrgId: org1},
		{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: org1},
		{Name: "delta", Paths: "alpha.delta", OrgId: org1},
		{Name: "echo", Paths: "echo", OrgId: org1},
		{Name: "foxtrot", Paths: "foxtrot", OrgId: org2},
	}

	tests := []struct {
		name         string
		orgID        uuid.UUID
		parentFolder string
		expected     []folder.Folder
		expectError  bool
		errorMessage string
	}{
		{
			name:         "Get all child folders of alpha",
			orgID:        org1,
			parentFolder: "alpha",
			expected: []folder.Folder{
				{Name: "bravo", Paths: "alpha.bravo", OrgId: org1},
				{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: org1},
				{Name: "delta", Paths: "alpha.delta", OrgId: org1},
			},
			expectError: false,
		},
		{
			name:         "Get all child folders of bravo",
			orgID:        org1,
			parentFolder: "bravo",
			expected: []folder.Folder{
				{Name: "charlie", Paths: "alpha.bravo.charlie", OrgId: org1},
			},
			expectError: false,
		},
		{
			name:         "No child folders for charlie",
			orgID:        org1,
			parentFolder: "charlie",
			expected:     []folder.Folder{},
			expectError:  false,
		},
		{
			name:         "No child folders for echo",
			orgID:        org1,
			parentFolder: "echo",
			expected:     []folder.Folder{},
			expectError:  false,
		},
		{
			name:         "Folder does not exist in the given org",
			orgID:        org1,
			parentFolder: "foxtrot",
			expected:     nil,
			expectError:  true,
			errorMessage: "Folder does not exist in the given org",
		},
		{
			name:         "Invalid folder name",
			orgID:        org1,
			parentFolder: "invalid_folder",
			expected:     nil,
			expectError:  true,
			errorMessage: "Folder does not exist",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := folder.NewDriver(sampleFolders)

			got, err := f.GetAllChildFolders(tt.orgID, tt.parentFolder)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Error() != tt.errorMessage {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMessage, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(got) != len(tt.expected) {
				t.Errorf("Expected %d child folders, got %d", len(tt.expected), len(got))
			}

			for _, expectedFolder := range tt.expected {
				found := false
				for _, actualFolder := range got {
					if actualFolder.Name == expectedFolder.Name &&
						actualFolder.Paths == expectedFolder.Paths &&
						actualFolder.OrgId == expectedFolder.OrgId {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected child folder %+v not found", expectedFolder)
				}
			}
		})
	}
}
