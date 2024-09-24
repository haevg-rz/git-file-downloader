package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/dhcgn/GitLabFileDownloader/internal"
	"github.com/dhcgn/GitLabFileDownloader/internal/api"
)

// AppName is the name of the application
const AppName = "GitLab File Downloader"

var (
	version  = "undef"
	commitID = "undef"

	flagTokenPtr = flag.String(internal.FlagNameToken, ``, `Private-Token with access right for "api" and "read_repository, role must be minimum "Reporter""`)

	flagOutPathPtr      = flag.String(internal.FlagNameOutPath, ``, "Path to write file to disk")
	flagRepoFilePathPar = flag.String(internal.FlagNameRepoFilePath, ``, "File path in repo, like src/main.go")

	flagBranchPtr = flag.String(internal.FlagNameBranch, `main`, "Branch")

	flagOutFolderPtr      = flag.String(internal.FlagNameOutFolder, ``, "Folder to write file to disk")
	flagRepoFolderPathPtr = flag.String(internal.FlagNameRepoFolderPathEscaped, ``, "Folder to write file to disk")

	flagUrlPtr           = flag.String(internal.FlagNameUrl, ``, "Url to Api v4, like https://my-git-lab-server.local/api/v4/")
	flagProjectNumberPtr = flag.Int(internal.FlagNameProjectNumber, 0, "The Project ID from your project")

	flagIncludeOnlyPtr = flag.String(internal.IncludeOnly, ``, "Include only these regex pattern")
	flagExcludePtr     = flag.String(internal.Exclude, ``, "Exclude these regex pattern")
)

func main() {
	mainSub()
}

func mainSub() {
	log.Println(AppName, "Version:", version, "Commit:", commitID)
	log.Println(`Project: https://github.com/dhcgn/GitLabFileDownloader/`)

	flag.Parse()

	settings := getSettingsFromFlags()
	isValid, args, msgs := settings.IsValid()
	if !isValid {
		log.Println("Arguments are missing:", args)
		log.Println("Messages:", msgs)
		flag.PrintDefaults()
		return
	}

	branches, err := api.GetBranches(settings)
	if err != nil || len(branches) == 0 {
		log.Println("Error GetBranches:", err)
		return
	}

	found := false
	for _, branch := range branches {
		if branch.Name == settings.Branch {
			found = true
		}
	}

	if !found {
		log.Println("Branch not found:", settings.Branch)
		log.Println("Available branches:", branches)
		return
	}

	switch settings.Mode() {
	case internal.ModeFile:
		log.Println("Mode: File")
		fileModeHandling(settings)
	case internal.ModeFolder:
		log.Println("Mode: Folder")
		folderModeHandling(settings)
	}
}

// exists returns whether the given file or directory exists
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func folderModeHandling(settings internal.Settings) {
	if !exists(settings.OutFolder) {
		err := os.Mkdir(settings.OutFolder, 0755)
		if err != nil {
			log.Println("Error:", err)
			return
		}
	}

	files, err := api.GetFilesFromFolder(settings)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	log.Println("Sync", len(files), "files, from remote folder", settings.RepoFolderPath)

	for _, file := range files {

		if settings.IncludeOnly != "" {
			m, err := regexp.MatchString(settings.IncludeOnly, file.Name)
			if err == nil {
				if !m {
					log.Println("Skip:", file.Name, "because include only rule:", settings.IncludeOnly)
					continue
				}
			}
		}
		if settings.Exclude != "" {
			m, err := regexp.MatchString(settings.Exclude, file.Name)
			if err == nil {
				if m {
					log.Println("Skip:", file.Name, "because exclude rule:", settings.Exclude)
					continue
				}
			}
		}

		if file.Type == "tree" {
			folderSettings := settings
			outFolder := path.Join(folderSettings.OutFolder, file.Name)
			folderSettings.OutFolder = outFolder
			folderSettings.RepoFolderPath = file.Path
			folderModeHandling(folderSettings)

			continue
		}

		fileSettings := settings
		outFile := path.Join(fileSettings.OutFolder, path.Base(file.Path))
		fileSettings.OutFile = outFile
		fileSettings.RepoFilePath = file.Path

		fileModeHandling(fileSettings)
	}
}

func fileModeHandling(settings internal.Settings) {
	new, err := fileModeHandlingInternal(settings)

	if err != nil {
		log.Println("Error at", settings.RepoFilePath, ":", err)
		return
	}
	if new {
		log.Println("Wrote file:", settings.RepoFilePath, ", because is new or changed")
		return
	}
	log.Println("Skip:", settings.RepoFilePath, ", because content is equal")
}

func fileModeHandlingInternal(settings internal.Settings) (bool, error) {
	exists, dir := testTargetFolder(settings.OutFile)
	if !exists {
		return false, fmt.Errorf("Target folder %v doesn't exists", dir)
	}

	gitLapFile, err := api.GetFile(settings)
	if err != nil {
		return false, fmt.Errorf("API Call error: %v", err)
	}

	fileData, err := base64.StdEncoding.DecodeString(gitLapFile.Content)
	if err != nil {
		return false, fmt.Errorf("DecodeString: %v", err)
	}

	isEqual, err := isOldFileEqual(gitLapFile, settings)
	if err != nil {
		return false, fmt.Errorf("isOldFileEqual: %v", err)
	}

	if isEqual {
		return false, nil
	}

	err = os.WriteFile(settings.OutFile, fileData, 0644)
	if err != nil {
		return false, fmt.Errorf("WriteFile: %v", err)
	}
	return true, nil
}

func isOldFileEqual(gitLapFile api.GitLapFile, settings internal.Settings) (bool, error) {
	if _, err := os.Stat(settings.OutFile); err == nil {
		file, err := os.Open(settings.OutFile)
		if err != nil {
			return false, err
		}
		defer file.Close()
		hash := sha256.New()
		if _, err := io.Copy(hash, file); err != nil {
			return false, err
		}
		sha256Hex := hex.EncodeToString(hash.Sum(nil))

		return sha256Hex == gitLapFile.ContentSha256, nil
	}
	return false, nil
}

func testTargetFolder(outFile string) (bool, string) {
	dir := filepath.Dir(outFile)
	if _, err := os.Stat(dir); err == nil {
		return true, dir
	} else if os.IsNotExist(err) {
		return false, dir
	}
	return false, dir
}

func getSettingsFromFlags() internal.Settings {
	return internal.Settings{
		PrivateToken:   *flagTokenPtr,
		OutFile:        *flagOutPathPtr,
		OutFolder:      *flagOutFolderPtr,
		Branch:         *flagBranchPtr,
		ApiUrl:         *flagUrlPtr,
		ProjectNumber:  strconv.Itoa(*flagProjectNumberPtr),
		RepoFilePath:   *flagRepoFilePathPar,
		RepoFolderPath: *flagRepoFolderPathPtr,
		UserAgent:      AppName + " " + version,
		IncludeOnly:    *flagIncludeOnlyPtr,
		Exclude:        *flagExcludePtr,
	}
}
