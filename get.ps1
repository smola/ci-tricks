#!/usr/bin/env pwsh

Write-Output "Looking up latest release"
# Fix SSL/TLS error on Windows PowerShell
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$response = try {
    Invoke-WebRequest "https://github.com/smola/ci-tricks/releases/latest" -UseBasicParsing -Method Head -MaximumRedirection 0 -ErrorAction 'Ignore'
} catch {
    $_.Exception.Response
}

# Windows PowerShell and PowerShell Core seem to vary their behaviour
# here: $response.Headers.Location might be a String or a System.Uri.
$tag = $response.Headers.Location.ToString().Split("/")[-1]

$path = (Get-Item -Path '.\' -Verbose).FullName
$output = "$path/ci-tricks.exe"
$arch = "amd64"
$url = "https://github.com/smola/ci-tricks/releases/download/$tag/ci-tricks_windows_$arch.exe"
Write-Output "Downloading $url"
(New-Object System.Net.WebClient).DownloadFile($url, $output)

Write-Output "Running $output"
# AppVeyor does not capture output from Start-Process, so we call
# the executable directly.
iex "$output"

if ($LastExitCode -ne 0) {
    $host.SetShouldExit($LastExitCode )
}