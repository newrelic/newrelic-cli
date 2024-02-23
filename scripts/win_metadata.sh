# Prepares the 'versioninfo.json' file inside the 'cmd/newrelic/' folder so that the
#'go generate cmd/newrelic/main.go' command in the goreleaser 'before' hook grabs it
# to create the 'resource.syso' file which is needed to embed versioning information
# into the resulting Windows OS 'newrelic.exe' binary.

VERSION_FILE=cmd/newrelic/versioninfo.json
SYSO_FILE=cmd/newrelic/resource.syso
TPL_FILE=templates/versioning/versioninfo.json.template
YEAR=$(date +%Y)

if [ -f $VERSION_FILE ]; then
  rm $VERSION_FILE
fi

if [ -f $SYSO_FILE ]; then
  rm $SYSO_FILE
fi

cp $TPL_FILE $VERSION_FILE

sed -i "s/{CLIMajorVersion}/$1/g" $VERSION_FILE
sed -i "s/{CLIMinorVersion}/$2/g" $VERSION_FILE
sed -i "s/{CLIPatchVersion}/$3/g" $VERSION_FILE
sed -i "s/{Year}/$YEAR/g" $VERSION_FILE
