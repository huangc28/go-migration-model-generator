package genmodel

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

const (
	DefaultSourcePath   = "db/migrations"
	DefaultDestFilename = "db/schema.sql"
	DefaultHost         = "127.0.0.1"
	DefaultPort         = 5432
	DefaultUsername     = "postgres"
	DefaultPassword     = ""
	DefaultDBName       = ""
)

var (
	Source     string
	PGHost     string
	PGPort     int
	PGUsername string
	PGPassword string
	PGDBname   string
)

func init() {
	genModelCmd.Flags().StringVarP(&Source, "source", "s", DefaultSourcePath, "Migration directory to read from")
	genModelCmd.Flags().StringVarP(&PGHost, "host", "", DefaultHost, "PG host")
	genModelCmd.Flags().IntVarP(&PGPort, "port", "", DefaultPort, "PG Host port")
	genModelCmd.Flags().StringVarP(&PGUsername, "username", "u", DefaultUsername, "PG username")
	genModelCmd.Flags().StringVarP(&PGPassword, "password", "", DefaultPassword, "PG password")
	genModelCmd.Flags().StringVarP(&PGDBname, "dbname", "", PGDBname, "PG database name")

}

func combineProjectPath(source string) string {
	cwd, _ := os.Getwd()

	return filepath.Join(cwd, source)
}

func GetMigrationInfo(db *sql.DB) (int, bool, error) {
	var (
		version int
		dirty   bool
	)

	err := db.QueryRow(`
		SELECT version, dirty
		FROM schema_migrations
	`).Scan(&version, &dirty)

	if err != nil {
		return 0, false, err
	}

	return version, dirty, nil
}

// pick list of file names that are <= specific version.
func pickMigrationsByVersion(files []os.FileInfo, version int) []os.FileInfo {
	var suitedFiles []os.FileInfo

	for _, file := range files {
		migVer, err := strconv.Atoi(strings.Split(file.Name(), "_")[0])

		if err != nil {
			log.Printf("failed to parse version number of %s, skipping...", file.Name())
			continue
		}

		sufSegs := strings.Split(file.Name(), ".")
		migType := sufSegs[len(sufSegs)-2 : len(sufSegs)-1][0]

		if migVer <= version && migType == "up" {
			suitedFiles = append(suitedFiles, file)
		}
	}

	return suitedFiles
}

func appendFileContentToDestFile(files []os.FileInfo, src string, dest string) {
	destFile, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	defer destFile.Close()

	if err != nil {
		log.Fatalf("failed to open / create dest file %s", err.Error())
	}

	for _, file := range files {
		func(file os.FileInfo) {
			cByte, err := ioutil.ReadFile(filepath.Join(src, file.Name()))

			if err != nil {
				log.Fatalf("Failed to read bytes from src file %s", err.Error())
			}

			if _, err := destFile.Write(cByte); err != nil {
				log.Fatalf("Failed when piping content from %s, exiting... with error %s", filepath.Join(src, file.Name()), err.Error())
			}
		}(file)
	}
}

func Gen(cmd *cobra.Command, args []string) error {
	// read latest migration info from database
	psqlDSN := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		PGHost,
		PGPort,
		PGUsername,
		PGPassword,
		PGDBname,
	)

	db, err := sql.Open("postgres", psqlDSN)

	if err != nil {
		log.Fatalf("Failed to connect to psql %v", err.Error())
		return err
	}

	// ---------- read migration info of the project  ----------
	// read version info from schema_migrations table
	version, dirty, err := GetMigrationInfo(db)

	if err != nil {
		log.Fatalf("Failed to read migration info %s", err.Error())

		return err
	}

	log.Printf("\n version: %d\n dirty: %t \n", version, dirty)

	if dirty {
		log.Fatal("migration seems dirty! Please fix the migration first")

		return err
	}

	// ---------- read `migration up` files from migration directory ----------
	sourceAbsolutePath := combineProjectPath(Source)

	log.Printf("reading from migration source path...  %s", sourceAbsolutePath)

	files, err := ioutil.ReadDir(sourceAbsolutePath)

	if err != nil {
		log.Fatalf("failed to read migrations from path %s, %s", sourceAbsolutePath, err.Error())

		return err
	}

	mFiles := pickMigrationsByVersion(files, version)

	appendFileContentToDestFile(
		mFiles,
		Source,
		combineProjectPath(DefaultDestFilename),
	)

	// ---------- execute sqlc generate command ----------
	osCmd := exec.Command("sqlc", "generate")

	osCmd.Env = os.Environ()
	osCmd.Stdout = os.Stdout
	osCmd.Stderr = os.Stderr

	if err := osCmd.Start(); err != nil {
		log.Fatalf("error executing command: 'sqlc generate' %s", err.Error())

		return err
	}

	if err := osCmd.Wait(); err != nil {
		log.Fatalf("error waiting for command: 'sqlc generate' %s", err.Error())

		return err
	}

	return nil
}

var genModelCmd = &cobra.Command{
	Use:   "gen",
	Short: "read / collect SQL from list of migrations files to genderate models in go code.",
	RunE:  Gen,
}
