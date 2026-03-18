param(
    [string]$BaseUrl = "http://localhost:8080",
    [string]$ExecutivePassword = "Aa1!SmokeTestPass"
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function ConvertTo-JsonBody {
    param([Parameter(Mandatory = $true)]$InputObject)
    return $InputObject | ConvertTo-Json -Depth 10 -Compress
}

function Read-ErrorResponseBody {
    param($Exception)
    if ($null -eq $Exception.Response) {
        return ""
    }

    $stream = $Exception.Response.GetResponseStream()
    if ($null -eq $stream) {
        return ""
    }

    $reader = New-Object System.IO.StreamReader($stream)
    try {
        return $reader.ReadToEnd()
    }
    finally {
        $reader.Close()
    }
}

function Invoke-Api {
    param(
        [Parameter(Mandatory = $true)][string]$Method,
        [Parameter(Mandatory = $true)][string]$Path,
        [hashtable]$Headers = @{},
        $BodyObject = $null,
        [int[]]$ExpectedStatus = @(200)
    )

    $url = "$BaseUrl$Path"
    $requestBody = $null
    if ($null -ne $BodyObject) {
        $requestBody = ConvertTo-JsonBody -InputObject $BodyObject
    }

    $statusCode = 0
    $rawBody = ""

    try {
        $response = Invoke-WebRequest -Uri $url -Method $Method -Headers $Headers -ContentType "application/json" -Body $requestBody -UseBasicParsing
        $statusCode = [int]$response.StatusCode
        $rawBody = [string]$response.Content
    }
    catch {
        if ($null -eq $_.Exception.Response) {
            throw "Request failed for $Method ${Path}: $($_.Exception.Message)"
        }
        $statusCode = [int]$_.Exception.Response.StatusCode
        $rawBody = Read-ErrorResponseBody -Exception $_.Exception
    }

    $passed = $ExpectedStatus -contains $statusCode
    $marker = if ($passed) { "PASS" } else { "FAIL" }
    Write-Host "[$marker] $Method $Path -> $statusCode"

    if (-not $passed) {
        throw "Unexpected status for $Method $Path. Expected: $($ExpectedStatus -join ', '), actual: $statusCode, body: $rawBody"
    }

    $jsonBody = $null
    if ($rawBody -ne "") {
        try {
            $jsonBody = $rawBody | ConvertFrom-Json
        }
        catch {
            $jsonBody = $null
        }
    }

    return [pscustomobject]@{
        StatusCode = $statusCode
        RawBody = $rawBody
        JsonBody = $jsonBody
    }
}

$seed = [DateTime]::UtcNow.ToString("yyyyMMddHHmmss") + (Get-Random -Maximum 9999)
$execEmail = "smoke.exec.$seed@example.com"
$execStudentID = "EXE$seed"
$memberEmail = "smoke.member.$seed@example.com"
$memberStudentID = "MEM$seed"
$managedExecEmail = "managed.exec.$seed@example.com"
$managedExecStudentID = "MEX$seed"

Write-Host "Running smoke checks against $BaseUrl"

Invoke-Api -Method "GET" -Path "/swagger/index.html" -ExpectedStatus @(200)
Invoke-Api -Method "GET" -Path "/health" -ExpectedStatus @(200)

$registerExec = Invoke-Api -Method "POST" -Path "/api/register" -BodyObject @{
    name = "Smoke Exec"
    email = $execEmail
    student_id = $execStudentID
    password = $ExecutivePassword
    source_dashboard = "executives"
} -ExpectedStatus @(201)

$login = Invoke-Api -Method "POST" -Path "/api/login" -BodyObject @{
    identifier = $execEmail
    password = $ExecutivePassword
} -ExpectedStatus @(200)

$token = $login.JsonBody.token
if ($null -eq $token) {
    throw "Login response does not include token"
}

$accessToken = [string]$token.access_token
$refreshToken = [string]$token.refresh_token
if ($accessToken -eq "" -or $refreshToken -eq "") {
    throw "Login response missing access_token or refresh_token"
}

$refreshParts = $refreshToken -split "\.", 2
if ($refreshParts.Count -ne 2) {
    throw "Refresh token format is unexpected"
}
$refreshTokenID = $refreshParts[0]

Invoke-Api -Method "POST" -Path "/api/refresh" -BodyObject @{
    refresh_token_id = $refreshTokenID
    refresh_token = $refreshToken
} -ExpectedStatus @(200)

$authHeaders = @{
    Authorization = "Bearer $accessToken"
}

$memberCreate = Invoke-Api -Method "POST" -Path "/api/members" -Headers $authHeaders -BodyObject @{
    name = "Smoke Member"
    email = $memberEmail
    student_id = $memberStudentID
    password = $ExecutivePassword
    source_dashboard = "members"
} -ExpectedStatus @(201)

$memberID = [string]$memberCreate.JsonBody.member_id
if ($memberID -eq "") {
    throw "Create member response missing member_id"
}

Invoke-Api -Method "GET" -Path "/api/members" -Headers $authHeaders -ExpectedStatus @(200)
Invoke-Api -Method "GET" -Path "/api/members/$memberID" -Headers $authHeaders -ExpectedStatus @(200)
Invoke-Api -Method "PUT" -Path "/api/members/$memberID" -Headers $authHeaders -BodyObject @{
    course = "BSCS"
    contact_number = "09171234567"
} -ExpectedStatus @(200)
Invoke-Api -Method "DELETE" -Path "/api/members/$memberID" -Headers $authHeaders -ExpectedStatus @(200)

$createManagedExec = Invoke-Api -Method "POST" -Path "/api/executives" -Headers $authHeaders -BodyObject @{
    name = "Managed Exec"
    email = $managedExecEmail
    student_id = $managedExecStudentID
    password = $ExecutivePassword
} -ExpectedStatus @(201)

$managedExecutiveID = [string]$createManagedExec.JsonBody.executive_id
if ($managedExecutiveID -eq "") {
    throw "Create executive response missing executive_id"
}

Invoke-Api -Method "GET" -Path "/api/executives" -Headers $authHeaders -ExpectedStatus @(200)
Invoke-Api -Method "GET" -Path "/api/executives/$managedExecutiveID" -Headers $authHeaders -ExpectedStatus @(200)
Invoke-Api -Method "PUT" -Path "/api/executives/$managedExecutiveID" -Headers $authHeaders -BodyObject @{
    name = "Managed Exec Updated"
} -ExpectedStatus @(200)
Invoke-Api -Method "DELETE" -Path "/api/executives/$managedExecutiveID" -Headers $authHeaders -ExpectedStatus @(200)

$sessionToRevoke = Invoke-Api -Method "POST" -Path "/api/login" -BodyObject @{
    identifier = $execEmail
    password = $ExecutivePassword
} -ExpectedStatus @(200)

$revokeToken = [string]$sessionToRevoke.JsonBody.token.refresh_token
$revokeTokenParts = $revokeToken -split "\.", 2
if ($revokeTokenParts.Count -ne 2) {
    throw "Second refresh token format is unexpected"
}
$revokeSessionID = $revokeTokenParts[0]

Invoke-Api -Method "POST" -Path "/api/sessions/$revokeSessionID/revoke" -Headers $authHeaders -ExpectedStatus @(200)
Invoke-Api -Method "POST" -Path "/api/logout" -Headers $authHeaders -ExpectedStatus @(200)

Write-Host "All endpoint smoke checks completed successfully."
