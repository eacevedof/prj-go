# install.ps1 -- instala el siem-client (agente del SIEM) en Windows, de principio a fin:
#   1. pide lo que falte: la url del panel y el token del agente (sale al crear el agente en el panel)
#   2. descarga el binario de ese panel y comprueba su sha256 (SHA256SUMS del propio panel)
#   3. escribe el .env (a donde envia y con que token) y prueba la conexion (config-sync + config-check)
#   4. lo deja arrancando con Windows (tarea programada como SYSTEM)
# Lo que se vigila no se configura aqui: son los monitores del proyecto, que el agente descarga solo del panel.
#
# Uso (PowerShell como administrador):
#   .\install.ps1                                   # pregunta url y token
#   .\install.ps1 -Url https://siem.dominio.com -Token <API_TOKEN>
#   .\install.ps1 -Dir C:\programs\gosiem-client -NoService -Yes
# O con doble clic en install.bat (lanza este script). Volver a lanzarlo actualiza el binario y conserva el token.
param(
    [string]$Url = "",
    [string]$Token = "",
    [string]$Dir = "C:\programs\gosiem-client",
    [switch]$NoService,
    [switch]$Yes
)
$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"
$TaskName = "siem-client"

function Write-Step([string]$text) { Write-Host ""; Write-Host "== $text" -ForegroundColor Cyan }
function Stop-WithError([string]$text) { Write-Host "ERROR: $text" -ForegroundColor Red; exit 1 }

$isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole(
    [Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $NoService -and -not $isAdmin) {
    Stop-WithError "hay que lanzarlo como administrador para dejarlo arrancando con Windows (o -NoService)"
}

Write-Step "1 - panel y agente"
if (-not $Url) { $Url = Read-Host "url del panel (p. ej. https://siem.dominio.com)" }
$Url = $Url.Trim().TrimEnd("/")
if ($Url -notmatch "^https?://") { Stop-WithError "la url del panel tiene que empezar por http:// o https://" }
$envPath = Join-Path $Dir ".env"
if (-not $Token -and (Test-Path $envPath)) {
    $existing = Select-String -Path $envPath -Pattern "^APP_SIEM_API_TOKEN=(.+)$" | Select-Object -First 1
    if ($existing) {
        $Token = $existing.Matches[0].Groups[1].Value
        Write-Host "se conserva el token del .env que ya habia en $Dir"
    }
}
if (-not $Token) {
    $secure = Read-Host "token del agente (sale una vez al crearlo en el panel)" -AsSecureString
    $Token = [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR($secure))
}
if (-not $Token) { Stop-WithError "sin token no se puede enviar nada: crea el agente en el panel" }
$serviceLabel = if ($NoService) { "no" } else { "si" }
Write-Host "panel: $Url - carpeta: $Dir - arranque con Windows: $serviceLabel"
if (-not $Yes) {
    $confirm = Read-Host "seguir? (s/n) [s]"
    if ($confirm -and $confirm -notmatch "^(s|si|y|yes)$") { Stop-WithError "cancelado" }
}

Write-Step "2 - binario"
try {
    $sums = (Invoke-WebRequest -UseBasicParsing -Uri "$Url/downloads/SHA256SUMS").Content
}
catch { Stop-WithError "no se puede descargar $Url/downloads/SHA256SUMS: comprueba la url y que este equipo llega al panel" }
if ($sums -is [byte[]]) { $sums = [Text.Encoding]::UTF8.GetString($sums) }
# la version mas alta publicada para Windows (por numeros: la 0.10.0 va detras de la 0.9.0)
$candidates = foreach ($line in ($sums -split "`n")) {
    if ($line -match "^([0-9a-f]{64})\s+\*?(siem-client-windows-(\d+\.\d+\.\d+)\.exe)\s*$") {
        [pscustomobject]@{ Sha256 = $Matches[1]; FileName = $Matches[2]; Version = [version]$Matches[3] }
    }
}
$binary = $candidates | Sort-Object Version | Select-Object -Last 1
if (-not $binary) { Stop-WithError "el panel no publica binario para Windows (make dist en el cliente)" }
Write-Host "descargando $($binary.FileName)"
New-Item -ItemType Directory -Force $Dir | Out-Null
# La carpeta no hereda permisos: bajo C:\ "Usuarios autentificados" podria modificar el .exe, el .env (token) y el
# siem-client.json, que la tarea programada ejecuta como SYSTEM. Solo SYSTEM y Administradores, y dueno Administradores
# (si la creo otro usuario antes, como dueno podria volver a cambiar los permisos). Sin admin (-NoService), el usuario
# actual en lugar de Administradores.
$systemSid = "*S-1-5-18"
$administratorsSid = "*S-1-5-32-544"
$ownerSid = if ($isAdmin) { $administratorsSid } else { "*" + [Security.Principal.WindowsIdentity]::GetCurrent().User.Value }
if ($isAdmin) { icacls $Dir /setowner $ownerSid /T /C /Q | Out-Null }
icacls $Dir /inheritance:r /grant:r "${systemSid}:(OI)(CI)F" "${administratorsSid}:(OI)(CI)F" "${ownerSid}:(OI)(CI)F" /Q | Out-Null
if ($LASTEXITCODE -ne 0) { Stop-WithError "no se pueden fijar los permisos de $Dir (icacls)" }
$tempFile = Join-Path $env:TEMP "siem-client-$([guid]::NewGuid()).exe"
try {
    Invoke-WebRequest -UseBasicParsing -Uri "$Url/downloads/$($binary.FileName)" -OutFile $tempFile
    $actualSha256 = (Get-FileHash $tempFile -Algorithm SHA256).Hash.ToLower()
    if ($actualSha256 -ne $binary.Sha256) { Stop-WithError "el sha256 no coincide: descarga corrupta o manipulada" }
    Write-Host "sha256 correcto"
    if (-not $NoService) { Stop-ScheduledTask -TaskName $TaskName -ErrorAction SilentlyContinue }
    Get-Process -Name "siem-client" -ErrorAction SilentlyContinue | Where-Object { $_.Path -eq (Join-Path $Dir "siem-client.exe") } | Stop-Process -Force
    Move-Item -Force $tempFile (Join-Path $Dir "siem-client.exe")
}
finally { Remove-Item -Force $tempFile -ErrorAction SilentlyContinue }
$exe = Join-Path $Dir "siem-client.exe"
Write-Host "instalado en $exe"

Write-Step "3 - configuracion"
# LF y ASCII: el .env que lee el agente
$envText = "APP_ENV=production`nAPP_SIEM_ENDPOINT_URL=$Url/api/v1/ingest-batches`nAPP_SIEM_API_TOKEN=$Token`n"
[IO.File]::WriteAllText($envPath, $envText, [Text.Encoding]::ASCII)
Push-Location $Dir
try {
    & $exe config-sync
    if ($LASTEXITCODE -ne 0) { Stop-WithError "el panel no acepta el token o no responde: revisa la url y el token (config-sync)" }
    & $exe config-check | Out-Null
    if ($LASTEXITCODE -ne 0) { Stop-WithError "la configuracion no es valida: $exe config-check" }
}
finally { Pop-Location }
Write-Host "conexion y configuracion correctas"

Write-Step "4 - arranque con Windows"
if ($NoService) {
    Write-Host "sin tarea programada (-NoService). Para arrancarlo: cd $Dir; .\siem-client.exe agent-run"
}
else {
    Unregister-ScheduledTask -TaskName $TaskName -Confirm:$false -ErrorAction SilentlyContinue
    $action = New-ScheduledTaskAction -Execute $exe -Argument "agent-run" -WorkingDirectory $Dir
    $trigger = New-ScheduledTaskTrigger -AtStartup
    $settings = New-ScheduledTaskSettingsSet -RestartCount 999 -RestartInterval (New-TimeSpan -Minutes 1) -ExecutionTimeLimit ([TimeSpan]::Zero)
    Register-ScheduledTask -TaskName $TaskName -Action $action -Trigger $trigger -Settings $settings -User "SYSTEM" -RunLevel Highest | Out-Null
    Start-ScheduledTask -TaskName $TaskName
    Write-Host "tarea programada $TaskName activa (arranca con Windows como SYSTEM)"
}

Write-Step "listo"
Write-Host "en el panel, el agente pasa a online en menos de un minuto"
