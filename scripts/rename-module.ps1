param(
  [Parameter(Mandatory=$true)]
  [string]$NewMod,
  [string]$OldMod
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

if (-not $OldMod) {
  $line = (Get-Content go.mod | Select-String -Pattern '^module ').Line
  if (-not $line) {
    Write-Error "Cannot find 'module' line in go.mod"
  }
  $OldMod = ($line -split '\s+')[1]
}

Write-Host "OLD_MOD: $OldMod"
Write-Host "NEW_MOD: $NewMod"

# 1) go.mod
go mod edit -module $NewMod

# 2) 替换源码文本
$files = git ls-files | Select-String -Pattern '\.(go|md|yaml|yml|toml|json|mk|Makefile)$'
$files = $files.ToString().Split([Environment]::NewLine, [System.StringSplitOptions]::RemoveEmptyEntries)

foreach ($f in $files) {
  (Get-Content $f) -replace [regex]::Escape($OldMod), $NewMod | Set-Content $f
}

# 3) 整理依赖
go mod tidy

Write-Host "Module renamed successfully. Try: go build ./..."
