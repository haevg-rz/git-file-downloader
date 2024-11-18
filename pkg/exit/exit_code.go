package exit

const (
	/*
		Success:
		Everything went according to plan
	*/
	Success = iota

	/*
		InternalError:
		Unexpected internal error. Please contact maintainer for help.
	*/
	InternalError

	/*
		ReceivedSignal:
		Signifies that the program received a signal while operating. for example: "^C"
	*/
	ReceivedSignal

	/*
		MissingFlags:
		One or more mandatory flags are missing
	*/
	MissingFlags

	/*
		InvalidOutPath:
		Given out-(file/folder) path is invalid
	*/
	InvalidOutPath

	/*
		UnknownGitProvider
		Specified git provider is unknown. Valid options are: GITHUB, GITLAB, AZURE
	*/
	UnknownGitProvider

	/*
		BranchOrRepoNotFound:
		Could not find remote branch
	*/
	BranchOrRepoNotFound

	/*
		FailedToCreateFile:
		Failed to create local file in the case it didn't exist before
	*/
	FailedToCreateFile

	/*
		FailedToOpenFile:
		Failed to open local file when trying to compare hashes
	*/
	FailedToOpenFile

	/*
		FailedToWriteToFile:
		Failed to write remote file data to local file
	*/
	FailedToWriteToFile

	/*
		FailedToCreateFolder:
		Failed to create non-existing folder.
		Happens if you provide a path with non-existing folders that are not the last element of the path
		example: foo/bar where foo doesn't exist.
	*/
	FailedToCreateFolder

	/*
		FailedToDecodeRemoteFileContent:
		Failed to decode base64 file content from remote file
	*/
	FailedToDecodeRemoteFileContent

	/*
		FailedToRetrieveRemoteFile:
		Failed to get file from remote repository. Can happen due to faulty remote path.
	*/
	FailedToRetrieveRemoteFile

	/*
		FailedToGetFilesFromRemoteFolder:
		Failed to get files from remote repository folder
	*/
	FailedToGetFilesFromRemoteFolder
)

var (
	Code = Success
)
