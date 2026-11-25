# New Relic CLI Guided Install Guide

## Overview

The `newrelic install` command provides an automated, guided installation experience for New Relic agents and integrations. This guide helps you understand when to use guided install and when to use alternative methods.

## Intended Use

### Guided Install is Designed For

The CLI-based guided installation provides a simple, single-command experience for **initial installation on new hosts**:

- Customers who want an easy, automated installation process
- Standard agent deployments with default configurations
- Embedding in automated deployment scripts to ensure new servers are monitored from day one
- Integration with configuration management tools (Ansible, Puppet, Chef)
- Scenarios where a quick, consistent installation is preferred

Advanced users who require fine-grained control over installation settings, specific version pinning, or custom configuration during installation may prefer manual installation procedures.

### How Guided Install Works

When you run `newrelic install`, the command:

1. Detects your system configuration and running services
2. Recommends appropriate New Relic agents and integrations
3. Installs agents with default configuration settings
4. Validates that agents are reporting data successfully

To ensure a consistent installation experience, the command configures agents with standard default settings. This approach works well for new installations where no prior configuration exists.

## Upgrading Existing Agents

If you have existing New Relic agents with custom configurations that you want to preserve, use the manual upgrade procedures instead:

### Infrastructure Agent Upgrades

For detailed upgrade instructions specific to your operating system, see the New Relic documentation:

**Documentation:** [Update Infrastructure Agent](https://docs.newrelic.com/docs/infrastructure/install-infrastructure-agent/update-or-uninstall/update-infrastructure-agent/)

### APM Agent Upgrades

Each APM agent has specific upgrade procedures that preserve your configurations. See the New Relic documentation for language-specific upgrade instructions:

**Documentation:** [APM Agent Updates](https://docs.newrelic.com/docs/apm/)

## Configuration Considerations

### What Happens to Existing Configurations

If you run `newrelic install` on a system with existing New Relic agents:

- Agent configuration files are set to default values
- Custom settings (log paths, custom attributes, labels) are not preserved
- Integration configurations are reset to defaults

This behavior ensures a clean, consistent installation but may not be suitable for systems with custom configurations you want to keep.

### Backing Up Configurations

If you plan to use guided install on a system with existing agents, consider backing up your configuration files first:

**Linux:**
```bash
# Back up main config file
sudo cp /etc/newrelic-infra.yml /etc/newrelic-infra.yml.backup

# Back up integrations directory
sudo cp -r /etc/newrelic-infra /etc/newrelic-infra.backup
```

**Windows:**
```powershell
Copy-Item -Path "C:\Program Files\New Relic\newrelic-infra" -Destination "C:\Program Files\New Relic\newrelic-infra.backup" -Recurse
```

## Frequently Asked Questions

### Can I use guided install to upgrade my existing agents?

While the command will run successfully, it's designed for new installations. For upgrades that preserve your custom configurations, we recommend using the manual upgrade procedures for your specific agent type.

### Will my telemetry data be affected?

No. Guided install only affects local agent configuration files. Your historical data in New Relic remains unchanged.

### Can I customize the installation?

The guided install uses standard default configurations. For custom setups, you can either:
1. Use guided install and then manually edit configuration files afterward
2. Perform a manual installation following agent-specific documentation

### What if I need help?

- **Documentation:** [New Relic Docs](https://docs.newrelic.com/)
- **Issues and Questions:** [GitHub Issues](https://github.com/newrelic/newrelic-cli/issues)
- **General Support:** [New Relic Support](https://support.newrelic.com/)

## Command Reference

For complete command options, run:

```bash
newrelic install --help
```

See also:
- [CLI Documentation](newrelic.md)
- [Getting Started Guide](GETTING_STARTED.md)
- [New Relic Infrastructure Agent Documentation](https://docs.newrelic.com/docs/infrastructure/)
