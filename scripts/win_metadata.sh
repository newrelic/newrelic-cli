# Prepares the 'versioninfo.json' file inside the 'cmd/newrelic/' folder so that the
#'go generate cmd/newrelic/main.go' command in the goreleaser 'before' hook grabs it
# to create the 'resource_windows.syso' file which is needed to embed versioning information
# into the resulting Windows OS 'newrelic.exe' binary.

VERSION_FILE=cmd/newrelic/versioninfo.json
SYSO_FILE=cmd/newrelic/resource_windows.syso
TPL_FILE=templates/versioning/versioninfo.json.template

VERSION=$(git describe --tags --abbrev=0 | sed 's/v//g')
echo "VERSION: $VERSION"

# Get semver values.
MAJOR=$(echo $VERSION | cut -d'.' -f1)
MINOR=$(echo $VERSION | cut -d'.' -f2)
PATCH=$(echo $VERSION | cut -d'.' -f3)

YEAR=$(date +%Y)

if [ -f $VERSION_FILE ]; then
  rm $VERSION_FILE
fi

if [ -f $SYSO_FILE ]; then
  rm $SYSO_FILE
fi

cp $TPL_FILE $VERSION_FILE

sed -i "s/{CLIMajorVersion}/$MAJOR/g" $VERSION_FILE
sed -i "s/{CLIMinorVersion}/$MINOR/g" $VERSION_FILE
sed -i "s/{CLIPatchVersion}/$PATCH/g" $VERSION_FILE
sed -i "s/{Year}/$YEAR/g" $VERSION_FILE
