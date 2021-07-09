// +build integration

package recipes

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/config"
)

func TestLocalRecipeFetcher_FetchRecipes_EmptyManifest(t *testing.T) {
	tmp, err := ioutil.TempDir("/tmp", "newrelic")
	defer os.RemoveAll(tmp)
	require.NoError(t, err)

	config.Init(tmp)

	path := filepath.Join(tmp, "recipes")
	err = os.MkdirAll(path, 0750)
	require.NoError(t, err)
	err = ioutil.WriteFile(filepath.Join(path, "infra.yml"), []byte(sampleRecipe), 0600)
	require.NoError(t, err)

	f := LocalRecipeFetcher{
		Path: path,
	}
	recipes, err := f.FetchRecipes(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, recipes)
}

func TestLocalRecipeFetcher_FetchRecipes(t *testing.T) {
	tmp, err := ioutil.TempDir("/tmp", "newrelic")
	defer os.RemoveAll(tmp)
	require.NoError(t, err)

	config.Init(tmp)

	path := filepath.Join(tmp, "recipes")

	err = os.MkdirAll(path, 0750)
	require.NoError(t, err)

	err = ioutil.WriteFile(filepath.Join(path, "infra.yml"), []byte(sampleRecipe), 0600)
	require.NoError(t, err)

	f := LocalRecipeFetcher{
		Path: path,
	}
	recipes, err := f.FetchRecipes(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, recipes)

}

var sampleRecipe = `---
name: infrastructure-agent-installer
displayName: Infra

installTargets:
  - type: host
    os: linux
    platform: "debian"

keywords: []
processMatch: []

validationNrql: "SELECT count(*) from SystemSample where hostname like '{{.HOSTNAME}}' FACET entityGuid SINCE 10 minutes ago"

install:
  version: "3"
  silent: true

  tasks:
    default:
      cmds:
        - task: assert_pre_req
        - task: setup_license
        - task: update_apt
        - task: add_gpg_key
        - task: add_nr_source
        - task: update_apt_nr_source
        - task: install_infra
        - task: restart
        - task: assert_agent_started

    assert_pre_req:
      cmds:
        - |
          IS_GREP_INSTALLED=$(which grep | wc -l)
          if [ $IS_GREP_INSTALLED -eq 0 ] ; then
            echo "grep is required to run the newrelic install. Please install grep and re-run the installation." >> /dev/stderr
            exit 10
          fi
        - |
          IS_SED_INSTALLED=$(which sed | wc -l)
          if [ $IS_SED_INSTALLED -eq 0 ] ; then
            echo "sed is required to run the newrelic install. Please install sed and re-run the installation." >> /dev/stderr
            exit 11
          fi
        - |
          IS_AWK_INSTALLED=$(which awk | wc -l)
          if [ $IS_AWK_INSTALLED -eq 0 ] ; then
            echo "awk is required to run the newrelic install. Please install awk and re-run the installation." >> /dev/stderr
            exit 12
          fi
        - |
          IS_CAT_INSTALLED=$(which cat | wc -l)
          if [ $IS_CAT_INSTALLED -eq 0 ] ; then
            echo "cat is required to run the newrelic install. Please install cat and re-run the installation." >> /dev/stderr
            exit 13
          fi
        - |
          IS_TEE_INSTALLED=$(which tee | wc -l)
          if [ $IS_TEE_INSTALLED -eq 0 ] ; then
            echo "tee is required to run the newrelic install. Please install tee and re-run the installation." >> /dev/stderr
            exit 14
          fi
        - |
          IS_TOUCH_INSTALLED=$(which touch | wc -l)
          if [ $IS_TOUCH_INSTALLED -eq 0 ] ; then
            echo "touch is required to run the newrelic install. Please install touch and re-run the installation." >> /dev/stderr
            exit 15
          fi
        - |
          IS_DOCKER_CONTAINER=$(sudo grep 'docker\|lxc' /proc/1/cgroup | wc -l)
          if [ $IS_DOCKER_CONTAINER -gt 0 ] ; then
            echo "We’ve detected that you are installing our infrastructure agent inside a docker container. This agent is not designed to be installed within a container, but rather on the host it is running on or as a separate container. For additional information visit: https://docs.newrelic.com/docs/infrastructure/install-infrastructure-agent/linux-installation/docker-container-infrastructure-monitoring/" >> /dev/stderr
            exit 18
          fi
        - |
          IS_WSL_CONTAINER=$(sudo grep -i 'Microsoft' /proc/version | wc -l)
          if [ $IS_WSL_CONTAINER -gt 0 ] ; then
            echo "Sorry, our infrastructure agent cannot be installed for Microsoft Windows Subsystem for Linux, this is an unsupported operating system." >> /dev/stderr
            exit 19
          fi
        - |
          if [ -n "{{.DEBIAN_CODENAME}}" ]; then
            IS_INFRA_AVAILABLE=$(curl -Is https://download.newrelic.com/infrastructure_agent/linux/apt/dists/{{.DEBIAN_CODENAME}}/InRelease | grep " 2[0-9][0-9] " | wc -l)
            if [ $IS_INFRA_AVAILABLE -eq 0 ] ; then
              echo "there is no newrelic infrastructure agent available for the distribution with codename '{{.DEBIAN_CODENAME}}'." >> /dev/stderr
              exit 21
            fi
          else
            if [ -n "{{.DEBIAN_VERSION_CODENAME}}" ]; then
              IS_INFRA_AVAILABLE=$(curl -Is https://download.newrelic.com/infrastructure_agent/linux/apt/dists/{{.DEBIAN_VERSION_CODENAME}}/InRelease | grep " 2[0-9][0-9] " | wc -l)
              if [ $IS_INFRA_AVAILABLE -eq 0 ] ; then
                echo "there is no newrelic infrastructure agent available for the distribution with version codename '{{.DEBIAN_VERSION_CODENAME}}'." >> /dev/stderr
                exit 21
              fi
            else
              echo "there is no newrelic infrastructure agent available for the distribution, no version codename was found." >> /dev/stderr
              exit 21
            fi
          fi
      vars:
        DEBIAN_CODENAME:
          sh: awk -F= '/VERSION_CODENAME/ {print $2}' /etc/os-release
        DEBIAN_VERSION_CODENAME:
          sh: cat /etc/os-release | grep "VERSION=\"[0-9] " | awk -F " " '{print $2}' | sed 's/[()"]//g'

    setup_license:
      cmds:
        - |
          if [ -f /etc/newrelic-infra.yml ]; then
            sudo rm /etc/newrelic-infra.yml;
          fi
          sudo touch /etc/newrelic-infra.yml;
        - |
          sudo tee /etc/newrelic-infra.yml > /dev/null <<"EOT"
          license_key: {{.NEW_RELIC_LICENSE_KEY}}
          enable_process_metrics: true
          EOT
        - |
          if [ $(echo {{.NEW_RELIC_REGION}} | grep -i staging | wc -l) -gt 0 ]; then
            echo 'staging: true' | sudo tee -a /etc/newrelic-infra.yml > /dev/null
          fi

    update_apt:
      cmds:
        - |
          sudo apt-get update -yq
      silent: true

    add_gpg_key:
      cmds:
        - |
          sudo apt-get install gnupg2 -y
          sudo curl -s https://download.newrelic.com/infrastructure_agent/gpg/newrelic-infra.gpg | sudo apt-key add -
      silent: true

    add_nr_source:
      cmds:
        - |
          if [ -n "{{.DEBIAN_CODENAME}}" ]; then
            printf "deb [arch=amd64] https://download.newrelic.com/infrastructure_agent/linux/apt {{.DEBIAN_CODENAME}} main" | sudo tee /etc/apt/sources.list.d/newrelic-infra.list > /dev/null
          else
            printf "deb [arch=amd64] https://download.newrelic.com/infrastructure_agent/linux/apt {{.DEBIAN_VERSION_CODENAME}} main" | sudo tee /etc/apt/sources.list.d/newrelic-infra.list > /dev/null
          fi
      vars:
        DEBIAN_CODENAME:
          sh: awk -F= '/VERSION_CODENAME/ {print $2}' /etc/os-release
        DEBIAN_VERSION_CODENAME:
          sh: cat /etc/os-release | grep "VERSION=\"[0-9] " | awk -F " " '{print $2}' | sed 's/[()"]//g'
      silent: true

    update_apt_nr_source:
      cmds:
        - |
          sudo apt-get update -yq

    install_infra:
      cmds:
        - |
          sudo apt-get install newrelic-infra -y -qq
      silent: true

    restart:
      cmds:
        - |
          if [ {{.IS_SYSTEMCTL}} -gt 0 ]; then
            sudo systemctl restart newrelic-infra
          else 
            if [ {{.IS_INITCTL}} -gt 0 ]; then
              sudo initctl restart newrelic-infra
            else
              sudo /etc/init.d/newrelic-infra restart
            fi
          fi
      vars:
        IS_SYSTEMCTL:
          sh: which systemctl | wc -l
        IS_INITCTL:
          sh: which initctl | wc -l

    assert_agent_started:
      cmds:
        - |
          # Ensure agent has enough time to start
          sleep 10
          IS_INFRA_INSTALLED=$(sudo ps aux | grep newrelic-infra-service | grep -v grep | wc -l)
          if [ $IS_INFRA_INSTALLED -eq 0 ] ; then
            echo "The infrastructure agent has not started after installing. Please try again later, or see our documentation for installing manually https://docs.newrelic.com/docs/using-new-relic/cross-product-functions/install-configure/install-new-relic" >> /dev/stderr
            exit 31
          fi

postInstall:
  info: |2
      ⚙️  The Infrastructure Agent configuration file can be found in /etc/newrelic-infra.yml
      Edit this file to make changes or configure advanced features for the agent. See the docs for options:
      https://docs.newrelic.com/docs/infrastructure/install-infrastructure-agent/configuration/infrastructure-agent-configuration-settings
      
      Note: Process monitoring has been enabled by default - all other config options are left to the user.
`
