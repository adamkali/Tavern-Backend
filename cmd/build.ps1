# Get the commit message
$COMMIT_MESSAGE = $args[0]

# check if there is a VERSION.yaml file
# if not make one
if (!(Test-Path -Path ../cmd/VERSION.yaml)) {
    "major: 0" | Out-File -FilePath ./VERSION.yaml
    "minor: 1" | Out-File -FilePath ./VERSION.yaml -Append
}

# check if -M is set
if ($args[1] -eq "-M") {
    # increment major version
    Write-Host "Incrementing major version"
    # Get the current major version
    $MAJOR = Get-Content ./VERSION.yaml | Select-String -Pattern "major: " | Select-Object -ExpandProperty Line | Select-String -Pattern "\d+" | Select-Object -ExpandProperty Line
    # Increment the major version
    $MAJOR = $MAJOR + 1
    # Set the minor version to 0
    $MINOR = 0
} else {
    # increment minor version
    Write-Host "Incrementing minor version"
    # Get the current major version
    $MAJOR = Get-Content ./VERSION.yaml | Select-String -Pattern "major: " | Select-Object -ExpandProperty Line | Select-String -Pattern "\d+" | Select-Object -ExpandProperty Line
    # Get the current minor version
    $MINOR = Get-Content ./VERSION.yaml | Select-String -Pattern "minor: " | Select-Object -ExpandProperty Line | Select-String -Pattern "\d+" | Select-Object -ExpandProperty Line
    # Increment the minor version
    $MINOR = $MINOR + 1
}

# git add -A
# git commit -m "$COMMIT_MESSAGE"
# then checkout the Beor
git checkout Beor
git merge main

git add -A
git commit -m "$COMMIT_MESSAGE"

# git push origin beor
git push origin Beor

# git checkout main
git checkout main
