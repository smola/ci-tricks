package main

import (
	"fmt"
	"os"
)

func init() {
	Tricks = append(Tricks, InstallPostgres)
}

const PostgresVersionKey = "POSTGRESQL_VERSION"

func InstallPostgres(env *Env) error {
	version := os.Getenv(PostgresVersionKey)
	if version == "" {
		return nil
	}

	fmt.Println("Installing PostgreSQL", version)

	var err error
	switch env.Provider {
	case Appveyor:
		err = installPostgresAppveyor(env, version)
	case Travis:
		err = installPostgresTravis(env, version)
	default:
		err = fmt.Errorf("unsupported provider: %s", env.Provider)
	}

	if err == nil {
		err = WaitForTCP("127.0.0.1:5432")
	}

	if err != nil {
		return fmt.Errorf("error starting PostgreSQL: %s", err)
	}

	return nil
}

func installPostgresAppveyor(env *Env, version string) error {
	if env.OS != Windows {
		return fmt.Errorf("PostgreSQL is supported only on Windows")
	}

	if env.Arch != Amd64 {
		return fmt.Errorf("PostgreSQL is supported only on amd64")
	}

	return Run("net", "start", fmt.Sprintf("postgresql-x64-%s", version))
}

func installPostgresTravis(env *Env, version string) error {
	switch env.OS {
	case Linux:
		return installPostgresTravisLinux(env, version)
	case OSX:
		return installPostgresTravisOSX(env, version)
	default:
		return fmt.Errorf("PostgreSQL is only supported on Linux and macOS")
	}
}

func installPostgresTravisLinux(env *Env, version string) error {
	if err := Run("sudo", "service", "postgresql", "stop"); err != nil {
		return err
	}

	if err := Run("sudo", "service", "postgresql", "start", version); err != nil {
		return err
	}

	return nil
}

func installPostgresTravisOSX(env *Env, version string) error {
	const defaultVersion = "9.6"
	if err := Run("rm", "-rf", "/usr/local/var/postgres"); err != nil {
		return err
	}

	if version != defaultVersion {
		pkg := fmt.Sprintf("postgresql@%s", version)

		if err := Run("brew", "unlink", "postgresql"); err != nil {
			return err
		}

		if err := Run("brew", "install", pkg); err != nil {
			return err
		}

		if err := Run("brew", "link", "--force", pkg); err != nil {
			return err
		}
	}

	if _, err := os.Stat("/usr/local/var/postgres"); os.IsNotExist(err) {
		if err := Run("initdb", "/usr/local/var/postgres"); err != nil {
			return err
		}
	}

	if err := Run("pg_ctl", "-D", "/usr/local/var/postgres", "start"); err != nil {
		return err
	}

	if err := WaitForTCP("127.0.0.1:5432"); err != nil {
		return err
	}

	if err := Run("createuser", "-s", "-p", "5432", "postgres"); err != nil {
		return err
	}

	if err := Run("psql", "--version"); err != nil {
		return err
	}

	return nil
}
