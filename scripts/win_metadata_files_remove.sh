# Removes:
#  1) the 'cmd/newrelic/versioninfo.json' created by the 'win_metadata.sh' script
#  2) the 'cmd/newrelic/resource_windows.syso' file created by 'go generate cmd/newrelic/main.go'

VERSION_FILE=cmd/newrelic/versioninfo.json
SYSO_FILE=cmd/newrelic/resource_windows.syso

if [ -f $VERSION_FILEn ]; then
  rm $VERSION_FILE
fi

if [ -f $SYSO_FILE ]; then
  rm $SYSO_FILE
fi
