param(
    [string]$InstallDir = ".",
    [string]$BinaryName = "tsplay.exe",
    [switch]$SkipRun
)

$ErrorActionPreference = "Stop"

function Fail($Message) {
    throw $Message
}

function Resolve-ArchTag {
    try {
        $arch = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString().ToUpperInvariant()
    } catch {
        $arch = [string]$env:PROCESSOR_ARCHITECTURE
        if (-not $arch) {
            Fail "Unable to determine Windows architecture."
        }
        $arch = $arch.ToUpperInvariant()
    }

    switch ($arch) {
        "X64" { return "amd64" }
        "AMD64" { return "amd64" }
        "ARM64" {
            Write-Host "Windows ARM64 detected. Using the amd64 TSPlay build via x64 emulation."
            return "amd64"
        }
        default { Fail "Unsupported Windows architecture: $arch" }
    }
}

$repo = "tensafe/tsplay"
$baseUrl = "https://github.com/$repo/releases/latest/download"
$archTag = Resolve-ArchTag
$downloadUrl = "$baseUrl/tsplay-windows-$archTag.exe"

New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
$installDirFull = (Resolve-Path -Path $InstallDir).ProviderPath
$targetPath = Join-Path $installDirFull $BinaryName

Write-Host "Downloading TSPlay for windows/$archTag"
Invoke-WebRequest -Uri $downloadUrl -OutFile $targetPath

Write-Host "Installed: $targetPath"

if ($SkipRun) {
    Write-Host "Skipped quickstart run."
    exit 0
}

Write-Host "Running quickstart demo..."
& $targetPath -action quickstart-demo

Write-Host ""
Write-Host "Next steps:"
Write-Host "  $targetPath -action file-srv -addr :8000"
Write-Host "  $targetPath -flow script/tutorials/10_assert_page_state.flow.yaml"
