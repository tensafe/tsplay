param(
    [string]$CodexHome = $env:CODEX_HOME,
    [string]$InstallDir,
    [string]$Tsplay,
    [string]$ExtractAssets,
    [switch]$Force
)

$ErrorActionPreference = "Stop"

function Write-Log {
    param([string]$Message)
    Write-Host $Message
}

function Fail {
    param([string]$Message)
    throw $Message
}

function Resolve-Dir {
    param([string]$Path)
    if (-not (Test-Path -LiteralPath $Path -PathType Container)) {
        Fail "Directory does not exist: $Path"
    }
    return [System.IO.Path]::GetFullPath((Resolve-Path -LiteralPath $Path).Path)
}

function Get-TsplayPath {
    if ($Tsplay) {
        if (-not (Test-Path -LiteralPath $Tsplay -PathType Leaf)) {
            Fail "tsplay binary was not found: $Tsplay"
        }
        return [System.IO.Path]::GetFullPath((Resolve-Path -LiteralPath $Tsplay).Path)
    }

    $command = Get-Command tsplay -ErrorAction SilentlyContinue
    if ($command) {
        return $command.Source
    }

    return $null
}

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$ScriptDir = Resolve-Dir $ScriptDir
$SkillName = Split-Path -Leaf $ScriptDir

if (-not (Test-Path -LiteralPath (Join-Path $ScriptDir "SKILL.md") -PathType Leaf)) {
    Fail "SKILL.md not found next to install.ps1"
}

if (-not $InstallDir) {
    if ($CodexHome) {
        $InstallDir = Join-Path $CodexHome "skills"
    } else {
        $InstallDir = Join-Path $HOME ".codex/skills"
    }
}

New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
$InstallDir = Resolve-Dir $InstallDir
$TargetDir = Join-Path $InstallDir $SkillName

if (Test-Path -LiteralPath $TargetDir -PathType Container) {
    $TargetReal = [System.IO.Path]::GetFullPath((Resolve-Path -LiteralPath $TargetDir).Path)
    if ($TargetReal.TrimEnd('\').ToLowerInvariant() -eq $ScriptDir.TrimEnd('\').ToLowerInvariant()) {
        Write-Log "Skill is already installed at $TargetDir"
    } else {
        if ($Force) {
            Remove-Item -LiteralPath $TargetDir -Recurse -Force
            Write-Log "Removed existing install at $TargetDir"
        } else {
            $BackupDir = "$TargetDir.backup.$((Get-Date).ToString('yyyyMMddHHmmss'))"
            Move-Item -LiteralPath $TargetDir -Destination $BackupDir
            Write-Log "Backed up existing install to $BackupDir"
        }
        New-Item -ItemType Directory -Force -Path $TargetDir | Out-Null
        Get-ChildItem -LiteralPath $ScriptDir -Force | ForEach-Object {
            Copy-Item -LiteralPath $_.FullName -Destination $TargetDir -Recurse -Force
        }
        Write-Log "Installed $SkillName to $TargetDir"
    }
} else {
    New-Item -ItemType Directory -Force -Path $TargetDir | Out-Null
    Get-ChildItem -LiteralPath $ScriptDir -Force | ForEach-Object {
        Copy-Item -LiteralPath $_.FullName -Destination $TargetDir -Recurse -Force
    }
    Write-Log "Installed $SkillName to $TargetDir"
}

if ($ExtractAssets) {
    $TsplayPath = Get-TsplayPath
    if (-not $TsplayPath) {
        Fail "tsplay was not found on PATH; rerun with -Tsplay C:\path\to\tsplay.exe or install tsplay first"
    }
    & $TsplayPath -action extract-assets -extract-root $ExtractAssets
    Write-Log "Extracted bundled assets to $ExtractAssets"
}

$DetectedTsplay = Get-TsplayPath
if ($DetectedTsplay) {
    Write-Log "Detected tsplay: $DetectedTsplay"
    Write-Log "Optional next step: $DetectedTsplay -action list-assets"
} else {
    Write-Log "tsplay was not found on PATH."
    Write-Log "Install a matching tsplay release binary, or use .\\tsplay.exe / go run . when working inside the repo."
}

Write-Log "Skill install complete."
