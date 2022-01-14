$Addr="{{GetProperty "init-server-address"}}"
$DeployYaml="{{.InitApplyYaml}}"
$DbPath='{{.DBPath}}'
$Token="{{.Token}}"
$Mode='{{.Mode}}'

function Check-User {
     $currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
     if (!$currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
        Write-Warning "baetyl installation should using administration principal"
        Break Script
     }
}

function Remove-DbFile {
    $CoreDbFile = Join-Path $DbPath 'store\core.db'
    $InitDbFile = Join-Path $DbPath 'init\store\core.db'
    if (Test-Path $InitDbFile) {
        Remove-Item $InitDbFile
    }
    if (Test-Path $CoreDbFile) {
        Remove-Item $CoreDbFile
    }
}

function Install-Baetyl {
    Remove-DbFile
    if ($Mode -eq "native") {
        Write-Host "baetyl install in native mode"
        if (Get-Command baetyl -ErrorAction SilentlyContinue) {
            baetyl delete
            baetyl apply -f "$Addr/v1/init/$($DeployYaml)?token=$Token" --skip-verify=true
        } else {
            Write-Warning "baetyl not installed yet, please install baetyl firstly"
            Break Script
        }
    }
}

Check-User
Install-Baetyl
