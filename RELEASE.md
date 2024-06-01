# Releasing

This documents the release process of Goodwill. It's mostly a reminder for future Justen about how releases are done.

This work is typically done manually outside of GitHub infrastructure; but at some point it may be automated.

## Prerequisites

### Required Software

- [Maven 3.x](https://maven.apache.org/download.cgi)
- [Go 1.22+](https://go.dev/doc/install)
- [Java 17+](https://adoptium.net/temurin/releases/)
- [GPG](https://www.gnupg.org/download/index.html)
- [Signify](https://github.com/aperezdc/signify?tab=readme-ov-file#installation)
    - Homebrew: `brew install signify-osx`

### Configuration

- `~/.m2/settings.xml` - Required for all publishing to Sonatype Nexus/Maven Central
    ```xml
    <settings>
        <profiles>
            <profile>
                <id>maven-central-gpg</id>
                <properties>
                    <gpg.keyname>GPG_KEY_EMAIL_OR_ID</gpg.keyname>
                </properties>
            </profile>
        </profiles>
        <activeProfiles>
            <activeProfile>maven-central-gpg<activeProfile>
        </activeProfiles>
        <servers>
            <server>
                <id>ossrh</id>
                <username>OSSRH_TOKEN_USER</username>
                <password>OSSRH_TOKEN_PASS</password>
            </server>
        </servers>
    </settings>
    ```
- `SIGNIFY_KEY` environment variable - Required to sign checksums for GitHub Release
    ```sh
    export SIGNIFY_KEY=~/.signify/goodwill.sec
    ```
- `test/tofu/terraform.tfvars` - Required only to run `mage e2e:testPublished` a staged release before publishing.

    ```sh
    sonatype_username = "OSSRH_TOKEN_USER"
    sonatype_password = "OSSRH_TOKEN_PASS"
    sonatype_staging_repo = "techjustenconcord-00000"
    ```

## Pre-Release

### Publish Snapshot to Sonatype Nexus / Maven Central

A snapshot version is used for publishing a test release that can be consumed before committing to the actual released version. 
To start a snapshot, use the `mage snapshot` target:

```sh
VERSION=0.7.0 mage snapshot
```

This changes the `pom.xml` by writing a snapshot version

```xml
    <groupId>tech.justen.concord</groupId>
    <artifactId>goodwill</artifactId>
    <version>0.7.0-SNAPSHOT</version>
```

Afterwards, you can publish an artifact directly to Nexus using the `nexus:deploy` target.

To E2E Test the SNAPSHOT artifact, you can run `mage e2e:testPublished`, which will use the `pom.xml` version for testing.

### Publish Pre-Release to GitHub

## Release

### Publish Release to Sonatype Nexus / Maven Central

A release can be pushed first to a staging repository before it is committed to the Maven Central repository.
Once it is deployed, it cannot be removed; so use the staging repository to test a release before publishing. 

First, create a release version using `mage release`
```sh
VERSION=0.7.0 mage release
```

This changes the `pom.xml` by writing the new version

```xml
    <groupId>tech.justen.concord</groupId>
    <artifactId>goodwill</artifactId>
    <version>0.7.0</version>
```

Afterwards, you can publish an artifact to Nexus staging using the `nexus:deploy` target.

To E2E Test the staged artifact, you can run `mage e2e:testPublished`, which will use the `pom.xml` version for testing.

Once testing is completed to satisfaction, release the artifact from staging using `mage nexus:release`.
To drop the staged artifact, run `mage nexus:drop`.

### Publishing Release to GitHub

Run `mage sign` to package and sign SHA256 sums for the artifacts.

The `mage release` command should have made a tag in the repository.
Push this tag using: `git push --tags`.

Create a Release based off this tag in GitHub.

Upload the artifacts in the `dist` folder as assets to the release.
