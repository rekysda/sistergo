<#
PowerShell script to install Apache virtual host in Laragon and add hosts entry.
Usage (as Administrator):
  Right click PowerShell -> Run as Administrator
  cd C:\laragon\www\sistergo\scripts
  .\setup-laragon.ps1
#>

param(
    [string]$LaragonRoot = 'C:\laragon',
    [string]$Domain = 'sistergo.test',
    [string]$ConfName = 'sistergo.conf',
    [int]$GoPort = 8080
)

function Ensure-Admin {
    $isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
    if (-not $isAdmin) {
        Write-Error "Script must be run as Administrator. Right-click PowerShell -> Run as Administrator."
        exit 1
    }
}

function Write-HostsEntry($domain) {
    $hostsPath = "$env:SystemRoot\System32\drivers\etc\hosts"
    $entry = "127.0.0.1 `t$domain"
    $hosts = Get-Content $hostsPath -ErrorAction Stop
    if ($hosts -notcontains $entry -and ($hosts -notcontains ("127.0.0.1 `t$domain"))) {
        Add-Content -Path $hostsPath -Value $entry
        Write-Output "Added hosts entry: $entry"
    } else {
        Write-Output "Hosts entry already present: $domain"
    }
}

function Install-Conf($laragonRoot, $domain, $confName, $port) {
    $sitesEnabledPath = Join-Path $laragonRoot "etc\apache2\sites-enabled"
    if (-not (Test-Path $sitesEnabledPath)) {
        Write-Error "Cannot find Laragon Apache sites-enabled folder: $sitesEnabledPath. Is Laragon installed at $laragonRoot?"
        exit 1
    }

    # Resolve template relative to this script's folder: project-root/laragon/sistergo.conf.template
    $scriptRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
    $candidate = Join-Path $scriptRoot "..\laragon\sistergo.conf.template"
    try {
        $templatePath = (Resolve-Path $candidate -ErrorAction Stop).Path
    } catch {
        Write-Error "Template file not found at $candidate. Please ensure the template exists at project-root\laragon\sistergo.conf.template or pass a valid laragon template path.";
        return
    }
    if (-not (Test-Path $templatePath)) {
        Write-Error "Template file not found: $templatePath"
        exit 1
    }

    $targetPath = Join-Path $sitesEnabledPath $confName
    if (-not (Test-Path $templatePath)) {
        Write-Error "Template file missing: $templatePath. Aborting vhost install."
        return
    }
    $confText = Get-Content $templatePath -Raw
    $confText = $confText -replace "sistergo.test", $domain
    $confText = $confText -replace "127.0.0.1:8080", "127.0.0.1:$port"

    Set-Content -Path $targetPath -Value $confText -Force
    Write-Output "Installed vhost config to $targetPath"
}

function Enable-Apache-ProxyModules($laragonRoot) {
    # Look for httpd.conf in Laragon etc dir
    $httpdFiles = Get-ChildItem -Path $laragonRoot -Recurse -Filter httpd.conf -ErrorAction SilentlyContinue | Select-Object -First 1
    if (-not $httpdFiles) { Write-Output "httpd.conf not found automatically; ensure mod_proxy and mod_proxy_http are enabled in your Apache config."; return }

    $httpdPath = $httpdFiles.FullName
    Write-Output "Found httpd.conf at $httpdPath"

    $httpdContent = Get-Content $httpdPath -Raw
    $httpdContent = $httpdContent -replace "#(LoadModule\s+proxy_module\s+modules\\mod_proxy.so)", '$1'
    $httpdContent = $httpdContent -replace "#(LoadModule\s+proxy_http_module\s+modules\\mod_proxy_http.so)", '$1'

    Set-Content -Path $httpdPath -Value $httpdContent
    Write-Output "Attempted to enable mod_proxy and mod_proxy_http in httpd.conf (search/replace). Please verify the httpd.conf manually if needed."
}

# Main
Ensure-Admin
Write-Output "Installing Laragon vhost for domain $Domain"
Write-HostsEntry -domain $Domain
Install-Conf -laragonRoot $LaragonRoot -domain $Domain -confName $ConfName -port $GoPort
Enable-Apache-ProxyModules -laragonRoot $LaragonRoot

Write-Output "Done. Restart Apache via Laragon UI (Menu -> Stop -> Start) to apply changes."
Write-Output "If your Go server runs on port $GoPort, you can visit http://$Domain/api/health to verify."