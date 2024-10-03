package folder

import (
	"errors"
	"strings"

)

var (
	ErrSourceFolderNotFound = errors.New("source folder does not exist")
	ErrDestinationFolderNotFound = errors.New("destination folder does not exist")
	ErrCannotMoveToChild = errors.New("cannot move a folder to a child of itself")
	ErrCannotMoveToSelf = errors.New("cannot move a folder to itself")
	ErrCannotMoveBetweenOrgs = errors.New("cannot move a folder to a different organization")
)

func (f *driver) MoveFolder(name string, dst string) ([]Folder, error) {
	var sourceFolder *Folder
	var destinationFolder *Folder

	// Find source folder
	for i, folder := range f.folders {
		if folder.Name == name {
			sourceFolder = &f.folders[i]
			break
		}
	}

	if sourceFolder == nil {
		return nil, ErrSourceFolderNotFound
	}

	// Find destination folder
	for i, folder := range f.folders {
		if folder.Name == dst {
			destinationFolder = &f.folders[i]
			break
		}
	}

	if destinationFolder == nil {
		return nil, ErrDestinationFolderNotFound
	}

	// if source and destination are the same
	if sourceFolder.Paths == destinationFolder.Paths {
		return nil, ErrCannotMoveToSelf
	}
	
	// if destination is a child of source
	if strings.HasPrefix(destinationFolder.Paths, sourceFolder.Paths+".") {
		return nil, ErrCannotMoveToChild
	}

	// if source and destination are in different orgs
	if sourceFolder.OrgId != destinationFolder.OrgId {
		return nil, ErrCannotMoveBetweenOrgs
	}

	// update paths
	oldParentPath := sourceFolder.Paths
	newParentPath := destinationFolder.Paths + "." + sourceFolder.Name

	sourceFolder.Paths = newParentPath

	// update paths of all children
	for i := range f.folders {
		folder := &f.folders[i]
		if folder.OrgId == sourceFolder.OrgId && folder.Paths != oldParentPath {
			if strings.HasPrefix(folder.Paths, oldParentPath+".") {
				suffix := strings.TrimPrefix(folder.Paths, oldParentPath)
				folder.Paths = newParentPath + suffix
			}
		}
	}

	return f.folders, nil
}
