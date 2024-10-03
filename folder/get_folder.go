package folder

import (
	"strings"
	"errors"

	"github.com/gofrs/uuid"
)

// errors
var (
	ErrFolderNotFound = errors.New("Folder does not exist")
	ErrFolderNotInOrg = errors.New("Folder does not exist in the given org")
)

func GetAllFolders() []Folder {
	return GetSampleData()
}



func (f *driver) GetFoldersByOrgID(orgID uuid.UUID) []Folder {
	folders := f.folders

	res := []Folder{}
	for _, f := range folders {
		if f.OrgId == orgID {
			res = append(res, f)
		}
	}

	return res

}



func (f *driver) GetAllChildFolders(orgID uuid.UUID, name string) ([]Folder, error) {
    var parentPath string
    childFolders := []Folder{}

	folderFound := false
	folderInOrg := false

    // Finding parent folder
    for _, folder := range f.folders { // similar to python for folder in self.folders
		if folder.Name == name {
			folderFound = true
			if folder.OrgId == orgID {
				parentPath = folder.Paths
				folderInOrg = true
				break
			}
		} 
    }

    // If parent not found
    if !folderFound {
        return nil, ErrFolderNotFound
    }
	if !folderInOrg {
		return nil, ErrFolderNotInOrg
	}


	// Finding child folders
	for _, folder := range f.folders {
		if folder.OrgId != orgID {
			continue // skip if not the same org
		}
		// if path starts with parent and isn't parent
		if strings.HasPrefix(folder.Paths, parentPath+".") && folder.Paths != parentPath {
			childFolders = append(childFolders, folder)
		}
	}

    return childFolders, nil
}
