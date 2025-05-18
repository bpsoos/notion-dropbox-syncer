package syncer

import (
	"fmt"
	"slices"
	"time"
)

func (s *Syncer) Sync(databaseId string, folderPath string) error {
	s.logger.Info("fetching rows from notion database")
	rows, err := s.notionAdapter.ListDatabaseRows(databaseId)
	if err != nil {
		return fmt.Errorf("list notion database rows error: %v", err)
	}

	s.logger.Info("fetching files from dropbox")
	files, err := s.dropboxAdapter.ListFolder(folderPath)
	if err != nil {
		return fmt.Errorf("list dropbox folder error: %v", err)
	}

	s.logger.Info("syncing")
	for _, file := range files {
		s.logger.Debug(file.Name)
		if slices.Contains(rows, file.Name) {
			s.logger.Debug("already present")
			continue
		}

		s.logger.Info(
			"file not yet in notion",
			"file_path", file.Path,
			"file_name", file.Name,
			"file_hash", file.Hash,
		)
		time.Sleep(time.Duration(s.syncDelayMs) * time.Millisecond)
		s.logger.Info("getting link")
		fileLink, err := s.dropboxAdapter.GetLink(file.Path)
		if err != nil {
			return fmt.Errorf("get dropbox link error: %v", err)
		}
		s.logger.Info("adding to notion")
		err = s.notionAdapter.AddDatabaseRow(databaseId, file.Name, fileLink)
		if err != nil {
			return fmt.Errorf("add notion database row error: %v", err)
		}
	}

	return nil
}
